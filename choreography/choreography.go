package choreography

type Choreography struct {
	Name          string
	Duration      uint32
	SoundDuration uint32
	HasSounds     bool
	Events        []*ChoreographyEvent
	Actors        []*ChoreographyActor
}

func NewChoreography() *Choreography {
	return &Choreography{
		Events: make([]*ChoreographyEvent, 0, 10),
		Actors: make([]*ChoreographyActor, 0, 10),
	}
}

func (c *Choreography) AddEvent(event *ChoreographyEvent) {
	event.choreography = c
	c.Events = append(c.Events, event)
}

func (c *Choreography) AddActor(actor *ChoreographyActor) {
	actor.choreography = c
	c.Actors = append(c.Actors, actor)
}

type ChoreographyEvent struct {
	choreography  *Choreography
	EventType     EventType
	Name          string
	PreferredName string
	StartTime     float32
	EndTime       float32
	Param1        string
	Param2        string
	Param3        string
	*CurveData
	Flags              uint8
	DistToTarget       float32
	RelativeTags       []*ChoreographyTag
	FlexTimingTags     []*ChoreographyTag
	AbsoluteTags       [2][]*ChoreographyTag
	UsesTag            bool
	TagName            string
	TagWavName         string
	Tracks             map[string]*FlexAnimationTrack
	LoopCount          int
	CloseCaptionType   int8
	CloseCaptionToken  string
	SpeakFlags         uint8
	SoundStartDelay    float32
	ConstrainedEventId int32
	EventId            int32
}

func NewChoreographyEvent() *ChoreographyEvent {
	return &ChoreographyEvent{
		Tracks: make(map[string]*FlexAnimationTrack),
	}
}

func (ce *ChoreographyEvent) AddTrack(name string) *FlexAnimationTrack {
	track := &FlexAnimationTrack{}
	ce.Tracks[name] = track
	return track
}

type CurveData struct {
	Ramp      []*CurveDataSample
	LeftEdge  *CurveDataEdge
	RightEdge *CurveDataEdge
}

func NewCurveData() *CurveData {
	return &CurveData{
		Ramp: make([]*CurveDataSample, 0, 10),
	}
}

func (c *CurveData) AddSample(sample *CurveDataSample) {
	c.Ramp = append(c.Ramp, sample)
}

type CurveDataSample struct {
	Time      float32
	Value     float32
	Selected  bool
	Bezier    *CurveDataSampleBezier
	CurveType *CurveDataSampleType
}

type CurveDataSampleBezier struct {
	Flags     uint8
	InDeg     float32
	InWeight  float32
	OutDeg    float32
	OutWeight float32
}

type CurveDataSampleType struct {
	InType  uint8
	OutType uint8
}

type CurveDataEdge struct {
	CurveType CurveDataSampleType
	ZeroValue float32
}

type ChoreographyActor struct {
	choreography *Choreography
	Name         string
	Active       bool
	Channels     []*ChoreographyChannel
}

func NewChoreographyActor() *ChoreographyActor {
	return &ChoreographyActor{}
}

func (a *ChoreographyActor) AddChannel(channel *ChoreographyChannel) {
	channel.choreography = a.choreography
	a.Channels = append(a.Channels, channel)
}

type ChoreographyChannel struct {
	choreography *Choreography
	Name         string
	Active       bool
	Events       []*ChoreographyEvent
}

func NewChoreographyChannel() *ChoreographyChannel {
	return &ChoreographyChannel{
		Events: make([]*ChoreographyEvent, 0, 10),
	}
}

func (c *ChoreographyChannel) AddEvent(event *ChoreographyEvent) {
	event.choreography = c.choreography
	c.Events = append(c.Events, event)
}

type ChoreographyTag struct {
	choreography *Choreography
	Name         string
	Value        float32
}

type FlexAnimationTrack struct {
	Flags             uint8
	Min               float32
	Max               float32
	SamplesCurve      *CurveData
	ComboSamplesCurve *CurveData
}

func (fat *FlexAnimationTrack) IsTrackActive() bool {
	return (fat.Flags & (1 << 0)) == 1
}

func (fat *FlexAnimationTrack) IsComboType() bool {
	return (fat.Flags & (1 << 1)) == 1<<1
}
