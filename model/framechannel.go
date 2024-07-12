package model

type frameChannelData struct {
	name  string
	datas any
}

type frameChannel struct {
	channelClass string
	variableName string
	datas        []frameChannelData
}

func newFrameChannel(channelClass string, variableName string, dc *DataChannel) *frameChannel {
	fc := frameChannel{
		channelClass: channelClass,
		variableName: variableName,
		datas:        make([]frameChannelData, len(dc.elements)),
	}

	for k, v := range dc.elements {
		fc.datas[k].name = v.name
	}

	return &fc
}
