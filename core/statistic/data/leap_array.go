package data

import (
	"errors"
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/util"
	"runtime"
)

type WindowWrap struct {
	windowLengthInMs uint32
	windowStart      uint64
	value            interface{}
}

func (ww *WindowWrap) resetTo(startTime uint64) {
	ww.windowStart = startTime
}

func (ww *WindowWrap) isTimeInWindow(timeMillis uint64) bool {
	return ww.windowStart <= timeMillis && timeMillis < ww.windowStart+uint64(ww.windowLengthInMs)
}

// The basic data structure of sliding windows
//
type LeapArray struct {
	windowLengthInMs uint32
	sampleCount      uint32
	intervalInMs     uint32
	array            []*WindowWrap     //实际保存的数据
	mux              util.TriableMutex // lock
}

func (la *LeapArray) WindowLengthInMs() uint32 {
	return la.windowLengthInMs
}

func (la *LeapArray) CurrentWindow(sw BucketGenerator) (*WindowWrap, error) {
	return la.CurrentWindowWithTime(util.GetTimeMilli(), sw)
}

func (la *LeapArray) CurrentWindowWithTime(timeMillis uint64, sw BucketGenerator) (*WindowWrap, error) {
	if timeMillis < 0 {
		return nil, errors.New("timeMillion is less than 0")
	}

	idx := la.calculateTimeIdx(timeMillis)
	windowStart := la.calculateStartTime(timeMillis)

	for {
		old := la.array[idx]
		if old == nil {
			newWrap := &WindowWrap{
				windowLengthInMs: la.windowLengthInMs,
				windowStart:      windowStart,
				value:            sw.newEmptyBucket(windowStart),
			}
			if la.mux.TryLock() && la.array[idx] == nil {
				la.array[idx] = newWrap
				la.mux.Unlock()
				return la.array[idx], nil
			} else {
				runtime.Gosched()
			}
		} else if windowStart == old.windowStart {
			return old, nil
		} else if windowStart > old.windowStart {
			// reset WindowWrap
			if la.mux.TryLock() {
				old, _ = sw.resetWindowTo(old, windowStart)
				la.mux.Unlock()
				return old, nil
			} else {
				runtime.Gosched()
			}
		} else if windowStart < old.windowStart {
			// Should not go through here,
			return nil, errors.New(fmt.Sprintf("provided time timeMillis=%d is already behind old.windowStart=%d", windowStart, old.windowStart))
		}
	}
}

func (la *LeapArray) calculateTimeIdx(timeMillis uint64) uint32 {
	timeId := (int)(timeMillis / uint64(la.windowLengthInMs))
	return uint32(timeId % len(la.array))
}

func (la *LeapArray) calculateStartTime(timeMillis uint64) uint64 {
	return timeMillis - (timeMillis % uint64(la.windowLengthInMs))
}

//  Get all the bucket in sliding window for current time;
func (la *LeapArray) Values() []*WindowWrap {
	return la.valuesWithTime(util.GetTimeMilli())
}

func (la *LeapArray) valuesWithTime(timeMillis uint64) []*WindowWrap {
	if timeMillis <= 0 {
		return nil
	}
	wwp := make([]*WindowWrap, 0)
	for _, wwPtr := range la.array {
		if wwPtr == nil || la.isWindowDeprecated(timeMillis, wwPtr) {
			continue
		}
		newWW := &WindowWrap{
			windowLengthInMs: wwPtr.windowLengthInMs,
			windowStart:      wwPtr.windowStart,
			value:            wwPtr.value,
		}
		wwp = append(wwp, newWW)
	}
	return wwp
}

func (la *LeapArray) isWindowDeprecated(timeMillis uint64, ww *WindowWrap) bool {
	return timeMillis-ww.windowStart > uint64(la.intervalInMs)
}

type BucketGenerator interface {
	// 根据开始时间，创建一个新的统计bucket, bucket的具体数据结构可以有多个
	newEmptyBucket(startTime uint64) interface{}

	// 将窗口ww重置startTime和空的统计bucket
	resetWindowTo(ww *WindowWrap, startTime uint64) (*WindowWrap, error)
}
