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
	case int32:
		return strconv.Itoa(int(v))
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 32)
	case string:
		return v
	case nil:
	case []Kv3Value:
		ret := "["
		for _, v2 := range v {
			/*switch v2 := v2.(type) {
				case Uint
			default:
				log.Println("Unknown type2:", v2)
			}*/
			ret += "\n\t" + valueToString(v2) + ","
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
