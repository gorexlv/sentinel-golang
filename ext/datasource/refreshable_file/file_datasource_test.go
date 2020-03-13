package refreshable_file

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/ext/datasource/plugin"
)

func TestFileDataSourceStarter_SimpleRule(t *testing.T) {
	var ds datasource.DataSource = FileDataSourceStarter("../../../tests/testdata/extension/SystemRule.json")
	ds.SetDecoderBuilder(func(reader io.Reader) datasource.Decoder {
		// customize decoder
		return json.NewDecoder(reader)
	})

	// ds.AddPropertyHandler(plugin.UpdateFlowRules)
	ds.AddPropertyHandler(plugin.UpdateSystemRules)
	ds.Initialize()
}

func TestFileDataSourceStarter_CompoundRule(t *testing.T) {
	var ds datasource.DataSource = FileDataSourceStarter("../../../tests/testdata/extension/CompundSource.json")
	ds.SetDecoderBuilder(func(reader io.Reader) datasource.Decoder {
		// customize decoder
		return json.NewDecoder(reader)
	})

	ds.AddPropertyHandler(plugin.UpdateCompoundRules)
	ds.Initialize()
}

func TestFileDataSourceStarter_CompoundRule_CustomizeHandler(t *testing.T) {
	updater := func (decoder datasource.Decoder) error {
		var compound map[string]json.RawMessage
		if err := decoder.Decode(&compound); err != nil {
			return err
		}

		_ = plugin.UpdateFlowRules(json.NewDecoder(bytes.NewBuffer(compound["flowRules"])))
		_ = plugin.UpdateSystemRules(json.NewDecoder(bytes.NewBuffer(compound["systemRules"])))
		return nil
	}

	var ds datasource.DataSource = FileDataSourceStarter("../../../tests/testdata/extension/CompundSource.json")
	ds.AddPropertyHandler(updater)
	ds.Initialize()

}