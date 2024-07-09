package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type Animation struct {
	group           *AnimGroup
	Name            string
	fps             float64
	FrameCount      int32
	frameblockArray []*frameBlock
}

func newAnimation(group *AnimGroup) *Animation {
	return &Animation{
		group:           group,
		frameblockArray: make([]*frameBlock, 0),
	}
}

func (anim *Animation) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	anim.Name, ok = datas.GetStringAttribute("m_name")
	if !ok {
		return errors.New("")
	}

	anim.fps, ok = datas.GetFloatAttribute("fps")
	if !ok {
		anim.fps = 30 //TODO: not sure if we should set a default value
	}

	pData := datas.GetKv3ElementAttribute("m_pData")
	if pData != nil {
		//log.Println(pData)
		anim.FrameCount, _ = pData.GetInt32Attribute("m_nFrames")
		frameblockArray, _ := pData.GetKv3ValueArrayAttribute("m_frameblockArray")
		for _, v := range frameblockArray {
			fb := new(frameBlock)
			fb.initFromDatas(v.(*kv3.Kv3Element))
		}
	}
	return nil
}

func (anim *Animation) GetFps() float64 {
	return anim.fps
}

func (anim *Animation) GetDuration() float64 {
	return float64(anim.FrameCount-1) / anim.fps
}

func (anim *Animation) GetFrame(frameIndex int) error {
	return nil
}
