package plugin


import (
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/ext/datasource"
)

func UpdateFlowRules(decoder datasource.Decoder) error {
	var rules []*flow.FlowRule
	if err := decoder.Decode(&rules); err != nil {
		return err
	}

	return flow.LoadRules(rules)
}

