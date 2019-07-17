package base

import (
	"context"
	"github.com/sentinel-group/sentinel-golang/core/node"
)

type Context struct {
	name         string
	entranceNode node.DefaultNode
	origin       string
	context      context.Context
}

func NewContext() *Context {
	return &Context{}
}
