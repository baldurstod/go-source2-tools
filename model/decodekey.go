package model

import (
	"fmt"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type DecodeKey struct {
	Bones        []*DecodeKeyBone
	DataChannels []*DataChannel
}

func newDecodeKey() *DecodeKey {
	return &DecodeKey{
		Bones:        make([]*DecodeKeyBone, 0),
		DataChannels: make([]*DataChannel, 0),
	}
}

func (dk *DecodeKey) initFromDatas(datas *kv3.Kv3Element) error {
	boneArray, ok := datas.GetKv3ValueArrayAttribute("m_boneArray")
	if !ok {
		return fmt.Errorf("can't get bone array while initializing decode key")
	}

	for _, v := range boneArray {
		bone := newDecodeKeyBone()
		bone.initFromDatas(v.(*kv3.Kv3Element))
		dk.addBone(bone)
	}

	dataChannelArray, ok := datas.GetKv3ValueArrayAttribute("m_dataChannelArray")
	if !ok {
		return fmt.Errorf("can't get data channel array while initializing decode key")
	}

	for _, v := range dataChannelArray {
		channel := newDataChannel()
		channel.initFromDatas(v.(*kv3.Kv3Element))
		dk.addDataChannel(channel)
	}

	return nil
}

func (dk *DecodeKey) addBone(bone *DecodeKeyBone) {
	dk.Bones = append(dk.Bones, bone)
}

func (dk *DecodeKey) addDataChannel(dc *DataChannel) {
	dk.DataChannels = append(dk.DataChannels, dc)
}
