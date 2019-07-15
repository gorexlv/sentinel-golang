package base

import (
	"context"
	"github.com/sentinel-group/sentinel-golang/core/node"
)

type Context struct {
	name         string
	entranceNode node.DefaultNode
	curEntry     Entry
	origin       string
	context      context.Context
}
