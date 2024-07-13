package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/animations"
	"github.com/baldurstod/go-source2-tools/kv3"
)

type Model struct {
	file                 *source2.File
	skeleton             *Skeleton
	sequences            map[string]*Sequence
	animations           map[string]*Animation
	activities           map[string]map[*Sequence]struct{}
	sequencesInitialized bool
	animGroups           map[*AnimGroup]struct{}
	animBlock            *AnimBlock
}

func NewModel() *Model {
	return &Model{
		sequences:  map[string]*Sequence{},
		animations: map[string]*Animation{},
		activities: make(map[string]map[*Sequence]struct{}),
		animGroups: make(map[*AnimGroup]struct{}),
		animBlock:  newAnimBlock(),
	}
}

func (m *Model) SetFile(f *source2.File) {
	m.file = f
}

func (m *Model) GetSkeleton() (*Skeleton, error) {
	if m.skeleton != nil {
		return m.skeleton, nil
	}

	s, err := m.initSkeleton()
	if err != nil {
		return nil, err
	}
	m.skeleton = s

	return s, nil
}

func (m *Model) initSkeleton() (*Skeleton, error) {
	if m.file == nil {
		return nil, errors.New("model don't have a file")
	}

	skeleton, err := m.file.GetBlockStruct("DATA.m_modelSkeleton")
	if err != nil {
		return nil, errors.New("can't find m_modelSkeleton attribute")
	}

	boneNames := skeleton.(*kv3.Kv3Element).GetAttribute("m_boneName")
	bonePosParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_bonePosParent")
	boneRotParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_boneRotParent")
	boneParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_nParent")

	if boneNames == nil || bonePosParents == nil || boneRotParents == nil || boneParents == nil {
		return nil, errors.New("can't find m_modelSkeleton sub attributes")
	}

	var boneNamesArray []kv3.Kv3Value
	var bonePosParentArray []kv3.Kv3Value
	var boneRotParentArray []kv3.Kv3Value
	var boneParentArray []kv3.Kv3Value
	var ok bool

	if boneNamesArray, ok = boneNames.([]kv3.Kv3Value); !ok {
		return nil, errors.New("m_boneName is not an array")
	}
	if bonePosParentArray, ok = bonePosParents.([]kv3.Kv3Value); !ok {
		return nil, errors.New("m_bonePosParent is not an array")
	}
	if boneRotParentArray, ok = boneRotParents.([]kv3.Kv3Value); !ok {
		return nil, errors.New("m_boneRotParent is not an array")
	}
	if boneParentArray, ok = boneParents.([]kv3.Kv3Value); !ok {
		return nil, errors.New("m_nParent is not an array")
	}

	if len(boneNamesArray) != len(bonePosParentArray) ||
		len(boneNamesArray) != len(boneRotParentArray) ||
		len(boneNamesArray) != len(boneParentArray) {
		return nil, errors.New("bone arrays have different sizes")
	}

	s := NewSkeleton(len(boneNamesArray))
	var boneName string
	var bonePosParent []kv3.Kv3Value
	var boneRotParent []kv3.Kv3Value
	for i := 0; i < len(boneNamesArray); i++ {
		if boneName, ok = boneNamesArray[i].(string); !ok {
			return nil, errors.New("bone name is not a string")
		}
		if bonePosParent, ok = bonePosParentArray[i].([]kv3.Kv3Value); !ok {
			return nil, errors.New("bone pos is not an array")
		}
		if boneRotParent, ok = boneRotParentArray[i].([]kv3.Kv3Value); !ok {
			return nil, errors.New("bone rot is not an array")
		}

		bone, _ := s.CreateBone(boneName)
		if err := KvArrayToVector3(bonePosParent, &bone.PosParent); err != nil {
			return nil, fmt.Errorf("error while decoding bone parent position: <%w>", err)
		}
		if err := KvArrayToQuaternion(boneRotParent, &bone.RotParent); err != nil {
			return nil, fmt.Errorf("error while decoding bone parent rotation: <%w>", err)
		}
	}

	// Phase 2: parenting
	for i := 0; i < len(boneNamesArray); i++ {
		bone, err := s.GetBoneById(i)
		if err != nil {
			return nil, fmt.Errorf("can't get bone id: %d <%w>", i, err)
		}

		p := boneParentArray[i].(int32)
		if p != -1 {
			parentBone, err := s.GetBoneById(int(p))
			if err != nil {
				return nil, fmt.Errorf("can't get parent bone id: %d <%w>", p, err)
			}

			bone.ParentBone = parentBone
		}

	}

	return s, nil
}

func (m *Model) GetSequence(activity string, modifiers map[string]struct{}) (*Sequence, error) {
	if !m.sequencesInitialized {
		err := m.initAnimations()
		if err != nil {
			return nil, fmt.Errorf("error in getsequence: <%w>", err)
		}
	}

	if m.file == nil {
		return nil, errors.New("model don't have a file")
	}

	activities, ok := m.activities[activity]
	if !ok {
		return nil, errors.New("activity not found " + activity)
	}

	bestScore := -1
	var bestMatch *Sequence
	//TODO: weight similar animations
	for sequence := range activities {
		score := sequence.modifiersScore(modifiers)
		if score > bestScore {
			bestScore = score
			bestMatch = sequence
		}
	}
	/*
		for group := range m.animGroups {
			return group.GetSequence(activity, modifiers)
		}
	*/

	return bestMatch, nil
}

func (m *Model) GetSequenceByName(name string) (*Sequence, error) {
	if !m.sequencesInitialized {
		err := m.initAnimations()
		if err != nil {
			return nil, fmt.Errorf("error in getsequence: <%w>", err)
		}
	}

	sequence, ok := m.sequences[name]
	if !ok {
		return nil, errors.New("sequence not found " + name)
	}

	return sequence, nil
}

func (m *Model) GetAnimationData(animations []animations.AnimationParameter) error {
	if !m.sequencesInitialized {
		err := m.initAnimations()
		if err != nil {
			return fmt.Errorf("error in GetAnimationData: <%w>", err)
		}
	}
	for _, ap := range animations {
		log.Println(ap)
	}
	return nil
}

func (m *Model) initAnimations() error {
	err := m.initSequences()
	if err != nil {
		return err
	}

	err = m.initAnims()
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) initSequences() error {
	sequences, err := m.file.GetBlockStructAsKv3Element("ASEQ")
	if err != nil {
		return fmt.Errorf("can't find ASEQ block: <%w>", err)
	}

	localS1SeqDescArray, ok := sequences.GetKv3ValueArrayAttribute("m_localS1SeqDescArray")
	if !ok {
		return errors.New("key not found m_localS1SeqDescArray")
	}

	localSequenceNameArray, ok := sequences.GetKv3ValueArrayAttribute("m_localSequenceNameArray")
	if !ok {
		return errors.New("key not found m_localSequenceNameArray")
	}
	for _, v := range localS1SeqDescArray {
		seq := newSequence(m)
		seq.initFromDatas(v.(*kv3.Kv3Element), localSequenceNameArray)
		m.addSequence(seq)
	}

	return nil
}

func (m *Model) addSequence(seq *Sequence) {
	a, ok := m.activities[seq.Activity]

	if !ok {
		a = make(map[*Sequence]struct{})
		m.activities[seq.Activity] = a
	}
	a[seq] = struct{}{}
	m.sequences[seq.Name] = seq
}

func (m *Model) PrintSequences() {
	for k, v := range m.activities {
		log.Println(k)
		for k2 := range v {
			log.Println("\t" + k2.Name)
			for k3 := range k2.ActivityModifiers {
				log.Println("\t\t" + k3)
			}
		}
	}
}

func (m *Model) initAnims() error {
	anim, err := m.file.GetBlockStructAsKv3Element("ANIM")
	if err != nil {
		return fmt.Errorf("can't find ANIM block: <%w>", err)
	}
	//log.Println(anim.GetAttributes())

	err = m.animBlock.initFromDatas(anim)
	if err != nil {
		return fmt.Errorf("error while initializing anim block datas: <%w>", err)
	}

	err = m.initInternalAnimGroup()
	if err != nil {
		return fmt.Errorf("error while initializing model sequences: <%w>", err)
	}

	return nil
}

func (m *Model) initInternalAnimGroup() error {
	animGroup := newAnimGroup(m)
	animGroup.initFromFile(m.file)

	animDatas, err := m.file.GetBlockStructAsKv3Element("ANIM")
	if err != nil {
		return fmt.Errorf("can't find ANIM block: <%w>", err)
	}

	animArray, ok := animDatas.GetKv3ValueArrayAttribute("m_animArray")

	if ok {
		for _, v := range animArray {
			anim := animGroup.CreateAnimation(m.animBlock)
			anim.initFromDatas(v.(*kv3.Kv3Element))
			m.animations[anim.Name] = anim
		}
	}
	m.animGroups[animGroup] = struct{}{}

	return nil
}
