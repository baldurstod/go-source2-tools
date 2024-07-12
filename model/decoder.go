package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type Decoder struct {
	Name         string
	Version      int
	Type         int
	BytesPerBone int
}

func (dec *Decoder) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	dec.Name, ok = datas.GetStringAttribute("m_szName")
	if !ok {
		return errors.New("unable to get decoder name")
	}

	dec.Version, _ = datas.GetIntAttribute("m_nVersion")
	dec.Type, _ = datas.GetIntAttribute("m_nType")

	switch dec.Name {
	case "CCompressedStaticVector4D", "CCompressedFullVector4D":
		dec.BytesPerBone = 16
	case "CCompressedStaticFullVector3", "CCompressedFullVector3", "CCompressedDeltaVector3", "CCompressedAnimVector3":
		dec.BytesPerBone = 12
	case "CCompressedStaticVector2D", "CCompressedFullVector2D":
		dec.BytesPerBone = 8
	case "CCompressedAnimQuaternion", "CCompressedStaticVector3", "CCompressedStaticQuaternion", "CCompressedFullQuaternion":
		dec.BytesPerBone = 6
	case "CCompressedFullFloat", "CCompressedStaticFloat", "CCompressedStaticInt", "CCompressedFullInt", "CCompressedStaticColor32", "CCompressedFullColor32":
		dec.BytesPerBone = 4
	case "CCompressedFullShort":
		dec.BytesPerBone = 2
	case "CCompressedStaticChar", "CCompressedFullChar", "CCompressedStaticBool", "CCompressedFullBool":
		dec.BytesPerBone = 1
	default:
		return errors.New("unknown decoder type: " + dec.Name)
	}

	return nil
}

/*
	{
		m_szName = "CCompressedStaticFloat"
		m_nVersion = 0
		m_nType = 1
	},
*/
