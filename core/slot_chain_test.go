package core

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	log.Init()
}

type StatPrepareSlotMock1 struct {
	Name string
}

func (spl *StatPrepareSlotMock1) Prepare(ctx *Context) {
	fmt.Println(spl.Name)
	return
}

func TestSlotChain_addStatPrepareSlotFirstAndLast(t *testing.T) {
	sc := newSlotChain()
	for i := 9; i >= 0; i-- {
		sc.addStatPrepareSlotFirst(&StatPrepareSlotMock1{
			Name: "mock2" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addStatPrepareSlotFirst(&StatPrepareSlotMock1{
			Name: "mock1" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addStatPrepareSlotLast(&StatPrepareSlotMock1{
			Name: "mock3" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addStatPrepareSlotFirst(&StatPrepareSlotMock1{
			Name: "mock" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addStatPrepareSlotLast(&StatPrepareSlotMock1{
			Name: "mock4" + strconv.Itoa(i),
		})
	}

	spSlice := sc.statPres
	if len(spSlice) != 50 {
		t.Error("len error")
	}

	for idx, slot := range spSlice {
		n := "mock" + strconv.Itoa(idx)
		spsm, ok := slot.(*StatPrepareSlotMock1)
		if !ok {
			t.Error("type error")
		}
		reflect.DeepEqual(n, spsm.Name)
	}
}

type RuleCheckSlotMock1 struct {
	Name string
}

func (rcs *RuleCheckSlotMock1) Check(ctx *Context) *RuleCheckResult {
	fmt.Println(rcs.Name)
	return nil
}
func TestSlotChain_addRuleCheckSlotFirstAndLast(t *testing.T) {
	sc := newSlotChain()
	for i := 9; i >= 0; i-- {
		sc.addRuleCheckSlotFirst(&RuleCheckSlotMock1{
			Name: "mock2" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addRuleCheckSlotFirst(&RuleCheckSlotMock1{
			Name: "mock1" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addRuleCheckSlotLast(&RuleCheckSlotMock1{
			Name: "mock3" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addRuleCheckSlotFirst(&RuleCheckSlotMock1{
			Name: "mock" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addRuleCheckSlotLast(&RuleCheckSlotMock1{
			Name: "mock4" + strconv.Itoa(i),
		})
	}

	spSlice := sc.ruleChecks
	if len(spSlice) != 50 {
		t.Error("len error")
	}

	for idx, slot := range spSlice {
		n := "mock" + strconv.Itoa(idx)
		spsm, ok := slot.(*RuleCheckSlotMock1)
		if !ok {
			t.Error("type error")
		}
		reflect.DeepEqual(n, spsm.Name)
	}
}

type StatSlotMock1 struct {
	Name string
}

func (ss *StatSlotMock1) OnEntryPassed(ctx *Context) {
	fmt.Println(ss.Name)
}
func (ss *StatSlotMock1) OnEntryBlocked(ctx *Context) {
	fmt.Println(ss.Name)
}
func (ss *StatSlotMock1) OnCompleted(ctx *Context) {
	fmt.Println(ss.Name)
}
func TestSlotChain_addStatSlotFirstAndLast(t *testing.T) {
	sc := newSlotChain()
	for i := 9; i >= 0; i-- {
		sc.addStatSlotFirst(&StatSlotMock1{
			Name: "mock2" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addStatSlotFirst(&StatSlotMock1{
			Name: "mock1" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addStatSlotLast(&StatSlotMock1{
			Name: "mock3" + strconv.Itoa(i),
		})
	}
	for i := 9; i >= 0; i-- {
		sc.addStatSlotFirst(&StatSlotMock1{
			Name: "mock" + strconv.Itoa(i),
		})
	}
	for i := 0; i < 10; i++ {
		sc.addStatSlotLast(&StatSlotMock1{
			Name: "mock4" + strconv.Itoa(i),
		})
	}

	spSlice := sc.stats
	if len(spSlice) != 50 {
		t.Error("len error")
	}

	for idx, slot := range spSlice {
		n := "mock" + strconv.Itoa(idx)
		spsm, ok := slot.(*StatSlotMock1)
		assert.True(t, ok, "slot type must be StatSlotMock1")
		if !ok {
			t.Error("type error")
		}
		reflect.DeepEqual(n, spsm.Name)
	}
}

type ResourceBuilderSlotMock struct {
	mock.Mock
}

func (m *ResourceBuilderSlotMock) Prepare(ctx *Context) {
	m.Called(ctx)
	return
}

type FlowSlotMock struct {
	mock.Mock
}

func (m *FlowSlotMock) Check(ctx *Context) *RuleCheckResult {
	m.Called(ctx)
	return NewSlotResultPass()
}

type DegradeSlotMock struct {
	mock.Mock
}

func (m *DegradeSlotMock) Check(ctx *Context) *RuleCheckResult {
	m.Called(ctx)
	return NewSlotResultPass()
}

type StatisticSlotMock struct {
	mock.Mock
}

func (m *StatisticSlotMock) OnEntryPassed(ctx *Context) {
	m.Called(ctx)
	return
}
func (m *StatisticSlotMock) OnEntryBlocked(ctx *Context) {
	m.Called(ctx)
	return
}
func (m *StatisticSlotMock) OnCompleted(ctx *Context) {
	m.Called(ctx)
	return
}

func TestSlotChain_Entry(t *testing.T) {
	sc := newSlotChain()
	ctx := sc.GetContext()
	rw := &ResourceWrapper{
		ResourceName: "abc",
		FlowType:     InBound,
	}
	ctx.ResWrapper = rw
	ctx.Node = FindNode(rw)
	ctx.Count = 1
	ctx.Entry = NewCtEntry(ctx, rw, sc, ctx.Node)

	rbs := &ResourceBuilderSlotMock{}
	fsm := &FlowSlotMock{}
	dsm := &DegradeSlotMock{}
	ssm := &StatisticSlotMock{}
	sc.addStatPrepareSlotFirst(rbs)
	sc.addRuleCheckSlotFirst(fsm)
	sc.addRuleCheckSlotFirst(dsm)
	sc.addStatSlotFirst(ssm)

	rbs.On("Prepare", mock.Anything).Return()
	fsm.On("Check", mock.Anything).Return(NewSlotResultPass())
	dsm.On("Check", mock.Anything).Return(NewSlotResultPass())
	ssm.On("OnEntryPassed", mock.Anything).Return()

	sc.Entry(ctx)

	rbs.AssertNumberOfCalls(t, "Prepare", 1)
	fsm.AssertNumberOfCalls(t, "Check", 1)
	dsm.AssertNumberOfCalls(t, "Check", 1)
	ssm.AssertNumberOfCalls(t, "OnEntryPassed", 1)
	ssm.AssertNumberOfCalls(t, "OnEntryBlocked", 0)
	ssm.AssertNumberOfCalls(t, "OnCompleted", 0)
}
