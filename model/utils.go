package model

import (
	"errors"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/baldurstod/go-vector"
)

func KvArrayToVector3(src []kv3.Kv3Value, dst *vector.Vector3[float32]) error {
	if len(src) != 3 {
		return errors.New("length of vector3 must be 3")
	}

	for i := 0; i < 3; i++ {
		switch v := src[i].(type) {
		case float32:
			dst[i] = v
		case float64:
			dst[i] = float32(v)
		default:
			return errors.New("unknown value type")
		}
	}
	return nil
}
