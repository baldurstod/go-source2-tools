package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/x448/float16"
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

func (dec *Decoder) decode(reader *bytes.Reader, frameIndex int, boneCount int) error {
	switch dec.Name {
	case "CCompressedStaticVector3":
		//reader.Seek(int64(8+boneCount*(2+frameIndex*dec.BytesPerBone)), io.SeekStart)
		reader.Seek(int64(8+boneCount*2), io.SeekStart)
		v := [3]uint16{}

		err := binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return fmt.Errorf("failed to read segment bone count: <%w>", err)
		}

		log.Println(float16.Frombits(v[0]), float16.Frombits(v[1]), float16.Frombits(v[2]))

		return nil
	default:
		return errors.New("unknown decoder type: " + dec.Name)
	}
}

/*
	{
		m_szName = "CCompressedStaticFloat"
		m_nVersion = 0
		m_nType = 1
	},
*/
