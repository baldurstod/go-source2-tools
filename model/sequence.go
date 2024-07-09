package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type IAnimationResource interface{}

type Sequence struct {
	Name              string
	owner             *Model
	datas             *kv3.Kv3Element
	Activity          string
	ActivityModifiers map[string]struct{}
	animations        []string
	resource          IAnimationResource
}

func newSequence(owner *Model) *Sequence {
	return &Sequence{
		owner:      owner,
		animations: make([]string, 0, 3),
		//datas:             datas,
		//resource:          resource,
		ActivityModifiers: make(map[string]struct{}),
	}
}

func (seq *Sequence) initFromDatas(datas *kv3.Kv3Element, localSequenceNameArray []kv3.Kv3Value) error {
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

	fetch := datas.GetKv3ElementAttribute("m_fetch")
	if fetch != nil {
		localReferenceArray, _ := fetch.GetKv3ValueArrayAttribute("m_localReferenceArray")
		for _, v := range localReferenceArray {
			ref, ok := v.(int32)
			if !ok {
				continue
			}

			name := localSequenceNameArray[ref]
			n, ok := name.(string)
			if ok {
				seq.animations = append(seq.animations, n)
			}
		}
	}

	return nil
}

func (seq *Sequence) GetFps() float64 {
	for _, animName := range seq.animations {
		anim := seq.owner.animations[animName]
		if anim != nil {
			return anim.fps
		}
	}

	// 30 is the default when there is no underlying animation, for instance bind pose
	return 30
}

func (seq *Sequence) GetFrameCount() int {
	count := 0
	for _, animName := range seq.animations {
		anim := seq.owner.animations[animName]
		if anim != nil {
			count += int(anim.FrameCount)
		}
	}

	return count
}

func (seq *Sequence) GetFrame(frameIndex int) error {
	/*actualSequence := seq.GetActualSequence()
	if actualSequence != seq {
		return actualSequence.GetFrame(frameIndex)
	}*/

	if frameIndex > seq.GetFrameCount() {
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
