package kv3

type Kv3Value interface {
}

func Kv3ValueToInt(value Kv3Value) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case uint32:
		return int(v)
	case int64:
		return int(v)
	default:
		panic("Unknown type in Kv3ValueToInt64")
	}
}

func Kv3ValueToInt64(value Kv3Value) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int64:
		return v
	default:
		panic("Unknown type in Kv3ValueToInt64")
	}
}

func Kv3ValueToKv3Element(value Kv3Value) *Kv3Element {
	v, ok := value.(*Kv3Element)
	if ok {
		return v
	}

	return nil
}
