package choreography

type Choreography struct {
	Name          string
	Duration      uint32
	SoundDuration uint32
	HasSounds     bool
	events        []*ChoreographyEvent
	actors        []*ChoreographyActor
}

func NewChoreography() *Choreography {
	return &Choreography{
		events: make([]*ChoreographyEvent, 0, 10),
		actors: make([]*ChoreographyActor, 0, 10),
	}
}

func (c *Choreography) AddEvent(event *ChoreographyEvent) {
	event.Choreography = c
	c.events = append(c.events, event)
}

func (c *Choreography) AddActor(actor *ChoreographyActor) {
	actor.Choreography = c
	c.actors = append(c.actors, actor)
}

type ChoreographyEvent struct {
	Choreography *Choreography
	EventType    int8
	Name         string
	StartTime    float32
	EndTime      float32
	Param1       string
	Param2       string
	Param3       string
	*CurveData
	Flags        uint8
	DistToTarget float32
}

func NewChoreographyEvent() *ChoreographyEvent {
	return &ChoreographyEvent{}
}

type CurveData struct {
	Ramp []*CurveDataValue
}

func NewCurveData() *CurveData {
	return &CurveData{
		Ramp: make([]*CurveDataValue, 0, 10),
	}
}

func (c *CurveData) Add(time float32, value float32, selected bool) {
	c.Ramp = append(c.Ramp, &CurveDataValue{
		Time:     time,
		Value:    value,
		Selected: selected,
	})
}

type CurveDataValue struct {
	Time     float32
	Value    float32
	Selected bool
}

type ChoreographyActor struct {
	Choreography *Choreography
	Name         string
	Active       bool
	Channels     []*ChoreographyChannel
}

func NewChoreographyActor() *ChoreographyActor {
	return &ChoreographyActor{}
}

func (a *ChoreographyActor) AddChannel(channel *ChoreographyChannel) {
	channel.Choreography = a.Choreography
	a.Channels = append(a.Channels, channel)
}

type ChoreographyChannel struct {
	Choreography *Choreography
	Name         string
	Active       bool
	events       []*ChoreographyEvent
}

func NewChoreographyChannel() *ChoreographyChannel {
	return &ChoreographyChannel{
		events: make([]*ChoreographyEvent, 0, 10),
	}
}

func (c *ChoreographyChannel) AddEvent(event *ChoreographyEvent) {
	event.Choreography = c.Choreography
	c.events = append(c.events, event)
}
