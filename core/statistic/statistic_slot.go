package statistic

import (
	"errors"
	"github.com/sentinel-group/sentinel-golang/core/context"
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
)

type StatisticSlot struct {
	chain.LinkedSlot
}

func (fs *StatisticSlot) Entry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	var err error
	// fire next slot
	result, err := fs.FireEntry(ctx, resWrapper, count, prioritized)

	if result == nil {
		panic(errors.New("result is nil"))
	}
	if err != nil {
		return base.NewResultError("err is not nil"), err
	}
	defaultNode := resWrapper.DefaultNode()
	if defaultNode == nil {
		panic("should not nil")
	}
	switch result.Status {
	case base.ResultStatusPass:
		processPass(defaultNode, count)
	case base.ResultStatusBlocked:
		processBlock(defaultNode, count)
	case base.ResultStatusWait:
		processWait(defaultNode, count)
	case base.ResultStatusError:
		processError(defaultNode, count)
	default:
		panic("should not occur")
	}

	return result, err
}

func processWait(node *node.DefaultNode, count int) {
	panic("should not occur")
}

func processError(node *node.DefaultNode, count int) {
	node.AddErrorRequest(uint64(count))
}

func processBlock(node *node.DefaultNode, count int) {
	node.AddBlockRequest(uint64(count))
}

func processPass(node *node.DefaultNode, count int) {
	node.IncreaseGoroutineNum()
	node.AddPassRequest(uint64(count))
}

func (fs *StatisticSlot) Exit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error {
	defaultNode := resWrapper.DefaultNode()
	if defaultNode == nil {
		panic("DefaultNode is nil")
	}
	rt := resWrapper.CreateTime() - resWrapper.EndTime()

	defaultNode.AddRtAndSuccess(rt, 1)
	defaultNode.DecreaseGoroutineNum()

	return fs.FireExit(ctx, resWrapper, count)
}
