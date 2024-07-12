package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/baldurstod/go-vector"
	"github.com/x448/float16"
)

type Decoder struct {
	Name         string
	Version      int
	Type         int
	BytesPerBone int
}

func (dec *Decoder) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	dec.Name, ok = datas.GetStringAttribute("m_szName")
	if !ok {
		return errors.New("unable to get decoder name")
	}

	dec.Version, _ = datas.GetIntAttribute("m_nVersion")
	dec.Type, _ = datas.GetIntAttribute("m_nType")

	switch dec.Name {
	case "CCompressedStaticVector4D", "CCompressedFullVector4D":
		dec.BytesPerBone = 16
	case "CCompressedStaticFullVector3", "CCompressedFullVector3", "CCompressedDeltaVector3", "CCompressedAnimVector3":
		dec.BytesPerBone = 12
	case "CCompressedStaticVector2D", "CCompressedFullVector2D":
		dec.BytesPerBone = 8
	case "CCompressedAnimQuaternion", "CCompressedStaticVector3", "CCompressedStaticQuaternion", "CCompressedFullQuaternion":
		dec.BytesPerBone = 6
	case "CCompressedFullFloat", "CCompressedStaticFloat", "CCompressedStaticInt", "CCompressedFullInt", "CCompressedStaticColor32", "CCompressedFullColor32":
		dec.BytesPerBone = 4
	case "CCompressedFullShort":
		dec.BytesPerBone = 2
	case "CCompressedStaticChar", "CCompressedFullChar", "CCompressedStaticBool", "CCompressedFullBool":
		dec.BytesPerBone = 1
	default:
		return errors.New("unknown decoder type: " + dec.Name)
	}

	return nil
}

func (dec *Decoder) decode(reader *bytes.Reader, frameIndex int, boneIndex int, boneCount int) (any, error) {
	switch dec.Name {
	case "CCompressedStaticVector3":
		reader.Seek(int64(8+boneCount*2+boneIndex*dec.BytesPerBone), io.SeekStart)
		v := [3]uint16{}

		err := binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to read CCompressedStaticVector3: <%w>", err)
		}

		return vector.Vector3[float32]{
			float16.Frombits(v[0]).Float32(),
			float16.Frombits(v[1]).Float32(),
			float16.Frombits(v[2]).Float32(),
		}, nil
	case "CCompressedStaticFullVector3":
		reader.Seek(int64(8+boneCount*2+boneIndex*dec.BytesPerBone), io.SeekStart)
		v := vector.Vector3[float32]{}

		err := binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to read CCompressedStaticFullVector3: <%w>", err)
		}

		return v, nil
	case "CCompressedAnimQuaternion":
		reader.Seek(int64(8+boneCount*(2+frameIndex*dec.BytesPerBone)+boneIndex*dec.BytesPerBone), io.SeekStart)
		buf := make([]byte, 6)
		_, err := reader.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to read CCompressedAnimQuaternion: <%w>", err)
		}

		return readQuaternion48(buf), nil
	case "CCompressedDeltaVector3":
		//TODO: no sure how this block is supposed to work
		frameIndex = 0
		reader.Seek(int64(8+boneCount*(2+frameIndex*dec.BytesPerBone)+boneIndex*dec.BytesPerBone), io.SeekStart)
		v := vector.Vector3[float32]{}

		err := binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to read CCompressedStaticFullVector3: <%w>", err)
		}

		return v, nil
	case "CCompressedFullVector3":
		reader.Seek(int64(8+boneCount*(2+frameIndex*dec.BytesPerBone)+boneIndex*dec.BytesPerBone), io.SeekStart)
		v := vector.Vector3[float32]{}

		err := binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to read CCompressedStaticFullVector3: <%w>", err)
		}

		return v, nil
	default:
		return nil, errors.New("unknown decoder type: " + dec.Name)
	}
}

/*
	{
		m_szName = "CCompressedStaticFloat"
		m_nVersion = 0
		m_nType = 1
	},
*/

const QUATERNION48_SCALE = math.Sqrt2 / 0x8000

func readQuaternion48(buf []byte) *vector.Quaternion[float32] {
	// Values
	i1 := int(buf[0]) + ((int(buf[1]) & 127) << 8) - 0x4000
	i2 := int(buf[2]) + ((int(buf[3]) & 127) << 8) - 0x4000
	i3 := int(buf[4]) + ((int(buf[5]) & 127) << 8) - 0x4000

	// Signs
	s1 := buf[1] & 128
	s2 := buf[3] & 128
	s3 := buf[5] & 128

	x := QUATERNION48_SCALE * float32(i1)
	y := QUATERNION48_SCALE * float32(i2)
	z := QUATERNION48_SCALE * float32(i3)
	w := float32(math.Sqrt(float64(1 - (x * x) - (y * y) - (z * z))))

	// Apply sign 3
	if s3 == 128 {
		w *= -1
	}

	// Apply sign 1 and 2
	if s1 == 128 {
		if s2 == 128 {
			return &vector.Quaternion[float32]{y, z, w, x}
		} else {
			return &vector.Quaternion[float32]{z, w, x, y}
		}
	} else {
		if s2 == 128 {
			return &vector.Quaternion[float32]{w, x, y, z}
		} else {
			return &vector.Quaternion[float32]{x, y, z, w}
		}
	}
}
