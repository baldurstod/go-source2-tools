package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/baldurstod/go-vector"
)

type DecodeKeyBone struct {
	name       string
	parent     int
	pos        vector.Vector3[float32]
	quat       vector.Quaternion[float32]
	scale      float32
	alignement vector.Quaternion[float32]
	flags      int
}

func newDecodeKeyBone() *DecodeKeyBone {
	return &DecodeKeyBone{}
}

func (bone *DecodeKeyBone) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	bone.name, ok = datas.GetStringAttribute("m_name")
	if !ok {
		return errors.New("unable to get bone name")
	}

	bone.parent, ok = datas.GetIntAttribute("m_parent")
	if !ok {
		return errors.New("unable to get bone parent")
	}

	v, err := datas.GetVec3Attribute("m_pos")
	if err != nil {
		return fmt.Errorf("unable to get bone position: <%w>", err)
	}
	log.Println(v)

	q, err := datas.GetQuatAttribute("m_quat")
	if err != nil {
		return fmt.Errorf("unable to get bone quaternion: <%w>", err)
	}
	log.Println(q)

	return nil
}

/*
	{
		m_name = "Root_0"
		m_parent = -1
		m_pos = [ 13.248837, 0.0, 101.805443 ]
		m_quat = [ 0.0, 0.0, 0.0, 1.0 ]
		m_scale = 1.0
		m_qAlignment = [ 0.0, 0.0, 0.0, 0.0 ]
		m_flags = 64
	},
*/
