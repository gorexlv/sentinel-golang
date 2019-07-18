package main

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
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

	wg := &sync.WaitGroup{}

	for a := 0; a < 100000; a++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go test(wg, "test1")
		}

		//for j := 0; j < 10; j++ {
		//	wg.Add(1)
		//	go test(wg,"test2")
		//}
		wg.Wait()

		fmt.Println("done")
	}
}

func test(wg *sync.WaitGroup, res string) {
	r := rand.Uint32() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)
	result, e := core.Entry(nil, res)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	if result.Status == base.ResultStatusBlocked {
		fmt.Println("reason:", result.BlockedReason)
	}
	if result.Status == base.ResultStatusError {
		fmt.Println("reason:", result.ErrorMsg)
	}
	if result.Status == base.ResultStatusPass {
		_ = result.Exit()
	}
	wg.Done()
}
