package particles

import (
	"errors"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/kv3"
)

type ParticleSystem struct {
	file *source2.File
	//ControlPoints []*ControlPoint
	*ControlPointConfigurations
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		//ControlPointConfigurations: NewControlPointConfigurations(),
	}
}

func (ps *ParticleSystem) SetFile(f *source2.File) {
	ps.file = f
}

/*
func (ps *ParticleSystem) GetControlPoints() ([]*ControlPoint, error) {
	if ps.ControlPoints != nil {
		return ps.ControlPoints, nil
	}

	s, err := ps.initControlPoints()
	if err != nil {
		return nil, err
	}
	ps.ControlPoints = s

	return s, nil
}*/

func (ps *ParticleSystem) GetControlPointConfiguration(name string) (*ControlPointConfiguration, error) {
	if ps.ControlPointConfigurations == nil {
		if err := ps.initControlPointConfigurations(); err != nil {
			return nil, err
		}
	}

	return ps.ControlPointConfigurations.GetConfigurationByName(name)
}

func (ps *ParticleSystem) GetControlPointConfigurationById(id int) (*ControlPointConfiguration, error) {
	if ps.ControlPointConfigurations == nil {
		if err := ps.initControlPointConfigurations(); err != nil {
			return nil, err
		}
	}

	return ps.ControlPointConfigurations.GetConfigurationById(id)
}

func (ps *ParticleSystem) initControlPointConfigurations() error {
	if ps.file == nil {
		return errors.New("no file provided")
	}

	ps.ControlPointConfigurations = NewControlPointConfigurations()

	cpConfig, err := ps.file.GetBlockStruct("DATA.m_controlPointConfigurations")
	if err != nil {
		return errors.New("can't find m_controlPointConfigurations attribute")
	}

	ps.ControlPointConfigurations.initFromDatas(cpConfig.([]kv3.Kv3Value))

	return nil
}
