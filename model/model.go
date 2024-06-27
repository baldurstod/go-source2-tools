package model

import (
	"errors"
	"log"

	"github.com/baldurstod/go-source2-tools"
)

type Model struct {
	file *source2.File
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) SetFile(f *source2.File) {
	m.file = f
}

func (m *Model) GetBones() ([]*Bone, error) {
	if m.file == nil {
		return nil, errors.New("model don't have a file")
	}

	skeleton, err := m.file.GetBlockStruct("DATA.m_modelSkeleton")
	log.Println(skeleton, err)
	return nil, nil
}
