package data

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/model"
	"github.com/sentinel-group/sentinel-golang/core/slog"
	"go.uber.org/zap"
	"math"
	"sync/atomic"
)

type MetricEventType int8

const (
	MetricEventPass MetricEventType = iota
	MetricEventBlock
	MetricEventSuccess
	MetricEventError
	MetricEventRt
	// hack for getting length of enum
	metricEventNum
)

/**
MetricBucket store the metric statistic of each event
(MetricEventPass、MetricEventBlock、MetricEventError、MetricEventSuccess、MetricEventRt)
*/
type MetricBucket struct {
	// value of statistic
	counters [metricEventNum]uint64
	minRt    uint64
}

func newEmptyMetricBucket() *MetricBucket {
	mb := &MetricBucket{
		minRt: math.MaxUint64,
	}
	return mb
}

func (mb *MetricBucket) Add(event MetricEventType, count uint64) {
	if event > metricEventNum {
		panic("event is bigger then metricEventNum")
	}
	if event == MetricEventRt {
		if count < mb.minRt {
			mb.minRt = count
		}
	}
	atomic.AddUint64(&mb.counters[event], count)
}

func (mb *MetricBucket) Get(event MetricEventType) uint64 {
	if event > metricEventNum {
		panic("event is bigger then metricEventNum")
	}
	return mb.counters[event]
}

func (mb *MetricBucket) MinRt() uint64 {
	return mb.minRt
}

// The implement of sliding window based on struct LeapArray
type SlidingWindow struct {
	data       *LeapArray
	BucketType string
}

func NewSlidingWindow(sampleCount uint32, intervalInMs uint32) *SlidingWindow {
	if intervalInMs%sampleCount != 0 {
		panic(fmt.Sprintf("invalid parameters, intervalInMs is %d, sampleCount is %d.", intervalInMs, sampleCount))
	}
	winLengthInMs := intervalInMs / sampleCount
	arr := make([]*WindowWrap, sampleCount)
	return &SlidingWindow{
		data: &LeapArray{
			windowLengthInMs: winLengthInMs,
			sampleCount:      sampleCount,
			intervalInMs:     intervalInMs,
			array:            arr,
		},
		BucketType: "metrics",
	}
}

func (sw *SlidingWindow) newEmptyBucket(startTime uint64) interface{} {
	return newEmptyMetricBucket()
}

func (sw *SlidingWindow) resetWindowTo(ww *WindowWrap, startTime uint64) (*WindowWrap, error) {
	ww.windowStart = startTime
	ww.value = newEmptyMetricBucket()
	return ww, nil
}

func (sw *SlidingWindow) Count(event MetricEventType) uint64 {
	_, err := sw.data.CurrentWindow(sw)
	if err != nil {
		slog.GetLog(slog.Record).Error("get current window", zap.Error(err))
	}
	count := uint64(0)
	for _, ww := range sw.data.Values() {
		mb, ok := ww.value.(*MetricBucket)
		if !ok {
			panic("value is not type *MetricBucket")
		}
		count += mb.Get(event)
	}
	return count
}

func (sw *SlidingWindow) AddCount(event MetricEventType, count uint64) {
	curWindow, err := sw.data.CurrentWindow(sw)
	if err != nil || curWindow == nil || curWindow.value == nil {
		slog.GetLog(slog.Record).Error("sliding window fail to record success")
		return
	}

	mb, ok := curWindow.value.(*MetricBucket)
	if !ok {
		panic("value is not type *MetricBucket")
	}
	mb.Add(event, count)
}

func (sw *SlidingWindow) MaxSuccess() uint64 {
	_, err := sw.data.CurrentWindow(sw)
	if err != nil {
		slog.GetLog(slog.Record).Error("get current window", zap.Error(err))
	}

	succ := uint64(0)
	for _, ww := range sw.data.Values() {
		mb, ok := ww.value.(*MetricBucket)
		if !ok {
			panic("value is not type *MetricBucket")
		}
		s := mb.Get(MetricEventSuccess)
		if err != nil {
			slog.GetLog(slog.Record).Error("mb is not *MetricBucket")
		}
		succ = uint64(math.Max(float64(succ), float64(s)))
	}
	return succ
}

func (sw *SlidingWindow) MinSuccess() uint64 {
	_, err := sw.data.CurrentWindow(sw)
	if err != nil {
		slog.GetLog(slog.Record).Error(err.Error())
	}

	succ := uint64(math.MaxUint64)
	for _, ww := range sw.data.Values() {
		mb, ok := ww.value.(*MetricBucket)
		if !ok {
			panic("value is not type *MetricBucket")
		}
		s := mb.Get(MetricEventSuccess)
		if s < succ {
			succ = s
		}
	}
	return succ
}

func (sw *SlidingWindow) MinRt() uint64 {
	_, err := sw.data.CurrentWindow(sw)
	if err != nil {
		slog.GetLog(slog.Record).Error("get current window", zap.Error(err))
	}

	min := uint64(math.MaxUint64)
	for _, ww := range sw.data.Values() {
		mb, ok := ww.value.(*MetricBucket)
		if !ok {
			panic("value is not type *MetricBucket")
		}
		s := mb.MinRt()
		if s < min {
			min = s
		}
	}

	return min
}

func (sw *SlidingWindow) Details() []*model.MetricNode {
	if _, e := sw.data.CurrentWindow(sw); e != nil {
		panic(e)
	}
	windowWraps := sw.data.Values()
	metricNodes := make([]*model.MetricNode, 0, len(windowWraps))
	for _, ww := range windowWraps {
		mb, ok := ww.value.(*MetricBucket)
		if !ok {
			panic("value is not type *MetricBucket")
		}
		metricNode := new(model.MetricNode)
		metricNode.BlockQps = mb.Get(MetricEventBlock)
		metricNode.ErrorQps = mb.Get(MetricEventError)
		metricNode.PassQps = mb.Get(MetricEventPass)
		succ := mb.Get(MetricEventSuccess)
		metricNode.SuccessQps = succ
		if succ != 0 {
			metricNode.Rt = mb.Get(MetricEventRt) / succ
		} else {
			metricNode.Rt = mb.Get(MetricEventRt)
		}
		metricNode.Timestamp = ww.windowStart
		metricNodes = append(metricNodes, metricNode)
	}
	return metricNodes
}

func (sw *SlidingWindow) GetWindowIntervalInMs() uint32 {
	return sw.data.intervalInMs
}

func (sw *SlidingWindow) GetWindowIntervalInSec() float64 {
	return float64(sw.data.intervalInMs) / 1000.0
}
