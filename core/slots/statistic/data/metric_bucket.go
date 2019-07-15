package data

import (
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
