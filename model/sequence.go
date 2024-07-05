package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type IAnimationResource interface{}

type Sequence struct {
	Name              string
	owner             *Model
	datas             *kv3.Kv3Element
	fps               float32
	FrameCount        uint32
	lastFrame         uint32
	Activity          string
	ActivityModifiers map[string]struct{}
	frameblockArray   []kv3.Kv3Value
	resource          IAnimationResource
}

func newSequence(owner *Model, datas *kv3.Kv3Element, resource IAnimationResource) *Sequence {
	seq := &Sequence{
		owner:             owner,
		datas:             datas,
		resource:          resource,
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
		frameCount, _ := pData.GetInt32Attribute("m_nFrames")
		seq.FrameCount = uint32(frameCount)
		if frameCount > 0 {
			seq.lastFrame = seq.FrameCount - 1
		}
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

func (seq *Sequence) GetLastFrame() uint32 {
	return seq.GetActualSequence().lastFrame
}

func (seq *Sequence) GetFrame(frameIndex uint32) error {
	actualSequence := seq.GetActualSequence()
	if actualSequence != seq {
		return actualSequence.GetFrame(frameIndex)
	}

	if frameIndex > seq.lastFrame {
		frameIndex = 0
	}

	return nil
}
