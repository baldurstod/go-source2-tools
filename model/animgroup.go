package model

import (
	"log"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type AnimGroup struct {
	model          *Model
	localAnimArray []kv3.Kv3Value
	decodeKey      *kv3.Kv3Element
	animations     map[*Animation]struct{}
}

func newAnimGroup(model *Model, localAnimArray []kv3.Kv3Value, decodeKey *kv3.Kv3Element) *AnimGroup {
	return &AnimGroup{
		model:          model,
		localAnimArray: localAnimArray,
		decodeKey:      decodeKey,
		animations:     make(map[*Animation]struct{}),
	}
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
