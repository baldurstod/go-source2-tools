package model

type frameChannel struct {
	channelClass string
	variableName string
	datas        []any
}

func newFrameChannel(channelClass string, variableName string) *frameChannel {
	return &frameChannel{
		channelClass: channelClass,
		variableName: variableName,
		datas:        make([]any, 0),
	}
}
