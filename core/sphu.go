package core

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/common"
	"github.com/sentinel-group/sentinel-golang/core/context"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
	"github.com/sentinel-group/sentinel-golang/core/slots/cluster"
	"github.com/sentinel-group/sentinel-golang/core/slots/flow"
	"github.com/sentinel-group/sentinel-golang/core/statistic"
	"github.com/sentinel-group/sentinel-golang/core/util"
)

type DefaultSlotChainBuilder struct {
}

func (dsc *DefaultSlotChainBuilder) Build() chain.SlotChain {
	linkedChain := chain.NewLinkedSlotChain()
	linkedChain.AddLast(new(cluster.ClusterBuilderSlot))
	linkedChain.AddLast(new(flow.FlowSlot))
	linkedChain.AddLast(new(statistic.StatisticSlot))
	// add all slot
	return linkedChain
}

func NewDefaultSlotChainBuilder() *DefaultSlotChainBuilder {
	return &DefaultSlotChainBuilder{}
}

var defaultChain chain.SlotChain

func init() {
	defaultChain = NewDefaultSlotChainBuilder().Build()
}

func Entry(ctx *context.Context, resource string) (*common.Entry, error) {
	if nil == ctx {
		ctx = context.NewContext()
	}

	resourceWrap := &base.ResourceWrapper{
		ResourceName: resource,
		ResourceType: base.INBOUND,
	}
	resourceWrap.SetCtx(ctx)
	resourceWrap.SetCreateTime(util.GetTimeMilli())

	result, e := defaultChain.Entry(ctx, resourceWrap, 1, false)
	if e != nil {
		fmt.Println(e.Error())
	}
	if result == nil {
		panic("result is nil")
	}
	if result.Status == base.ResultStatusBlocked {
		if e := defaultChain.Exit(nil, resourceWrap, 1); e != nil {
			fmt.Println(e.Error())
		}
	}

	// 组装返回的 token
	token := new(common.Entry)
	token.TokenResult = result
	token.SetResWrapper(resourceWrap)
	token.SetSlotChain(defaultChain)
	return token, nil
}
