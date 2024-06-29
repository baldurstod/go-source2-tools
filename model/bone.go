package model

import (
	"github.com/baldurstod/go-vector"
)

type Bone struct {
	Name        string
	ParentBone  *Bone
	PosParent   vector.Vector3[float32]
	RotParent   vector.Quaternion[float32]
	ScaleParent float32
}

func NewBone(name string) *Bone {
	return &Bone{
		Name: name,
	}
}

func (b *Bone) String() string {
	var s = "{"

	s += b.Name

	if b.ParentBone != nil {
		s += " parent: "
		s += b.ParentBone.Name
	}
	s += b.PosParent.String()
	s += b.RotParent.String()

	s += "}"
	return s
}
