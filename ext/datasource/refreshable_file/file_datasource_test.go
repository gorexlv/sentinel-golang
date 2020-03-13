package refreshable_file

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/ext/datasource/plugin"
)

func TestFileDataSourceStarter(t *testing.T) {
	var ds datasource.DataSource = FileDataSourceStarter("../../../tests/testdata/extension/SystemRule.json")
	ds.SetDecoderBuilder(func(reader io.Reader) datasource.Decoder {
		return json.NewDecoder(reader)
	})

	ds.AddPropertyHandler(plugin.UpdateFlowRules)
	ds.AddPropertyHandler(plugin.UpdateSystemRules)
	ds.Initialize()
}
