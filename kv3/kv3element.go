package kv3

import (
	"log"
	"strconv"
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

	for k, _ := range e.attributes {
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

func (e *Kv3Element) GetFloat32Attribute(name string) (float32, bool) {
	value, ok := e.attributes[name]
	if !ok {
		return 0, false
	}

	f, ok := value.(float32)
	if !ok {
		return 0, false
	}

	return f, true
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

func (e *Kv3Element) String() string {
	var ret string

	for k, v := range e.attributes {
		ret += k + ":\t"
		ret += valueToString(v)
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

func valueToString(v any) string {
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
		return v
	case nil:
	case []Kv3Value:
		ret := "["
		for _, v2 := range v {
			ret += "\n\t" + valueToString(v2) + ","
		}
		ret += "\n]"
		return ret
	case []byte:
		ret := "["
		for _, v2 := range v {
			ret += "\n\t" + strconv.Itoa(int(v2)) + ","
		}
		ret += "\n]"
		return ret
	case *Kv3Element:
		return v.String()
	default:
		log.Println("Unknown type:", v)
	}

	return ""
}

type Kv3Value interface {
}
