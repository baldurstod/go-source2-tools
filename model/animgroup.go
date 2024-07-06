package model

import "github.com/baldurstod/go-source2-tools/kv3"

type AnimGroup struct {
	model          *Model
	localAnimArray []kv3.Kv3Value
	decodeKey      *kv3.Kv3Element
}

func newAnimGroup(model *Model, localAnimArray []kv3.Kv3Value, decodeKey *kv3.Kv3Element) *AnimGroup {
	return &AnimGroup{
		model:          model,
		localAnimArray: localAnimArray,
		decodeKey:      decodeKey,
	}
}
