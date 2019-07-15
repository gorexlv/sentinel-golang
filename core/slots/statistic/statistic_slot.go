package statistic

import (
	"errors"
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
)

type StatisticSlot struct {
	chain.LinkedSlot
}

func (fs *StatisticSlot) Entry(ctx *base.Context, resWrapper *base.ResourceWrapper, node *node.DefaultNode, count int, prioritized bool) (*base.TokenResult, error) {
	var r *base.TokenResult
	var err error
	defer func() {
		if e := recover(); e != nil {
			r = base.NewResultError("StatisticSlot")
			err = errors.New("panic occur")
		}
	}()
	// fire next slot
	result, err := fs.FireEntry(ctx, resWrapper, node, count, prioritized)
	if result == nil {
		return base.NewResultError("result is nil"), err
	}
	if err != nil {
		return base.NewResultError("err is not nil"), err
	}

	if result.Status == base.ResultStatusError {
		// TO DO
	}
	if result.Status == base.ResultStatusPass {
		node.AddPassRequest(1)
	}
	return result, err
}

func (fs *StatisticSlot) Exit(ctx *base.Context, resourceWrapper *base.ResourceWrapper, count int) error {
	return fs.FireExit(ctx, resourceWrapper, count)
}
