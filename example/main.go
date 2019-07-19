package main

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"github.com/sentinel-group/sentinel-golang/core/system"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:10000", nil))
	}()

	rule := system.SystemRule{
		MaxGoroutine:      19,
		MaxMemUsedPercent: 69.0,
		MaxCpuUsedPercent: 69.0,
	}
	system.LoadRules(rule)
	wg := &sync.WaitGroup{}

	for a := 0; a < 100000; a++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go test(wg, "test1")
		}

		for j := 0; j < 10; j++ {
			wg.Add(1)
			go test(wg, "test2")
		}
		wg.Wait()
	}
}

func test(wg *sync.WaitGroup, res string) {
	entry, e := core.Entry(nil, res)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	r := rand.Uint32() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)

	if entry.Status == base.ResultStatusBlocked {
		fmt.Println("reason:", entry.BlockedReason)
	}
	if entry.Status == base.ResultStatusError {
		fmt.Println("reason:", entry.ErrorMsg)
	}
	if entry.Status == base.ResultStatusPass {
		_ = entry.Exit()
	}
	wg.Done()
}
