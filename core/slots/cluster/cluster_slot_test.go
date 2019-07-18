package cluster

import (
	"github.com/sentinel-group/sentinel-golang/core/node"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestClusterBuilderSlot_Entry(t *testing.T) {
	slot := new(ClusterBuilderSlot)

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go testEntry(slot, wg, "test1")
		}

		for j := 0; j < 10; j++ {
			wg.Add(1)
			go testEntry(slot, wg, "test2")
		}
		wg.Wait()
	}

	value, ok := strResNodeMap.Load("test1")
	if !ok {
		t.Error("TestClusterBuilderSlot_Entry")
	}

	if _, ok := value.(*node.DefaultNode); !ok {
		t.Error("TestClusterBuilderSlot_Entry")
	}
}

func testEntry(slot *ClusterBuilderSlot, wg *sync.WaitGroup, resName string) {
	r := rand.Uint32() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)

	resource := new(base.ResourceWrapper)
	resource.ResourceName = resName
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
