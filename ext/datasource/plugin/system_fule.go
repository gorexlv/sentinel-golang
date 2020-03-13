package plugin

import (
	"fmt"

	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/ext/datasource"
)

func UpdateSystemRules(decoder datasource.Decoder) error {
	var rules []*system.SystemRule
	if err := decoder.Decode(&rules); err != nil {
		return err
	}

	fmt.Printf("rules = %+v\n", rules)

	return system.LoadRules(rules)
}

