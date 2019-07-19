package cluster

import (
	"github.com/sentinel-group/sentinel-golang/core/context"
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slog"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
	"go.uber.org/zap"
	"sync"
)

// 全局 StringResource -> DefaultNode
var strResNodeMap sync.Map

func StrResNodeMap() map[string]*node.DefaultNode {
	data := make(map[string]*node.DefaultNode)

	strResNodeMap.Range(func(key, value interface{}) bool {
		strRes, ok := key.(string)
		if !ok {
			panic("key is not key")
		}
		defaultNode, ok := value.(*node.DefaultNode)
		if !ok {
			panic("value is not *DefaultNode")
		}
		data[strRes] = defaultNode
		return true
	})

	return data
}

type ClusterBuilderSlot struct {
	chain.LinkedSlot
}

func (fs *ClusterBuilderSlot) Entry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	defaultNode, ok := strResNodeMap.Load(resWrapper.ResourceName)
	// not find make new DefaultNode
	if !ok {
		newNode := node.NewDefaultNode()
		actual, loaded := strResNodeMap.LoadOrStore(resWrapper.ResourceName, newNode)
		if !loaded {
			slog.GetLog(slog.Record).Info("new node add to map", zap.String("resource", resWrapper.ResourceName))
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

func (fs *ClusterBuilderSlot) Exit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error {
	return fs.FireExit(ctx, resWrapper, count)
}
