package model

import "github.com/baldurstod/go-source2-tools/kv3"

type Sequence struct {
	Name              string
	owner             *Model
	datas             *kv3.Kv3Element
	fps               float32
	LastFrame         int32
	Activity          string
	ActivityModifiers map[string]struct{}
}

func newSequence(owner *Model, datas *kv3.Kv3Element) *Sequence {
	seq := &Sequence{
		owner:             owner,
		datas:             datas,
		ActivityModifiers: make(map[string]struct{}),
	}

	var ok bool
	seq.Name, _ = datas.GetStringAttribute("m_name")
	seq.fps, ok = datas.GetFloat32Attribute("fps")
	if !ok {
		seq.fps = 30
	}

	return seq
}
