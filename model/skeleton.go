package model

import "errors"

type Skeleton struct {
	bones []*Bone
	names map[string]*Bone
}

func NewSkeleton(bones int) *Skeleton {
	return &Skeleton{
		bones: make([]*Bone, 0, bones),
		names: make(map[string]*Bone),
	}
}

func (s *Skeleton) AddBone(bone *Bone) error {
	s.bones = append(s.bones, bone)

	s.names[bone.Name] = bone

	return nil
}

func (s *Skeleton) GetBoneById(id int) (*Bone, error) {
	if id >= len(s.bones) {
		return nil, errors.New("bone id is out of bounds")
	}

	return s.bones[id], nil
}

func (s *Skeleton) GetBoneByName(name string) (*Bone, error) {
	bone, ok := s.names[name]

	if !ok {
		return nil, errors.New("unknown bone name " + name)
	}

	return bone, nil
}
