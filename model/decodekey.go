package model

type DecodeKey struct {
	Bones        []*Bone
	DataChannels []*DataChannel
}

func newDecodeKey() *DecodeKey {
	return &DecodeKey{
		Bones:        make([]*Bone, 0),
		DataChannels: make([]*DataChannel, 0),
	}
}

func (dk *DecodeKey) AddBone(bone *Bone) {
	dk.Bones = append(dk.Bones, bone)
}

func (dk *DecodeKey) AddDataChannel(dc *DataChannel) {
	dk.DataChannels = append(dk.DataChannels, dc)
}
