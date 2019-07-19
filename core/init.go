package core

import (
	"encoding/json"
	"fmt"
	"github.com/sentinel-group/sentinel-golang/core/slog"
	"github.com/sentinel-group/sentinel-golang/core/slots/cluster"
	"go.uber.org/zap"
	"time"
)

func init() {
	metricToLog()
}

func metricToLog() {
	go func() {
		defer func() {
			e := recover()
			slog.GetLog(slog.Record).Error("metricToLog error", zap.Any("e", e))
		}()
		for {
			time.Sleep(time.Second)
			nodeMap := cluster.StrResNodeMap()
			var allResInfo string
			for strRes, defNode := range nodeMap {
				metrics := defNode.Metrics()
				oneResMts := fmt.Sprintf("ResName:%s,", strRes)
				for _, oneMetric := range metrics {
					bytes, e := json.Marshal(*oneMetric)
					if e != nil {
						panic(e)
					}
					oneMeStr := string(bytes)
					oneResMts = oneResMts + oneMeStr
				}
				allResInfo = allResInfo + oneResMts + "||"
			}
			slog.GetLog(slog.MetricLog).Info(allResInfo)
		}
	}()
}
