package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
)

type DataChannelElement struct {
	name  string
	index int32
	mask  uint32
}

type DataChannel struct {
	channelClass string
	variableName string
	grouping     string
	description  string
	flags        int
	channelType  int32
	elements     []DataChannelElement
}

func newDataChannel() *DataChannel {
	return &DataChannel{
		elements: make([]DataChannelElement, 0),
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

	elementName, ok := datas.GetKv3ValueArrayAttribute("m_szElementNameArray")
	if !ok {
		return errors.New("unable to get data channel element name")
	}
	elementIndex, ok := datas.GetKv3ValueArrayAttribute("m_nElementIndexArray")
	if !ok {
		return errors.New("unable to get data channel element index")
	}
	elementMask, ok := datas.GetKv3ValueArrayAttribute("m_nElementMaskArray")
	if !ok {
		return errors.New("unable to get data channel element mask")
	}
	l := len(elementName)
	if len(elementIndex) != l {
		return errors.New("element name array is not the same length as element index")
	}
	if len(elementMask) != l {
		return errors.New("element name array is not the same length as element mask")
	}

	dc.elements = make([]DataChannelElement, len(elementName))
	for i := 0; i < l; i++ {
		elem := &dc.elements[i]
		elem.name = elementName[i].(string)
		elem.index = elementIndex[i].(int32)
		elem.mask = elementMask[i].(uint32)
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
