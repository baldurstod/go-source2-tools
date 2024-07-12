package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

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
	reader            *bytes.Reader
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
	seg.reader = bytes.NewReader(c)

	/*
		seg.decoderId = uint16(c[0]) + (uint16(c[1]) << 8)
		seg.bytesPerBone = uint16(c[2]) + (uint16(c[3]) << 8)
		seg.boneCount = uint16(c[4]) + (uint16(c[5]) << 8)
		seg.dataLength = uint16(c[6]) + (uint16(c[7]) << 8)
	*/

	err := binary.Read(seg.reader, binary.LittleEndian, &seg.decoderId)
	if err != nil {
		return fmt.Errorf("failed to read segment decoder id: <%w>", err)
	}
	binary.Read(seg.reader, binary.LittleEndian, &seg.bytesPerBone)
	if err != nil {
		return fmt.Errorf("failed to read segment bytes per bone: <%w>", err)
	}
	binary.Read(seg.reader, binary.LittleEndian, &seg.boneCount)
	if err != nil {
		return fmt.Errorf("failed to read segment bone count: <%w>", err)
	}
	binary.Read(seg.reader, binary.LittleEndian, &seg.dataLength)
	if err != nil {
		return fmt.Errorf("failed to read segment data length: <%w>", err)
	}

	seg.bones = make([]uint16, seg.boneCount)
	for i := 0; i < int(seg.boneCount); i++ {
		err := binary.Read(seg.reader, binary.LittleEndian, &seg.bones[i])
		if err != nil {
			return fmt.Errorf("failed to read segment bone id: <%w>", err)
		}
	}

	log.Println(seg.bones)

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
