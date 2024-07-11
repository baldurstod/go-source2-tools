package model

import (
	"fmt"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type AnimBlock struct {
	decoders []Decoder
	segments []Segment
}

func newAnimBlock() *AnimBlock {
	return &AnimBlock{
		decoders: make([]Decoder, 0),
		segments: make([]Segment, 0),
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

	segments, _ := datas.GetKv3ValueArrayAttribute("m_segmentArray")
	block.segments = make([]Segment, 0, len(segments))

	for _, v := range segments {
		seg := Segment{}
		err := seg.initFromDatas(v.(*kv3.Kv3Element))
		if err != nil {
			return fmt.Errorf("cannot init segment in AnimBlock.initFromDatas: <%w>", err)
		}
		block.segments = append(block.segments, seg)
	}

	return nil
}
