package model

type Bone struct {
	Name       string
	ParentBone *Bone
}

func NewBone(name string) *Bone {
	return &Bone{
		Name: name,
	}
}
