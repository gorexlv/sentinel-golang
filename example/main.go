package main

import (
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core"
	"github.com/sentinel-group/sentinel-golang/core/slots/base"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("=================start=================")
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go test(wg)
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go test2(wg)
	}
	wg.Wait()
	fmt.Println("=================end=================")
}

func test(wg *sync.WaitGroup) {
	rand.Seed(1000)
	r := rand.Int63() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)
	result, e := core.Entry(nil, "test")
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
	time.Sleep(time.Duration(r) * time.Millisecond)
	wg.Done()
}

func test2(wg *sync.WaitGroup) {
	rand.Seed(1000)
	r := rand.Int63() % 10
	time.Sleep(time.Duration(r) * time.Millisecond)
	result, e := core.Entry(nil, "test2")
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
	time.Sleep(time.Duration(r) * time.Millisecond)
	wg.Done()
}
