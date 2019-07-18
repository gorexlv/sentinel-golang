package data

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/util"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	WindowLengthImMs uint32 = 200
	SampleCount      uint32 = 5
	IntervalInMs     uint32 = 1000
)

//Test sliding windows create windows
func TestNewWindow(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)
	time := util.GetTimeMilli()

	wr, err := slidingWindow.data.CurrentWindowWithTime(time, slidingWindow)
	if wr == nil {
		t.Errorf("Unexcepted error")
	}
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	if wr.windowLengthInMs != WindowLengthImMs {
		t.Errorf("Unexcepted error, winlength is not same")
	}
	if wr.windowStart != (time - time%uint64(WindowLengthImMs)) {
		t.Errorf("Unexcepted error, winlength is not same")
	}
	if wr.value == nil {
		t.Errorf("Unexcepted error, value is nil")
	}
	if slidingWindow.Count(MetricEventPass) != 0 {
		t.Errorf("Unexcepted error, pass value is invalid")
	}
}

// Test the logic get window start time.
func TestLeapArrayWindowStart(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)
	firstTime := util.GetTimeMilli()
	previousWindowStart := firstTime - firstTime%uint64(WindowLengthImMs)

	wr, err := slidingWindow.data.CurrentWindowWithTime(firstTime, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	if wr.windowLengthInMs != WindowLengthImMs {
		t.Errorf("Unexpected error, winLength is not same")
	}
	if wr.windowStart != previousWindowStart {
		t.Errorf("Unexpected error, winStart is not same")
	}
}

// test sliding window has multi windows
func TestWindowAfterOneInterval(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)
	firstTime := util.GetTimeMilli()
	previousWindowStart := firstTime - firstTime%uint64(WindowLengthImMs)

	wr, err := slidingWindow.data.CurrentWindowWithTime(firstTime, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	if wr.windowLengthInMs != WindowLengthImMs {
		t.Errorf("Unexpected error, winLength is not same")
	}
	if wr.windowStart != previousWindowStart {
		t.Errorf("Unexpected error, winStart is not same")
	}
	if wr.value == nil {
		t.Errorf("Unexcepted error")
	}
	mb, ok := wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 1)
	mb.Add(MetricEventBlock, 1)
	mb.Add(MetricEventSuccess, 1)
	mb.Add(MetricEventError, 1)

	if mb.Get(MetricEventPass) != 1 {
		t.Errorf("Unexcepted error")
	}
	if mb.Get(MetricEventBlock) != 1 {
		t.Errorf("Unexcepted error")
	}
	if mb.Get(MetricEventSuccess) != 1 {
		t.Errorf("Unexcepted error")
	}
	if mb.Get(MetricEventError) != 1 {
		t.Errorf("Unexcepted error")
	}

	middleTime := previousWindowStart + uint64(WindowLengthImMs)/2
	wr2, err := slidingWindow.data.CurrentWindowWithTime(middleTime, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	if wr2.windowStart != previousWindowStart {
		t.Errorf("Unexpected error, winStart is not same")
	}
	mb2, ok := wr2.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	if wr != wr2 {
		t.Errorf("Unexcepted error")
	}
	mb2.Add(MetricEventPass, 1)
	if mb.Get(MetricEventPass) != 2 {
		t.Errorf("Unexcepted error")
	}
	if mb.Get(MetricEventBlock) != 1 {
		t.Errorf("Unexcepted error")
	}

	lastTime := middleTime + uint64(WindowLengthImMs)/2
	wr3, err := slidingWindow.data.CurrentWindowWithTime(lastTime, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	if wr3.windowLengthInMs != WindowLengthImMs {
		t.Errorf("Unexpected error")
	}
	if (wr3.windowStart - uint64(WindowLengthImMs)) != previousWindowStart {
		t.Errorf("Unexpected error")
	}
	mb3, ok := wr3.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	if &mb3 == nil {
		t.Errorf("Unexcepted error")
	}

	if mb3.Get(MetricEventPass) != 0 {
		t.Errorf("Unexcepted error")
	}
	if mb3.Get(MetricEventBlock) != 0 {
		t.Errorf("Unexcepted error")
	}
}

func TestNTimeMultiGoroutineUpdateOneWindow(t *testing.T) {
	for i := 0; i < 1000; i++ {
		multiGoroutineUpdateWindows(t)
	}
}

func _task(wg *sync.WaitGroup, slidingWindow *SlidingWindow, ti uint64, t *testing.T, ct *uint64) {
	wr, err := slidingWindow.data.CurrentWindowWithTime(ti, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	mb, ok := wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 1)
	mb.Add(MetricEventBlock, 1)
	mb.Add(MetricEventSuccess, 1)
	mb.Add(MetricEventError, 1)
	atomic.AddUint64(ct, 1)
	wg.Done()
}

func multiGoroutineUpdateWindows(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)
	firstTime := util.GetTimeMilli()

	const GoroutineNum = 10
	wg := &sync.WaitGroup{}
	wg.Add(GoroutineNum)
	var cnt = uint64(0)
	for i := 0; i < GoroutineNum; i++ {
		go _task(wg, slidingWindow, firstTime, t, &cnt)
	}
	wg.Wait()
	ww, err := slidingWindow.data.CurrentWindowWithTime(firstTime, slidingWindow)
	if err != nil {
		t.Errorf("Unexcepted error")
	}
	mb, ok := ww.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	if mb.Get(MetricEventPass) != GoroutineNum {
		t.Errorf("Unexcepted error, infact, %d", mb.Get(MetricEventPass))
	}
	if mb.Get(MetricEventBlock) != GoroutineNum {
		t.Errorf("Unexcepted error, infact, %d", mb.Get(MetricEventBlock))
	}
}

func TestMultiGoroutineUpdateOneWindow(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)

	wg := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go testAddCountRandTime(t, wg, slidingWindow)
		}
		wg.Wait()
		fmt.Println("done")
	}
}

func testAddCountRandTime(t *testing.T, wg *sync.WaitGroup, slidingWindow *SlidingWindow) {
	r := rand.Uint32() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)

	currentTime := util.GetTimeMilli()
	wr, _ := slidingWindow.data.CurrentWindowWithTime(currentTime, slidingWindow)
	mb, ok := wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 1)
	mb.Add(MetricEventError, 1)

	wg.Done()
}

func TestSlidingWindow_Details(t *testing.T) {
	slidingWindow := NewSlidingWindow(SampleCount, IntervalInMs)
	currentTime := util.GetTimeMilli()
	wr, _ := slidingWindow.data.CurrentWindowWithTime(currentTime, slidingWindow)
	mb, ok := wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 1)
	mb.Add(MetricEventError, 1)

	time.Sleep(time.Millisecond * time.Duration(WindowLengthImMs))
	currentTime = util.GetTimeMilli()
	wr, _ = slidingWindow.data.CurrentWindowWithTime(currentTime, slidingWindow)
	mb, ok = wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 2)
	mb.Add(MetricEventError, 2)

	time.Sleep(time.Millisecond * time.Duration(WindowLengthImMs))
	currentTime = util.GetTimeMilli()
	wr, _ = slidingWindow.data.CurrentWindowWithTime(currentTime, slidingWindow)
	mb, ok = wr.value.(*MetricBucket)
	if !ok {
		t.Errorf("Unexcepted error")
	}
	mb.Add(MetricEventPass, 3)
	mb.Add(MetricEventError, 3)

	details := slidingWindow.Details()
	if len(details) != 3 {
		t.Errorf("Unexcepted error")
	}

	pass := uint64(0)
	err := uint64(0)
	for _, mn := range details {
		pass = pass + mn.PassQps
		err = err + mn.ErrorQps
	}

	if pass != 6 || err != 6 {
		t.Errorf("Unexcepted error")
	}
}
