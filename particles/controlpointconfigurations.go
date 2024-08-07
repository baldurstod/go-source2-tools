package particles

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type ControlPointConfigurations struct {
	Configurations []*ControlPointConfiguration
}

func NewControlPointConfigurations() *ControlPointConfigurations {
	return &ControlPointConfigurations{
		Configurations: make([]*ControlPointConfiguration, 0),
	}

}

func (cpc *ControlPointConfigurations) initFromDatas(datas []kv3.Kv3Value) error {
	for _, c := range datas {
		e := kv3.Kv3ValueToKv3Element(c)
		if e != nil {
			config := NewControlPointConfiguration()
			config.initFromDatas(e)
			cpc.Configurations = append(cpc.Configurations, config)
		}
	}
	return nil
}

func (cpc *ControlPointConfigurations) GetConfigurationByName(name string) (*ControlPointConfiguration, error) {
	for _, config := range cpc.Configurations {
		if config.Name == name {
			return config, nil
		}
	}
	return nil, errors.New("config not found " + name)
}

func (cpc *ControlPointConfigurations) GetConfigurationById(id int) (*ControlPointConfiguration, error) {
	if id < 0 || id >= len(cpc.Configurations) {
		return nil, errors.New("control point configuration out of bounds")
	}
	return cpc.Configurations[id], nil
}
