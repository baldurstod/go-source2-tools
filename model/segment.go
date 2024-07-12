package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type Segment struct {
	UniqueFrameIndex  int
	LocalElementMasks int
	LocalChannel      int
	Container         []byte
	decoderId         uint16
	bytesPerBone      uint16
	boneCount         uint16
	dataLength        uint16
	bones             []uint16
}

func (seg *Segment) initFromDatas(datas *kv3.Kv3Element) error {
	seg.UniqueFrameIndex, _ = datas.GetIntAttribute("m_nUniqueFrameIndex")
	seg.LocalElementMasks, _ = datas.GetIntAttribute("m_nLocalElementMasks")
	seg.LocalChannel, _ = datas.GetIntAttribute("m_nLocalChannel")

	container := datas.GetAttribute("m_container")
	if container == nil {
		return errors.New("unable to get segment container")
	}

	c, ok := container.([]byte)
	if !ok {
		return errors.New("can't convert container to byte array")
	}

	seg.Container = c

	seg.decoderId = uint16(c[0]) + (uint16(c[1]) << 8)
	seg.bytesPerBone = uint16(c[2]) + (uint16(c[3]) << 8)
	seg.boneCount = uint16(c[4]) + (uint16(c[5]) << 8)
	seg.dataLength = uint16(c[6]) + (uint16(c[7]) << 8)

	seg.bones = make([]uint16, 0, seg.boneCount)

	return nil
}

/*
	{
		m_nUniqueFrameIndex = 0
		m_nLocalElementMasks = 1024
		m_nLocalChannel = 0
		m_container = #[ 02 00 03 00 01 00 10 00 11 00 9E 4E 20 00 00 01 ]
	},
*/
