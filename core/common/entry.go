package common

import (
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
	"github.com/sentinel-group/sentinel-golang/core/util"
)

type Entry struct {
	*base.TokenResult

	resWrapper *base.ResourceWrapper
	slotChain  chain.SlotChain
}

func (t *Entry) SlotChain() chain.SlotChain {
	return t.slotChain
}

func (t *Entry) SetSlotChain(slotChain chain.SlotChain) {
	t.slotChain = slotChain
}

func (t *Entry) ResWrapper() *base.ResourceWrapper {
	return t.resWrapper
}

func (t *Entry) SetResWrapper(resWrapper *base.ResourceWrapper) {
	t.resWrapper = resWrapper
}

func (t *Entry) Exit() error {
	t.resWrapper.SetEndTime(util.GetTimeMilli())
	if t.slotChain != nil {
		return t.slotChain.Exit(t.resWrapper.Ctx(), t.resWrapper, 1)
	} else {
		return nil
	}
}
