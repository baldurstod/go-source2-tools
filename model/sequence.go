package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type Sequence struct {
	Name              string
	owner             *Model
	datas             *kv3.Kv3Element
	fps               float32
	FrameCount        int32
	lastFrame         int32
	Activity          string
	ActivityModifiers map[string]struct{}
	frameblockArray   []kv3.Kv3Value
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
		seq.fps = 30 //TODO: not sure if we should set a default value
	}

	pData := datas.GetKv3ElementAttribute("m_pData")
	if pData != nil {
		//log.Println(pData)
		seq.FrameCount, _ = pData.GetInt32Attribute("m_nFrames")
		seq.frameblockArray, _ = pData.GetKv3ValueArrayAttribute("m_frameblockArray")

		///animArray, _ := anim.GetKv3ValueArrayAttribute("m_animArray")
	}

	activityArray, _ := datas.GetKv3ValueArrayAttribute("m_activityArray")
	for k, v := range activityArray {
		activity, ok := v.(*kv3.Kv3Element)
		if !ok {
			continue
		}

		name, _ := activity.GetStringAttribute("m_name")
		if k == 0 {
			seq.Activity = name
		} else {
			seq.ActivityModifiers[name] = struct{}{}
		}
	}

	return seq
}

func (seq *Sequence) GetActualSequence() *Sequence {
	fetch := seq.datas.GetKv3ElementAttribute("m_fetch")
	if fetch == nil {
		return seq
	}

	//localReferenceArray := fetch.GetKv3ElementAttribute("m_localReferenceArray")
	//log.Println(localReferenceArray)
	panic("code me")

	return seq
}

func (seq *Sequence) GetFps() float32 {
	return seq.GetActualSequence().fps
}

func (seq *Sequence) GetLastFrame() int32 {
	return seq.GetActualSequence().lastFrame
}
