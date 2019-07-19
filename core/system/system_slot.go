package system

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/context"
	"github.com/sentinel-group/sentinel-golang/core/slog"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/slots/chain"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"sync/atomic"
	"time"
)

// default value
const (
	MaxGoroutine      = uint64(2000)
	MaxMemUsedPercent = float64(99.9)
	MaxCpuUsedPercent = float64(99.9)
)

func init() {
	go func(r *RuntimeStatus) {
		for {
			cpuUsed, e := CpuUsed(time.Second)
			if e != nil {
				slog.GetLog(slog.Record).Error(e.Error())
			}
			r.SetCpuUsedPercent(cpuUsed)
			memUsed, e := MemUsed()
			if e != nil {
				slog.GetLog(slog.Record).Error(e.Error())
			}
			r.SetMemUsedPercent(memUsed)

			num := GoRoutineNum()
			r.SetGoroutine(num)
		}
	}(&RunStatus)
}

// system runtime status
type RuntimeStatus struct {
	goroutine      uint64
	memUsedPercent uint64
	cpuUsedPercent uint64
}

func (r *RuntimeStatus) CpuUsedPercent() float64 {
	memInt := atomic.LoadUint64(&r.memUsedPercent)
	return float64(memInt) / float64(1000)
}

func (r *RuntimeStatus) SetCpuUsedPercent(cpuUsedPercent float64) {
	cpuInt := uint64(cpuUsedPercent * 1000)
	atomic.StoreUint64(&r.cpuUsedPercent, cpuInt)
}

func (r *RuntimeStatus) MemUsedPercent() float64 {
	memInt := atomic.LoadUint64(&r.memUsedPercent)
	return float64(memInt) / float64(1000)
}

func (r *RuntimeStatus) SetMemUsedPercent(memUsedPercent float64) {
	memUseInt := uint64(memUsedPercent * 1000)
	atomic.StoreUint64(&r.memUsedPercent, memUseInt)
}

func (r *RuntimeStatus) Goroutine() uint64 {
	return r.goroutine
}

func (r *RuntimeStatus) SetGoroutine(goroutine uint64) {
	atomic.StoreUint64(&r.goroutine, goroutine)
}

// goable RuntimeStatus
var RunStatus RuntimeStatus

func MemUsed() (float64, error) {
	// 内存使用率
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0.0, err
	}
	return v.UsedPercent, nil
}

func CpuUsed(interval time.Duration) (float64, error) {
	cc, err := cpu.Percent(interval, false)
	if err != nil {
		return 0.0, err
	}
	return cc[0], nil
}

func GoRoutineNum() uint64 {
	return uint64(runtime.NumGoroutine())
}

//noinspection GoNameStartsWithPackageName
type SystemRule struct {
	MaxGoroutine      uint64
	MaxMemUsedPercent float64
	MaxCpuUsedPercent float64
}

var (
	maxGoroutine      = MaxGoroutine
	maxMemUsedPercent = MaxMemUsedPercent
	maxCpuUsedPercent = MaxCpuUsedPercent
)

// 在初始化的时候就加载规则
func LoadRules(systemRules []SystemRule) {
	for _, systemRule := range systemRules {
		maxGoroutine = systemRule.MaxGoroutine
		maxMemUsedPercent = systemRule.MaxMemUsedPercent
		maxCpuUsedPercent = systemRule.MaxCpuUsedPercent
	}
}

// 检查系统情况
func checkSystem(_ *base.ResourceWrapper) *base.TokenResult {
	// 系统协成个数
	curGoNum := RunStatus.Goroutine()
	if RunStatus.Goroutine() > maxGoroutine {
		return base.NewResultBlock(fmt.Sprintf("RunStatus.Goroutine:%d > maxGoroutine:%d", curGoNum, maxGoroutine))
	}
	// cpu 使用率,这里等会等改成异步获取
	curCpuUsed := RunStatus.CpuUsedPercent()
	if curCpuUsed > maxMemUsedPercent {
		return base.NewResultBlock(fmt.Sprintf("RunStatus.CpuUsedPercent:%f > MaxCpuUsedPercent:%f", curCpuUsed, maxCpuUsedPercent))
	}
	// 内存使用率
	memUsed := RunStatus.MemUsedPercent()
	if memUsed > maxMemUsedPercent {
		return base.NewResultBlock(fmt.Sprintf("RunStatus.MemUsedPercent:%f > MaxMemUsedPercent:%f", memUsed, maxMemUsedPercent))
	}
	return base.NewResultPass()
}

type SystemSlot struct {
	chain.LinkedSlot
}

func (ss *SystemSlot) Entry(ctx *context.Context, resWrapper *base.ResourceWrapper, count int, prioritized bool) (*base.TokenResult, error) {
	tokenResult := checkSystem(resWrapper)
	if tokenResult.Status == base.ResultStatusPass {
		return ss.FireEntry(ctx, resWrapper, count, prioritized)
	} else {
		return tokenResult, nil
	}
}

func (ss *SystemSlot) Exit(ctx *context.Context, resWrapper *base.ResourceWrapper, count int) error {
	return ss.FireExit(ctx, resWrapper, count)
}
