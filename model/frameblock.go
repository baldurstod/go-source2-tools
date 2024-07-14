package model

import (
	"errors"
	"fmt"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type frameBlock struct {
	group        *AnimGroup
	block        *AnimBlock
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

func newFrameBlock(group *AnimGroup, block *AnimBlock) *frameBlock {
	return &frameBlock{
		group: group,
		block: block,
	}
}

func (fb *frameBlock) initFromDatas(datas *kv3.Kv3Element) error {
	startFrame, _ := datas.GetInt32Attribute("m_nStartFrame")
	fb.startFrame = int(startFrame)

	endFrame, _ := datas.GetInt32Attribute("m_nEndFrame")
	fb.endFrame = int(endFrame)

	segmentIndex, _ := datas.GetKv3ValueArrayAttribute("m_segmentIndexArray")
	fb.segmentIndex = make([]int, 0, len(segmentIndex))
	for _, v := range segmentIndex {
		fb.segmentIndex = append(fb.segmentIndex, kv3.Kv3ValueToInt(v))
	}
	return nil
}

func (fb *frameBlock) GetFrame(frameIndex int) (*frame, error) {
	frame := newFrame()
	frameIndex -= fb.startFrame
	//log.Println(fb.segmentIndex)
	for _, v := range fb.segmentIndex {
		seg := fb.block.getSegment(int(v))
		err := fb.readSegment(frameIndex, seg, frame)
		if err != nil {
			return nil, fmt.Errorf("error in frameBlock.GetFrame: <%w>", err)
		}
		//frame.addChannel(fc)
		//log.Println(seg)
	}
	return frame, nil
}

func (fb *frameBlock) readSegment(frameIndex int, segment *Segment, f *frame) error {
	decoder := &fb.block.decoders[segment.decoderId]
	//log.Println(decoder)
	channel := fb.group.decodeKey.getDataChannel(segment.LocalChannel)
	if channel == nil {
		return errors.New("can't find channel in readSegment")
	}

	fc := f.getChannel(channel.channelClass, channel.variableName, channel)

	err := segment.decode(frameIndex, channel, decoder, fc)
	if err != nil {
		return fmt.Errorf("error while reading segment: <%w>", err)
	}

	//log.Println(channel)

	return nil
}
