package cluster

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
	"sync"
)

// 全局 StringResource -> DefaultNode
var strResMap sync.Map

type ClusterBuilderSlot struct {
	chain.LinkedSlot
}

func (fs *ClusterBuilderSlot) Entry(ctx *base.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	defaultNode, ok := strResMap.Load(resWrapper.ResourceName)
	// not find make new DefaultNode
	if !ok {
		newNode := node.NewDefaultNode()
		actual, loaded := strResMap.LoadOrStore(resWrapper.ResourceName, newNode)
		if !loaded {
			fmt.Println("new node add to strResMap")
		}
		defaultNode = actual
	}
	dfNode, ok := defaultNode.(*node.DefaultNode)
	if !ok {
		panic("type is not defaultNode")
	}
	resWrapper.SetDefaultNode(dfNode)
	return fs.FireEntry(ctx, resWrapper, count, prioritized)
}

func (fs *ClusterBuilderSlot) Exit(ctx *base.Context, resWrapper *base.ResourceWrapper, count int) error {
	return fs.FireExit(ctx, resWrapper, count)
}
