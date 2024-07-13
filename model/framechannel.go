package model

type frameChannelData struct {
	Name  string
	Datas any
}

type frameChannel struct {
	ChannelClass string
	VariableName string
	Datas        []frameChannelData
}

func newFrameChannel(channelClass string, variableName string, dc *DataChannel) *frameChannel {
	fc := frameChannel{
		ChannelClass: channelClass,
		VariableName: variableName,
		Datas:        make([]frameChannelData, len(dc.elements)),
	}

	for k, v := range dc.elements {
		fc.Datas[k].Name = v.name
	}

	return &fc
}
