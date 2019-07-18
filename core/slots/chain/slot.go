package chain

import (
	"github.com/sentinel-group/sentinel-golang/core/context"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
)

// a solt
type Slot interface {
	/**
	 * Entrance of this slots.
	 */
	Entry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error)

	Exit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error

	// 传递进入
	FireEntry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error)

	// 传递退出
	FireExit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error

	GetNext() Slot

	SetNext(next Slot)
}

// a slot can make slot compose linked
type LinkedSlot struct {
	// next linkedSlot
	next Slot
}

// 传递退出
func (s *LinkedSlot) Entry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	return s.FireEntry(ctx, resWrapper, count, prioritized)
}

// 传递进入
func (s *LinkedSlot) Exit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error {
	return s.FireExit(ctx, resWrapper, count)
}

// 传递进入, 没有下一个就返回 ResultStatusPass
func (s *LinkedSlot) FireEntry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	if s.next != nil {
		return s.next.Entry(ctx, resWrapper, count, prioritized)
	}
	return base.NewResultPass(), nil
}

// 传递退出，没有下一个就返回
func (s *LinkedSlot) FireExit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error {
	if s.next != nil {
		return s.next.Exit(ctx, resWrapper, count)
	} else {
		return nil
	}
}

func (s *LinkedSlot) GetNext() Slot {
	return s.next
}

func (s *LinkedSlot) SetNext(next Slot) {
	s.next = next
}
