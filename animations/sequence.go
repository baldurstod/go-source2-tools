package animations

type Sequence struct {
	Name              string
	Fps               float32
	Activity          string
	ActivityModifiers map[string]struct{}
}

func NewSequence() *Sequence {
	return &Sequence{
		ActivityModifiers: make(map[string]struct{}),
	}
}
