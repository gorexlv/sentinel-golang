package base

import (
	"github.com/sentinel-group/sentinel-golang/core/node"
)

type BaseResource struct {
	defaultNode    *node.DefaultNode
	ctx            *Context
	createTimeInMs uint64
	endTimeInMs    uint64
}

func (b *BaseResource) EndTime() uint64 {
	return b.endTimeInMs
}

func (b *BaseResource) SetEndTime(endTimeInMs uint64) {
	b.endTimeInMs = endTimeInMs
}

func (b *BaseResource) CreateTime() uint64 {
	return b.createTimeInMs
}

func (b *BaseResource) SetCreateTime(createTimeInMs uint64) {
	b.createTimeInMs = createTimeInMs
}

func (b *BaseResource) Ctx() *Context {
	return b.ctx
}

func (b *BaseResource) SetCtx(ctx *Context) {
	b.ctx = ctx
}

func (b *BaseResource) DefaultNode() *node.DefaultNode {
	return b.defaultNode
}

func (b *BaseResource) SetDefaultNode(defaultNode *node.DefaultNode) {
	b.defaultNode = defaultNode
}

type ResourceWrapper struct {
	BaseResource
	// unique resource name
	ResourceName string
	//
	ResourceType int
}
