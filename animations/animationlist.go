package animations

type AnimationList struct {
	Sequences map[*Sequence]struct{}
}

func (al AnimationList) AddSequence(seq *Sequence) {
	al.Sequences[seq] = struct{}{}
}
