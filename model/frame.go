package model

type frame struct {
	channels map[string]map[string]*frameChannel
}

func newFrame() *frame {
	return &frame{
		channels: make(map[string]map[string]*frameChannel),
	}
}

func (f *frame) getChannel(channelClass string, variableName string, dc *DataChannel) *frameChannel {
	m, ok := f.channels[channelClass]
	if !ok {
		m = make(map[string]*frameChannel)
		f.channels[channelClass] = m
	}

	fc, ok := m[variableName]
	if !ok {
		fc = newFrameChannel(channelClass, variableName, dc)
		m[variableName] = fc
	}

	return fc
}
