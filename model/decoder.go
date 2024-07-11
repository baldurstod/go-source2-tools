package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type Decoder struct {
	Name    string
	Version int
	Type    int
}

func (dec *Decoder) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	dec.Name, ok = datas.GetStringAttribute("m_szName")
	if !ok {
		return errors.New("unable to get decoder name")
	}

	dec.Version, _ = datas.GetIntAttribute("m_nVersion")
	dec.Type, _ = datas.GetIntAttribute("m_nType")

	return nil
}

/*
	{
		m_szName = "CCompressedStaticFloat"
		m_nVersion = 0
		m_nType = 1
	},
*/
