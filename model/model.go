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
	sequencesInitialized bool
	internalAnimGroup    *AnimGroup
	animGroups           map[*AnimGroup]struct{}
}

func NewModel() *Model {
	return &Model{
		sequences:  make(map[string]*Sequence),
		animGroups: make(map[*AnimGroup]struct{}),
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

func (m *Model) GetSequence(activity string, modifiers map[string]struct{}) error {
	if !m.sequencesInitialized {
		m.initSequences()
	}

	if m.file == nil {
		return errors.New("model don't have a file")
	}

	return nil
}

func (m *Model) GetAnimationData(animations []animations.AnimationParameter) error {
	for _, ap := range animations {
		log.Println(ap)
	}
	return nil
}

func (m *Model) initSequences() error {
	anim, err := m.file.GetBlockStructAsKv3Element("ANIM")
	if err != nil {
		return fmt.Errorf("can't find ANIM block: <%w>", err)
	}

	animArray, _ := anim.GetKv3ValueArrayAttribute("m_animArray")

	for _, v := range animArray {
		//log.Println(v)
		seq := newSequence(m, v.(*kv3.Kv3Element), nil)
		m.sequences[seq.Name] = seq
		//log.Println(seq.Name)
	}

	/*
		bonePosParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_bonePosParent")
		boneRotParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_boneRotParent")
		boneParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_nParent")
	*/

	log.Println(anim.GetAttributes())
	//log.Println(m.sequences)

	err = m.initInternalAnimGroup()
	if err != nil {
		return fmt.Errorf("error while initializing model sequences: <%w>", err)
	}

	return nil
}

func (m *Model) initInternalAnimGroup() error {
	localAnimArray, err := m.file.GetBlockStructAsKv3ValueArray("AGRP.m_localHAnimArray")
	if err != nil {
		return fmt.Errorf("can't get local anim array while initializing internal anim group: <%w>", err)
	}

	decodeKey, err := m.file.GetBlockStructAsKv3Element("AGRP.m_decodeKey")
	if err != nil {
		return fmt.Errorf("can't get decode key while initializing internal anim group: <%w>", err)
	}

	m.internalAnimGroup = newAnimGroup(m, localAnimArray, decodeKey)

	log.Println(localAnimArray, decodeKey, m.internalAnimGroup)

	anim, err := m.file.GetBlockStructAsKv3Element("ANIM")
	if err != nil {
		return fmt.Errorf("can't find ANIM block: <%w>", err)
	}

	loadedAnim := newAnimation(m.internalAnimGroup)
	//let anims = sourceFile.getBlockStruct('ANIM.keyValue.root');

	/*
		let anims = sourceFile.getBlockStruct('ANIM.keyValue.root');
		if (anims) {
			let loadedAnim = new Source2Animation(animGroup, '');
			loadedAnim.setAnimDatas(anims);
			animGroup._changemyname = animGroup._changemyname || [];
			animGroup._changemyname.push(loadedAnim);
		}
		this.animGroups.add(animGroup);
	*/
	log.Println(anim, loadedAnim)
	m.animGroups[m.internalAnimGroup] = struct{}{}

	return nil
}
