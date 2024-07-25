package kv3

import "github.com/baldurstod/go-vector"

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

func Kv3ValueToFloat32(value Kv3Value) float32 {
	switch v := value.(type) {
	case float32:
		return v
	case float64:
		return float32(v)
	}
	return 0
}

func Kv3ValueToFloat64(value Kv3Value) float64 {
	switch v := value.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	}
	return 0
}

func Kv3ValueToBool(value Kv3Value) bool { /*
		switch v := value.(type) {
		case float32:
			return float64(v)
		case float64:
			return v
		}*/
	return false
}

func Kv3ValueToKv3Element(value Kv3Value) *Kv3Element {
	v, ok := value.(*Kv3Element)
	if ok {
		return v
	}

	return nil
}

func Kv3ValueToString(value Kv3Value) string {
	v, ok := value.(string)
	if ok {
		return v
	}

	return ""
}

func Kv3ValueToVector3[T float32 | float64](value Kv3Value) vector.Vector3[T] {
	a, ok := value.([]T)
	if ok && len(a) == 3 {
		return vector.Vector3[T]{a[0], a[1], a[2]}
	}

	return vector.Vector3[T]{}
}

func Kv3ValueToQuaternion[T float32 | float64](value Kv3Value) vector.Quaternion[T] {
	a, ok := value.([]T)
	if ok && len(a) == 4 {
		return vector.Quaternion[T]{a[0], a[1], a[2], a[3]}
	}

	return vector.Quaternion[T]{}
}
