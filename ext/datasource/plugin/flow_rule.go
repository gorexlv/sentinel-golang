package plugin


import (
	"fmt"

	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/ext/datasource"
)

func UpdateFlowRules(decoder datasource.Decoder) error {
	var rules []*flow.FlowRule
	if err := decoder.Decode(&rules); err != nil {
		return err
	}

	fmt.Printf("rules = %+v\n", rules)

	return flow.LoadRules(rules)
}

