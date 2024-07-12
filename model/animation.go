package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type Animation struct {
	group       *AnimGroup
	block       *AnimBlock
	Name        string
	fps         float64
	FrameCount  int
	lastFrame   int
	frameBlocks []*frameBlock
}

func newAnimation(group *AnimGroup, block *AnimBlock) *Animation {
	return &Animation{
		group:       group,
		block:       block,
		frameBlocks: make([]*frameBlock, 0),
	}
}

func (anim *Animation) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	anim.Name, ok = datas.GetStringAttribute("m_name")
	if !ok {
		return errors.New("unable to get animation name")
	}

	anim.fps, ok = datas.GetFloat64Attribute("fps")
	if !ok {
		anim.fps = 30 //TODO: not sure if we should set a default value
	}

	pData := datas.GetKv3ElementAttribute("m_pData")
	if pData != nil {
		//log.Println(pData)
		anim.FrameCount, _ = pData.GetIntAttribute("m_nFrames")
		anim.lastFrame = anim.FrameCount - 1
		frameBlocks, _ := pData.GetKv3ValueArrayAttribute("m_frameblockArray")
		anim.frameBlocks = make([]*frameBlock, 0, len(frameBlocks))

		for _, v := range frameBlocks {
			fb := newFrameBlock(anim.group, anim.block)
			fb.initFromDatas(v.(*kv3.Kv3Element))
			anim.frameBlocks = append(anim.frameBlocks, fb)
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
	for _, fb := range anim.frameBlocks {
		if fb.startFrame <= frameIndex && fb.endFrame >= frameIndex {
			fb.GetFrame(frameIndex)
		}
	}

	//panic("todo")
	return nil
}
