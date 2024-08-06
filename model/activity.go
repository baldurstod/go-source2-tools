package model

type Activity struct {
	name      string
	modifiers map[string]struct{}
}

func NewActivity(activity string, modifiers ...string) *Activity {
	a := Activity{
		name:      activity,
		modifiers: make(map[string]struct{}),
	}

	for _, modifier := range modifiers {
		a.modifiers[modifier] = struct{}{}
	}

	return &a
}
