package plugin

import (
	"fmt"

	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/ext/datasource"
)

type CompoundRules struct {
	FlowRules []*flow.FlowRule
	SystemRules []*system.SystemRule
}

func UpdateCompoundRules(decoder datasource.Decoder) error {
	var cs CompoundRules
	if err := decoder.Decode(&cs); err != nil {
		return err
	}

	fmt.Printf("cs = %+v\n", cs)

	if len(cs.FlowRules) > 0 {
		_ = flow.LoadRules(cs.FlowRules)
	}

	if len(cs.SystemRules) > 0{
		_ = system.LoadRules(cs.SystemRules)
	}

	return nil
}


