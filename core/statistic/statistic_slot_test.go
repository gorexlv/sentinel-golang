package statistic

import (
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStatisticSlot_Entry(t *testing.T) {
	defaultNode1 := node.NewDefaultNode()
	defaultNode2 := node.NewDefaultNode()

	slot := new(StatisticSlot)

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go testEntry(slot, wg, "test1", defaultNode1)
		}

		for j := 0; j < 10; j++ {
			wg.Add(1)
			go testEntry(slot, wg, "test2", defaultNode2)
		}
		wg.Wait()
	}

}

func testEntry(slot *StatisticSlot, wg *sync.WaitGroup, resName string, defaultNode *node.DefaultNode) {
	r := rand.Uint32() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)

	resource := new(base.ResourceWrapper)
	resource.ResourceName = resName
	resource.SetDefaultNode(defaultNode)
	result, e := slot.Entry(nil, resource, 1, false)
	if e != nil {
		panic(e)
	}

	if result.Status == base.ResultStatusPass {
		e := slot.Exit(nil, resource, 1)
		if e != nil {
			panic(e)
		}
	}
	wg.Done()
}
