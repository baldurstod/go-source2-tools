package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type frameBlock struct {
	startFrame   int
	endFrame     int
	segmentIndex []int

	/*
		{
			m_nStartFrame = 0
			m_nEndFrame = 14
			m_segmentIndexArray =
			[
				0, 74, 75, 76,
				77, 5, 6,
			]
		},*/
}

func (fb *frameBlock) initFromDatas(datas *kv3.Kv3Element) error {
	startFrame, _ := datas.GetInt32Attribute("m_nStartFrame")
	fb.startFrame = int(startFrame)

	endFrame, _ := datas.GetInt32Attribute("m_nEndFrame")
	fb.endFrame = int(endFrame)

	segmentIndex, _ := datas.GetKv3ValueArrayAttribute("m_segmentIndexArray")
	fb.segmentIndex = make([]int, 0, len(segmentIndex))
	for _, v := range segmentIndex {
		fb.segmentIndex = append(fb.segmentIndex, int(v.(int32)))
	}
	return nil
}
