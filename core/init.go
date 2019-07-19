package core

import (
	"encoding/json"
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/slots/cluster"
	"time"
)

func init() {

	go func() {
		defer func() {
			e := recover()
			fmt.Println(e)
		}()
		for {
			time.Sleep(time.Second)
			nodeMap := cluster.StrResNodeMap()
			for strRes, defNode := range nodeMap {
				metrics := defNode.Metrics()
				mts := fmt.Sprintf("res:%s,", strRes)
				for _, me := range metrics {
					bytes, e := json.Marshal(*me)
					if e != nil {
						panic(e)
					}
					meStr := string(bytes)
					mts = mts + meStr
				}
				fmt.Println(mts)
			}
		}
	}()
}
