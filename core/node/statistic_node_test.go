package node

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/sentinel-group/sentinel-golang/core/statistic/data"
)

const maxDif = 0.1

func TestStatisticNode_Request(t *testing.T) {
	sn := NewStatisticNode()
	sn.AddPassRequest(1)
	sn.AddBlockRequest(1)

	if sn.RequestInMinute() != 2 {
		t.Error("TestStatisticNode_Request")
	}
}

func TestStatisticNode_Qps(t *testing.T) {
	sn := NewStatisticNode()
	sn.AddPassRequest(1)
	sn.AddBlockRequest(1)
	sn.AddErrorRequest(1)
	sn.AddRtAndSuccess(10, 1)

	if math.Dim(sn.BlockQps(), 1.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
	if math.Dim(sn.PassQps(), 1.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
	if math.Dim(sn.TotalQps(), 2.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
	if math.Dim(sn.ErrorQps(), 1.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
	if math.Dim(sn.SuccessQps(), 1.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
	if math.Dim(sn.MaxSuccessQps(), 1.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}

	sn.AddBlockRequest(2)
	if math.Dim(sn.BlockQps(), 3.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}

	if math.Dim(sn.TotalQps(), 4.0) > maxDif {
		t.Error("TestStatisticNode_PassQps")
	}
}

func TestStatisticNode_AvgRt(t *testing.T) {
	sn := NewStatisticNode()
	sn.AddRtAndSuccess(1, 1)
	sn.AddRtAndSuccess(1, 1)

	// 等待 1s, 使 秒级刷新
	time.Sleep(time.Second * 1)
	sn.AddRtAndSuccess(2, 1)
	sn.AddRtAndSuccess(2, 1)

	if math.Dim(sn.AvgRtInSecond(), 2.0) > maxDif {
		t.Error("TestStatisticNode_AvgRt")
	}

	if math.Dim(sn.AvgRtInMinute(), 1.5) > maxDif {
		t.Error("TestStatisticNode_AvgRt")
	}
}

func TestStatisticNode_MinRt(t *testing.T) {
	sn := NewStatisticNode()
	sn.AddRtAndSuccess(1, 1)
	sn.AddRtAndSuccess(1, 1)
	// 等待 1s, 使 秒级刷新
	time.Sleep(time.Second * 1)
	sn.AddRtAndSuccess(2, 1)
	sn.AddRtAndSuccess(2, 1)

	if sn.MinRtInSecond() != 2 {
		t.Error("TestStatisticNode_MinRt")
	}
	if sn.MinRtInMinute() != 1 {
		t.Error("TestStatisticNode_MinRt")
	}
}

func TestStatisticNode_AddRtAndSuccess(t *testing.T) {
	sn := NewStatisticNode()
	sn.AddRtAndSuccess(1, 1)
	sn.AddRtAndSuccess(2, 1)

	if sn.MinRtInSecond() != 1 {
		t.Error("TestStatisticNode_Reset")
	}
}

func TestStatisticNode_GoroutineNum(t *testing.T) {
	sn := NewStatisticNode()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			sn.IncreaseGoroutineNum()
			wg.Done()
		}(&wg)
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			sn.DecreaseGoroutineNum()
			wg.Done()
		}(&wg)
	}

	wg.Wait()
	if sn.currentGoroutineNum != 0 {
		t.Error("TestStatisticNode_DecreaseGoroutineNum")
	}
}

func TestStatisticNode_Reset(t *testing.T) {
	sn := &StatisticNode{
		rollingCounterInSecond: data.NewSlidingWindow(2, 1000),
		rollingCounterInMinute: data.NewSlidingWindow(60, 60*1000),
		currentGoroutineNum:    100,
		lastFetchTime:          -1,
	}
	oldRollingCounterInSecond := sn.rollingCounterInSecond
	oldRollingCounterInMinute := sn.rollingCounterInMinute

	sn.Reset()

	if sn.rollingCounterInSecond == oldRollingCounterInSecond {
		t.Error("TestStatisticNode_Reset")
	}

	if sn.rollingCounterInMinute == oldRollingCounterInMinute {
		t.Error("TestStatisticNode_Reset")
	}
	if sn.currentGoroutineNum != 0 {
		t.Error("TestStatisticNode_Reset")
	}
	if sn.lastFetchTime != -1 {
		t.Error("TestStatisticNode_Reset")
	}
}

func TestStatisticNode_MutilGorotuine(t *testing.T) {

}
