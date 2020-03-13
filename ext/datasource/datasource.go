package datasource

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

// The generic interface to describe the datasource
// Each DataSource instance listen in one property type.
type DataSource interface {
	// Add specified property handler in current datasource
	AddPropertyHandler(h PropertyHandler)
	SetDecoderBuilder(b DecoderBuilder)
	// ...
	Initialize()
	// Read original data from the data source.
	// return source bytes if succeed to read, if not, return error when reading
	ReadSource() ([]byte, error)
	// Close the data source.
	io.Closer
}

type (
	// 从数据源配置构建一个解码器
	DecoderBuilder func(io.Reader) Decoder
	// 延迟到数据源配置Handler中执行解码
	PropertyHandler func(Decoder) error
	Decoder interface {
		Decode(interface{}) error
	}
)

type Base struct {
	handlers     []PropertyHandler
	buildDecoder func(io.Reader) Decoder

	initOnce sync.Once
}

// SetDecoderBuilder reset datasource decoder's build method
// json.NewDecoder as default.
// You can set your own decoder builder to decode multi-rules
// or another data format like hcl/toml, by:
// datasourceInstance.SetDecoderBuilder(func(reader io.Reader) Decoder {
// 		return toml.NewDecoder(reader)
// }
// this can be called by biz user, datasource provider.
func (s *Base) SetDecoderBuilder(builder DecoderBuilder) {
	s.buildDecoder = builder
}

func (s *Base) AddPropertyHandler(h PropertyHandler) {
	s.initOnce.Do(s.init)
	s.handlers = append(s.handlers, h)
}

func (s Base) Handle(src []byte) error {
	s.initOnce.Do(s.init)
	for _, h := range s.handlers {
		decoder := s.buildDecoder(bytes.NewBuffer(src))
		if err := h(decoder); err != nil {
			return err
		}
	}
	return nil
}

func (s *Base) init() {
	s.handlers = make([]PropertyHandler, 0)
	if s.buildDecoder == nil {
		// default decoder builder
		s.buildDecoder = func(reader io.Reader) Decoder {
			return json.NewDecoder(reader)
		}
	}
}

