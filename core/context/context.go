package context

import (
	"context"
)

type Context struct {
	name    string
	origin  string
	context context.Context
}

func NewContext() *Context {
	return &Context{}
}
