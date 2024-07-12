package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type DataChannel struct {
	channelClass  string
	variableName  string
	grouping      string
	description   string
	flags         int
	channelType   int32
	elementsName  []string
	elementsIndex []int32
	elementsMask  []uint32
}

func newDataChannel() *DataChannel {
	return &DataChannel{
		elementsName: make([]string, 0),
	}
}

func (dc *DataChannel) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	dc.channelClass, ok = datas.GetStringAttribute("m_szChannelClass")
	if !ok {
		return errors.New("unable to get data channel class")
	}

	dc.variableName, ok = datas.GetStringAttribute("m_szVariableName")
	if !ok {
		return errors.New("unable to get data channel variable name")
	}

	dc.flags, _ = datas.GetIntAttribute("m_nFlags")
	dc.channelType, _ = datas.GetInt32Attribute("m_nType")
	dc.grouping, _ = datas.GetStringAttribute("m_szGrouping")
	dc.description, _ = datas.GetStringAttribute("m_szDescription")

	elementName, _ := datas.GetKv3ValueArrayAttribute("m_szElementNameArray")
	dc.elementsName = make([]string, 0, len(elementName))
	for _, v := range elementName {
		dc.elementsName = append(dc.elementsName, v.(string))
	}

	elementIndex, _ := datas.GetKv3ValueArrayAttribute("m_nElementIndexArray")
	dc.elementsIndex = make([]int32, 0, len(elementIndex))
	for _, v := range elementIndex {
		dc.elementsIndex = append(dc.elementsIndex, v.(int32))
	}

	elementMask, _ := datas.GetKv3ValueArrayAttribute("m_nElementMaskArray")
	dc.elementsMask = make([]uint32, 0, len(elementMask))
	for _, v := range elementMask {
		dc.elementsMask = append(dc.elementsMask, v.(uint32))
	}

	return nil
}

/*
	{
		m_szChannelClass = "BoneChannel"
		m_szVariableName = "Position"
		m_nFlags = 0
		m_nType = 3
		m_szGrouping = ""
		m_szDescription = ""
		m_szElementNameArray =
		[
			"root",
			"head_child",
		]
		m_nElementIndexArray = [ 0, 1 ]
		m_nElementMaskArray = [ 1, 65536 ]
	},
*/
