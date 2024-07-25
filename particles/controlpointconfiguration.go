package particles

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type ControlPointConfiguration struct {
	Name    string
	Drivers []*Driver
	*PreviewState
}

func NewControlPointConfiguration() *ControlPointConfiguration {
	return &ControlPointConfiguration{
		Drivers: make([]*Driver, 0),
	}
}

func (cpc *ControlPointConfiguration) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	cpc.Name, ok = datas.GetStringAttribute("m_name")
	if !ok {
		return errors.New("can't find name attribute")
	}

	drivers, _ := datas.GetKv3ValueArrayAttribute("m_drivers")
	for _, d := range drivers {
		e := kv3.Kv3ValueToKv3Element(d)
		if e != nil {
			driver := Driver{}
			driver.initFromDatas(e)
			cpc.Drivers = append(cpc.Drivers, &driver)
		}
	}

	previewState := datas.GetKv3ElementAttribute("m_previewState")
	if previewState != nil {
		cpc.PreviewState = &PreviewState{}
		cpc.PreviewState.initFromDatas(previewState)
	}

	return nil
}
