package model

import (
	"fmt"
	"log"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/kv3"
)

type AnimGroup struct {
	model          *Model
	localAnimArray []kv3.Kv3Value
	decodeKey      *kv3.Kv3Element
	animations     map[*Animation]struct{}
}

func newAnimGroup(model *Model) *AnimGroup {
	return &AnimGroup{
		model:      model,
		animations: make(map[*Animation]struct{}),
	}
}

func (ag *AnimGroup) initFromFile(f *source2.File) error {
	var err error
	ag.localAnimArray, err = f.GetBlockStructAsKv3ValueArray("AGRP.m_localHAnimArray")
	if err != nil {
		return fmt.Errorf("can't get local anim array while initializing anim group: <%w>", err)
	}

	ag.decodeKey, err = f.GetBlockStructAsKv3Element("AGRP.m_decodeKey")
	if err != nil {
		return fmt.Errorf("can't get decode key while initializing anim group: <%w>", err)
	}
	return nil
}

func (ag *AnimGroup) CreateAnimation(block *AnimBlock) *Animation {
	anim := newAnimation(ag, block)
	ag.AddAnimation(anim)

	return anim
}

func (ag *AnimGroup) AddAnimation(anim *Animation) {
	ag.animations[anim] = struct{}{}
}

func (ag *AnimGroup) GetSequence(activity string, modifiers map[string]struct{}) error {
	for k := range ag.animations {
		log.Println(k)
	}

	return nil
}
