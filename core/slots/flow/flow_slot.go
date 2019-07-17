package flow

import (
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
)

type FlowSlot struct {
	chain.LinkedSlot
	RuleManager *RuleManager
}

func (fs *FlowSlot) Entry(ctx *base.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	defaultNode := resWrapper.DefaultNode()
	if defaultNode == nil {
		panic("DefaultNode is nil")
	}

	if fs.RuleManager == nil {
		return fs.FireEntry(ctx, resWrapper, count, false)
	}
	rules := fs.RuleManager.getRuleBySource(resWrapper.ResourceName)
	if len(rules) == 0 {
		return fs.FireEntry(ctx, resWrapper, count, false)
	}
	success := checkFlow(ctx, resWrapper, rules, defaultNode, count)
	if success {
		return fs.FireEntry(ctx, resWrapper, count, false)
	} else {
		return base.NewResultBlock("FlowSlot"), nil
	}
}

func (fs *FlowSlot) Exit(ctx *base.Context, resWrapper *base.ResourceWrapper, count int) error {
	return fs.FireExit(ctx, resWrapper, count)
}

func checkFlow(ctx *base.Context, resourceWrap *base.ResourceWrapper, rules []*rule, node *node.DefaultNode, count int) bool {
	if rules == nil {
		return true
	}
	for _, rule := range rules {
		if !canPass(ctx, resourceWrap, rule, node, uint32(count)) {
			return false
		}
	}
	return true
}

func canPass(ctx *base.Context, resourceWrap *base.ResourceWrapper, rule *rule, node *node.DefaultNode, count uint32) bool {
	return rule.controller_.CanPass(ctx, node, count)
}
