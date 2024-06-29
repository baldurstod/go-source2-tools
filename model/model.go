package model

import (
	"errors"
	"log"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/kv3"
)

type Model struct {
	file     *source2.File
	skeleton *Skeleton
}

func NewModel() *Model {
	return &Model{}
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
	log.Println(skeleton, err)
	log.Println(skeleton.(*kv3.Kv3Element).GetAttribute("m_boneName"), err)

	boneNames := skeleton.(*kv3.Kv3Element).GetAttribute("m_boneName")
	bonePosParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_bonePosParent")
	boneRotParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_boneRotParent")
	boneParents := skeleton.(*kv3.Kv3Element).GetAttribute("m_nParent")

	if boneNames == nil || bonePosParents == nil || boneRotParents == nil || boneParents == nil {
		return nil, errors.New("can't find m_modelSkeleton attributes")
	}

	log.Println(boneNames)

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
	for i := 0; i < len(boneNamesArray); i++ {
		if boneName, ok = boneNamesArray[i].(string); !ok {
			return nil, errors.New("bone name is not a string")
		}
		if bonePosParent, ok = bonePosParentArray[i].([]kv3.Kv3Value); !ok {
			return nil, errors.New("bone pos is not an array")
		}

		bone, _ := s.CreateBone(boneName)
		err = KvArrayToVector3(bonePosParent, &bone.PosParent)

	}

	return s, nil
}
