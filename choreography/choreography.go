package choreography

type Choreography struct {
	Name          string
	Duration      uint32
	SoundDuration uint32
	HasSounds     bool
	events        []*ChoreographyEvent
}

func NewChoreography() *Choreography {
	return &Choreography{
		events: make([]*ChoreographyEvent, 10),
	}
}

func (c *Choreography) AddEvent(event *ChoreographyEvent) {
	event.Choreography = c
	c.events = append(c.events, event)
}

type ChoreographyEvent struct {
	*Choreography
	EventType int8
	Name      string
	StartTime float32
	EndTime   float32
	Param1    string
	Param2    string
	Param3    string
	*CurveData
	Flags        uint8
	DistToTarget float32
	/*
	   this.#repository = repository;
	   this.type = eventType;
	   this.name = name;
	   this.startTime = startTime;
	   this.endTime = endTime;
	   this.param1 = param1;
	   this.param2 = param2;
	   this.param3 = param3;
	   this.flags = flags;
	   this.distanceToTarget = distanceToTarget;
	   this.flexAnimTracks = {};
	*/
}

func NewChoreographyEvent() *ChoreographyEvent {
	return &ChoreographyEvent{
		//events: make([]*ChoreographyEvent, 10),
	}
}

type CurveData struct {
	Ramp []*CurveDataValue
}

func NewCurveData() *CurveData {
	return &CurveData{
		Ramp: make([]*CurveDataValue, 10),
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
