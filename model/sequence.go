package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type IAnimationResource interface{}

type Sequence struct {
	Name              string
	owner             *Model
	datas             *kv3.Kv3Element
	FrameCount        uint32
	lastFrame         uint32
	Activity          string
	ActivityModifiers map[string]struct{}
	resource          IAnimationResource
}

func newSequence(owner *Model) *Sequence {
	return &Sequence{
		owner: owner,
		//datas:             datas,
		//resource:          resource,
		ActivityModifiers: make(map[string]struct{}),
	}
}

func (seq *Sequence) initFromDatas(datas *kv3.Kv3Element) error {
	seq.Name, _ = datas.GetStringAttribute("m_sName")

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
	return nil
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
	panic("todo")
	// 30 is the default when there is no underlying animation, for instance bind pose
	return 30
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

func (seq *Sequence) modifiersScore(modifiers map[string]struct{}) int {
	ret := 0

	if len(modifiers) == 0 && len(seq.ActivityModifiers) == 0 {
		return 1
	}

	for modifier := range modifiers {
		if _, ok := seq.ActivityModifiers[modifier]; ok {
			ret++
		}
	}

	return ret
}
