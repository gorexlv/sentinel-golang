package datasource

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

type Base struct {
	handlers []PropertyHandler
	decoderBuild func(io.Reader) Decoder

	initOnce sync.Once
}

func (s *Base) SetDecoderBuilder(decoderBuild DecoderBuilder) {
	s.decoderBuild = decoderBuild
}

func (s *Base) init() {
	s.handlers = make([]PropertyHandler, 0)
	if s.decoderBuild == nil {
		// default decoder builder
		s.decoderBuild = func(reader io.Reader) Decoder {
			return json.NewDecoder(reader)
		}
	}
}

func (s *Base) AddPropertyHandler(h PropertyHandler) {
	s.initOnce.Do(s.init)
	s.handlers = append(s.handlers, h)
}

func (s Base) Handle(src []byte) error {
	s.initOnce.Do(s.init)
	for _, h := range s.handlers {
		decoder := s.decoderBuild(bytes.NewBuffer(src))
		if err := h(decoder); err != nil {
			return err
		}
	}
	return nil
}

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
