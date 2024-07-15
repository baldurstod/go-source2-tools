package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type FlexController struct {
	Name     string
	FlexType string
	Min      float32
	Max      float32
}

func (fc *FlexController) initFromDatas(datas *kv3.Kv3Element) error {

	var ok bool
	fc.Name, ok = datas.GetStringAttribute("m_szName")
	if !ok {
		return errors.New("unable to get flex controller name")
	}

	fc.FlexType, ok = datas.GetStringAttribute("m_szType")
	if !ok {
		fc.FlexType = "default"
	}

	fc.Min, _ = datas.GetFloat32Attribute("min")
	fc.Max, _ = datas.GetFloat32Attribute("max")

	return nil
}

// Return the root of a linear interpolation between min and max
func (fc *FlexController) GetDefaultValue() float32 {
	if fc.Max == fc.Min {
		return 0
	}
	return -fc.Min / (fc.Max - fc.Min)
}

/*
	{
		m_szName = "innerBrowRaiser"
		m_szType = "default"
		min = 0.0
		max = 1.0
	},
*/
