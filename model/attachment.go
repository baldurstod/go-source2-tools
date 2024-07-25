package model

import (
	"errors"
	"strconv"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/baldurstod/go-vector"
)

type Influence struct {
	Name          string
	Offset        vector.Vector3[float32]
	Rotation      vector.Quaternion[float32]
	Weights       float32
	RootTransform bool
}

type Attachment struct {
	Name           string
	Influences     []Influence
	IgnoreRotation bool
}

func NewAttachment(name string) *Attachment {
	return &Attachment{
		Name:       name,
		Influences: make([]Influence, 0, 3),
	}
}

func (a *Attachment) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	a.Name, ok = datas.GetStringAttribute("m_name")
	if !ok {
		return errors.New("unable to get attachment name")
	}

	influencesCount, ok := datas.GetIntAttribute("m_nInfluences")
	if !ok {
		return errors.New("unable to get influences count")
	}

	influenceNames, _ := datas.GetKv3ValueArrayAttribute("m_influenceNames")
	influenceRotations, _ := datas.GetKv3ValueArrayAttribute("m_vInfluenceRotations")
	influenceOffsets, _ := datas.GetKv3ValueArrayAttribute("m_vInfluenceOffsets")
	influenceWeights, _ := datas.GetKv3ValueArrayAttribute("m_influenceWeights")
	influenceRootTransform, _ := datas.GetKv3ValueArrayAttribute("m_bInfluenceRootTransform")

	if len(influenceNames) < influencesCount ||
		len(influenceRotations) < influencesCount ||
		len(influenceOffsets) < influencesCount ||
		len(influenceWeights) < influencesCount ||
		len(influenceRootTransform) < influencesCount {
		return errors.New("wrong influence array size, must be " + strconv.Itoa(influencesCount))
	}

	for i := 0; i < influencesCount; i++ {
		influence := Influence{}

		influence.Name = kv3.Kv3ValueToString(influenceNames[i])
		influence.Offset = kv3.Kv3ValueToVector3[float32](influenceOffsets[i])
		influence.Rotation = kv3.Kv3ValueToQuaternion[float32](influenceRotations[i])
		influence.Weights = kv3.Kv3ValueToFloat32(influenceWeights[i])
		influence.RootTransform = kv3.Kv3ValueToBool(influenceRootTransform[i])

		a.Influences = append(a.Influences, influence)
	}

	return nil
}

/*
	key = "attach_weapon"
	value =
	{
		m_name = "attach_weapon"
		m_influenceNames =
		[
			"Staff_1",
			"",
			"",
		]
		m_vInfluenceRotations =
		[
			[ 0.0, -0.0, -0.0, 1.0 ],
			[ 0.0, 0.0, 0.0, 1.0 ],
			[ 0.0, 0.0, 0.0, 1.0 ],
		]
		m_vInfluenceOffsets =
		[
			[ 0.0, 0.0, 0.0 ],
			[ 0.0, 0.0, 0.0 ],
			[ 0.0, 0.0, 0.0 ],
		]
		m_influenceWeights = [ 1.0, 0.0, 0.0 ]
		m_bInfluenceRootTransform = [ false, false, false ]
		m_nInfluences = 1
		m_bIgnoreRotation = false
	}*/
