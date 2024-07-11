package model

import (
	"fmt"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type AnimBlock struct {
	decoders []Decoder
}

func newAnimBlock() *AnimBlock {
	return &AnimBlock{
		decoders: make([]Decoder, 0),
	}
}

func (block *AnimBlock) initFromDatas(datas *kv3.Kv3Element) error {

	decoders, _ := datas.GetKv3ValueArrayAttribute("m_decoderArray")
	block.decoders = make([]Decoder, 0, len(decoders))

	for _, v := range decoders {
		dec := Decoder{}
		err := dec.initFromDatas(v.(*kv3.Kv3Element))
		if err != nil {
			return fmt.Errorf("cannot init decoder in AnimBlock.initFromDatas: <%w>", err)
		}
		block.decoders = append(block.decoders, dec)
	}

	return nil
}
