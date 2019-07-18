package node

import (
	"github.com/sentinel-group/sentinel-golang/core/model"
	"github.com/sentinel-group/sentinel-golang/core/statistic/data"
	"github.com/sentinel-group/sentinel-golang/core/util"
	"sync/atomic"
)

const (
	sampleCountOfSecond  = 2
	intervalInMsOfSecond = 1000

	sampleCountOfMin  = 60
	intervalInMsOfMin = 60 * 1000
)

type StatisticNode struct {
	rollingCounterInSecond *data.SlidingWindow
	rollingCounterInMinute *data.SlidingWindow
	currentGoroutineNum    uint64
	lastFetchTime          uint64
}

/*
* The statistic node keep three kinds of real-time statistics metrics:
*
* metrics in second level ({@code rollingCounterInSecond})
* metrics in minute level ({@code rollingCounterInMinute})
* goroutine count
 */
func NewStatisticNode() *StatisticNode {
	return &StatisticNode{
		rollingCounterInSecond: data.NewSlidingWindow(sampleCountOfSecond, intervalInMsOfSecond),
		rollingCounterInMinute: data.NewSlidingWindow(sampleCountOfMin, intervalInMsOfMin),
		currentGoroutineNum:    0,
		lastFetchTime:          0,
	}
}

func (sn *StatisticNode) RequestInMinute() uint64 {
	return sn.PassInMinute() + sn.BlockInMinute()
}

func (sn *StatisticNode) PassInMinute() uint64 {
	return sn.rollingCounterInMinute.Count(data.MetricEventPass)
}

func (sn *StatisticNode) SuccessInMinute() uint64 {
	return sn.rollingCounterInMinute.Count(data.MetricEventSuccess)
}

func (sn *StatisticNode) BlockInMinute() uint64 {
	return sn.rollingCounterInMinute.Count(data.MetricEventBlock)
}

func (sn *StatisticNode) ErrorInMinute() uint64 {
	return sn.rollingCounterInMinute.Count(data.MetricEventError)
}

func (sn *StatisticNode) PassQps() float64 {
	return float64(sn.rollingCounterInSecond.Count(data.MetricEventPass)) / sn.rollingCounterInSecond.GetWindowIntervalInSec()
}

func (sn *StatisticNode) BlockQps() float64 {
	return float64(sn.rollingCounterInSecond.Count(data.MetricEventBlock)) / sn.rollingCounterInSecond.GetWindowIntervalInSec()
}

func (sn *StatisticNode) TotalQps() float64 {
	return sn.PassQps() + sn.BlockQps()
}

func (sn *StatisticNode) SuccessQps() float64 {
	return float64(sn.rollingCounterInSecond.Count(data.MetricEventSuccess)) / sn.rollingCounterInSecond.GetWindowIntervalInSec()
}

func (sn *StatisticNode) MaxSuccessQps() float64 {
	return float64(sn.rollingCounterInSecond.MaxSuccess()) / sn.rollingCounterInSecond.GetWindowIntervalInSec()
}

func (sn *StatisticNode) ErrorQps() float64 {
	return float64(sn.rollingCounterInSecond.Count(data.MetricEventError)) / sn.rollingCounterInSecond.GetWindowIntervalInSec()
}

func (sn *StatisticNode) AvgRtInSecond() float64 {
	succ := sn.rollingCounterInSecond.Count(data.MetricEventSuccess)
	if succ == 0 {
		return 0.0
	}
	return float64(sn.rollingCounterInSecond.Count(data.MetricEventRt)) / float64(succ)
}

func (sn *StatisticNode) AvgRtInMinute() float64 {
	succ := sn.rollingCounterInMinute.Count(data.MetricEventSuccess)
	if succ == 0 {
		return 0.0
	}
	return float64(sn.rollingCounterInMinute.Count(data.MetricEventRt)) / float64(succ)
}

func (sn *StatisticNode) MinRtInSecond() uint64 {
	return sn.rollingCounterInSecond.MinRt()
}

func (sn *StatisticNode) MinRtInMinute() uint64 {
	return sn.rollingCounterInMinute.MinRt()
}

func (sn *StatisticNode) CurGoroutineNum() uint64 {
	return sn.currentGoroutineNum
}

func (sn *StatisticNode) PreviousBlockQps() uint64 {
	panic("implement me")
}

func (sn *StatisticNode) PreviousPassQps() uint64 {
	panic("implement me")
}

func (sn *StatisticNode) Metrics() map[uint64]*model.MetricNode {
	currentTime := util.GetTimeMilli()
	currentTime = currentTime - currentTime%1000

	newLastFetchTime := sn.lastFetchTime
	metricNodes := sn.rollingCounterInMinute.Details()
	timeMetricNode := make(map[uint64]*model.MetricNode)
	for _, metricNode := range metricNodes {
		if sn.isNodeInTime(metricNode, currentTime) && sn.isValidMetricNode(metricNode) {
			timeMetricNode[metricNode.Timestamp] = metricNode
			if metricNode.Timestamp > newLastFetchTime {
				newLastFetchTime = metricNode.Timestamp
			}
		}
	}
	sn.lastFetchTime = newLastFetchTime
	return timeMetricNode
}

func (sn *StatisticNode) isNodeInTime(mNode *model.MetricNode, currentTime uint64) bool {
	return mNode.Timestamp > sn.lastFetchTime && mNode.Timestamp < currentTime
}

func (sn *StatisticNode) isValidMetricNode(node *model.MetricNode) bool {
	return node.PassQps > 0 ||
		node.BlockQps > 0 ||
		node.SuccessQps > 0 ||
		node.ErrorQps > 0 ||
		node.Rt > 0
}

func (sn *StatisticNode) AddPassRequest(count uint64) {
	sn.rollingCounterInSecond.AddCount(data.MetricEventPass, count)
	sn.rollingCounterInMinute.AddCount(data.MetricEventPass, count)
}

func (sn *StatisticNode) AddRtAndSuccess(rt uint64, success uint64) {
	sn.rollingCounterInSecond.AddCount(data.MetricEventSuccess, success)
	sn.rollingCounterInSecond.AddCount(data.MetricEventRt, rt)

	sn.rollingCounterInMinute.AddCount(data.MetricEventSuccess, success)
	sn.rollingCounterInMinute.AddCount(data.MetricEventRt, rt)
}

func (sn *StatisticNode) AddBlockRequest(count uint64) {
	sn.rollingCounterInSecond.AddCount(data.MetricEventBlock, count)
	sn.rollingCounterInMinute.AddCount(data.MetricEventBlock, count)
}

func (sn *StatisticNode) AddErrorRequest(count uint64) {
	sn.rollingCounterInSecond.AddCount(data.MetricEventError, count)
	sn.rollingCounterInMinute.AddCount(data.MetricEventError, count)
}

func (sn *StatisticNode) IncreaseGoroutineNum() {
	atomic.AddUint64(&sn.currentGoroutineNum, 1)
}

func (sn *StatisticNode) DecreaseGoroutineNum() {
	atomic.AddUint64(&sn.currentGoroutineNum, ^uint64(1-1))
}

func (sn *StatisticNode) Reset() {
	sn.rollingCounterInSecond = data.NewSlidingWindow(sampleCountOfSecond, intervalInMsOfSecond)
	sn.rollingCounterInMinute = data.NewSlidingWindow(sampleCountOfMin, intervalInMsOfMin)
	sn.currentGoroutineNum = 0
	sn.lastFetchTime = 0
}
