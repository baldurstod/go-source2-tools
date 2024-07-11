package animations

import (
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
		dec.initFromDatas(v.(*kv3.Kv3Element))
		block.decoders = append(block.decoders, dec)
	}
	/*dec.Name, ok = datas.GetStringAttribute("m_szName")
	if !ok {
		return errors.New("unable to get decoder name")
	}

	dec.Version, _ = datas.GetIntAttribute("m_nVersion")
	dec.Type, _ = datas.GetIntAttribute("m_nType")*/

	return nil
}
