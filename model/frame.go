package model

type frame struct {
	channels map[string]map[string]*frameChannel
}

func newFrame() *frame {
	return &frame{
		channels: make(map[string]map[string]*frameChannel),
	}
}

func (f *frame) addChannel(fc *frameChannel) {
	m, ok := f.channels[fc.channelClass]
	if !ok {
		m = make(map[string]*frameChannel)
		f.channels[fc.channelClass] = m
	}

	m[fc.variableName] = fc
}
