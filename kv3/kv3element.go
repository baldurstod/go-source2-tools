package kv3

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/baldurstod/go-vector"
)

type Kv3Element struct {
	attributes map[string]any
}

func NewKv3Element() *Kv3Element {
	return &Kv3Element{
		attributes: make(map[string]any),
	}
}

func (e *Kv3Element) AddAttribute(name string, value any) {
	e.attributes[name] = value
}

func (e *Kv3Element) GetAttribute(name string) any {
	return e.attributes[name]
}

func (e *Kv3Element) GetAttributes() []string {
	attributes := make([]string, 0, len(e.attributes))

	for k := range e.attributes {
		attributes = append(attributes, k)
	}

	return attributes
}

func (e *Kv3Element) GetStringAttribute(name string) (string, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return "", false
	}

	f, ok := value.(string)
	if !ok {
		return "", false
	}

	return f, true
}

func (e *Kv3Element) GetIntAttribute(name string) (int, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	default:
		return 0, false
	}
}

func (e *Kv3Element) GetInt32Attribute(name string) (int32, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return int32(v), true
	case int32:
		return v, true
	default:
		return 0, false
	}
}

func (e *Kv3Element) GetFloat32Attribute(name string) (float32, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case float32:
		return v, true
	case float64:
		return float32(v), true
	default:
		return 0, false
	}
}

func (e *Kv3Element) GetFloat64Attribute(name string) (float64, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func (e *Kv3Element) GetKv3ValueArrayAttribute(name string) ([]Kv3Value, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return make([]Kv3Value, 0), false
	}

	v, ok := value.([]Kv3Value)
	if !ok {
		return make([]Kv3Value, 0), false
	}

	return v, true
}

func (e *Kv3Element) GetKv3ElementAttribute(name string) *Kv3Element {
	value, ok := e.attributes[name]
	if !ok {
		return nil
	}

	v, ok := value.(*Kv3Element)
	if !ok {
		return nil
	}

	return v
}

func (e *Kv3Element) GetVec3Attribute(name string) (*vector.Vector3[float32], error) {
	value, ok := e.GetKv3ValueArrayAttribute(name)
	if !ok {
		return nil, errors.New("can't get vec3 attribute")
	}

	if len(value) != 3 {
		return nil, errors.New("vec3 length must be 3")
	}

	var vec vector.Vector3[float32]
	for i := 0; i < 3; i++ {
		v := value[i]

		switch v := v.(type) {
		case float32:
			vec[i] = v
		case float64:
			vec[i] = float32(v)
		default:
			return nil, errors.New("vec3 item is not a float")
		}

	}
	return &vec, nil
}

func (e *Kv3Element) GetQuatAttribute(name string) (*vector.Quaternion[float32], error) {
	value, ok := e.GetKv3ValueArrayAttribute(name)
	if !ok {
		return nil, errors.New("can't get quaternion attribute")
	}

	if len(value) != 4 {
		return nil, errors.New("quaternion length must be 4")
	}

	var quat vector.Quaternion[float32]
	for i := 0; i < 4; i++ {
		v := value[i]

		switch v := v.(type) {
		case float32:
			quat[i] = v
		case float64:
			quat[i] = float32(v)
		default:
			return nil, errors.New("quaternion item is not a float")
		}

	}
	return &quat, nil
}

func (e *Kv3Element) String() string {
	return e.StringIndent(0)
}

func (e *Kv3Element) StringIndent(tabs int) string {
	var ret string

	for k, v := range e.attributes {
		ret += strings.Repeat("\t", tabs) + k + ": "
		ret += valueToString(v, tabs)
		/*switch v := v.(type) {
		case int32:
			ret += strconv.Itoa(int(v))
		case string:
			ret += v
		case nil:
		case []Kv3Value:
		default:
			log.Println("Unknown type:", v)
		}*/

		ret += "\n"
	}

	return ret
}

func valueToString(v any, tabs int) string {
	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return "\"" + v + "\""
	case nil:
	case []Kv3Value:
		ret := "["
		for _, v2 := range v {
			ret += "\n" + strings.Repeat("\t", tabs+1) + valueToString(v2, tabs+1) + ","
		}
		ret += "\n" + strings.Repeat("\t", tabs) + "]"
		return ret
	case []byte:
		ret := "["
		count := 0
		for _, v2 := range v {
			ret += strconv.Itoa(int(v2)) + ","
			count++
			if count > 100 {
				break
			}
		}
		ret += "\n" + strings.Repeat("\t", tabs) + "]"
		return ret
	case *Kv3Element:
		ret := "{\n"
		ret += v.StringIndent(tabs + 1)
		ret += strings.Repeat("\t", tabs) + "}"
		return ret
	default:
		log.Println("Unknown type:", v)
	}

	return ""
}
