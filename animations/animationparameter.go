package animations

import "math"

type AnimationParameter struct {
	Name   string
	weight float64
}

func (ap AnimationParameter) SetWeight(weight float64) {
	ap.weight = math.Max(math.Min(weight, 1.0), 0.0)
}
