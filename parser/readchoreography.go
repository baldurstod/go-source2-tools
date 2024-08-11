package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/choreography"
	"github.com/ulikunitz/xz/lzma"
)

func parseVcdList(context *parseContext, block *source2.FileBlock) error {
	vcdList := source2.FileBlockVcdList{}
	var choreoCount uint32
	var err error
	var strings []string
	var o int64

	if err = binary.Read(context.reader, binary.LittleEndian, &vcdList.Version); err != nil {
		return fmt.Errorf("failed to read vcd list version: <%w>", err)
	}
	if err = binary.Read(context.reader, binary.LittleEndian, &choreoCount); err != nil {
		return fmt.Errorf("failed to read vcd list count: <%w>", err)
	}

	o, _ = context.reader.Seek(0, io.SeekCurrent)
	if strings, err = readVcdStrings(context.reader); err != nil {
		return err
	}

	context.reader.Seek(o+16, io.SeekStart) //skip strings data
	if _, err = readChoreographies(context.reader, choreoCount, strings); err != nil {
		return err
	}

	return nil
}

func readVcdStrings(reader io.ReadSeeker) ([]string, error) {
	var err error
	var offset uint32
	var stringsCount uint32
	var stringsOffset uint32
	var o int64

	o, _ = reader.Seek(0, io.SeekCurrent)
	if err = binary.Read(reader, binary.LittleEndian, &offset); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings offset offset: <%w>", err)
	}
	offset += uint32(o)
	if err = binary.Read(reader, binary.LittleEndian, &stringsCount); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings count: <%w>", err)
	}
	o, _ = reader.Seek(0, io.SeekCurrent)
	if err = binary.Read(reader, binary.LittleEndian, &stringsOffset); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings offset: <%w>", err)
	}
	stringsOffset += uint32(o)

	stringOffsets := make([]uint32, stringsCount)
	strings := make([]string, stringsCount)

	reader.Seek(int64(offset), io.SeekStart)
	for i := uint32(0); i < stringsCount; i++ {
		if err = binary.Read(reader, binary.LittleEndian, &stringOffsets[i]); err != nil {
			return nil, fmt.Errorf("failed to read string offset %d: <%w>", i, err)
		}
	}

	for i := uint32(0); i < stringsCount; i++ {
		reader.Seek(int64(stringsOffset+stringOffsets[i]), io.SeekStart)

		if strings[i], err = readNullTerminatedString(reader); err != nil {
			return nil, fmt.Errorf("failed to read string %d: <%w>", i, err)
		}
	}
	return strings, nil
}

func readChoreographies(reader io.ReadSeeker, count uint32, strings []string) ([]*choreography.Choreography, error) {
	var err error
	var o int64
	choreos := make([]*choreography.Choreography, count)
	const CHOREO_LENGTH = 24
	o, _ = reader.Seek(0, io.SeekCurrent)
	for i := uint32(0); i < count; i++ {
		reader.Seek(o+int64(i*CHOREO_LENGTH), io.SeekStart)
		if choreos[i], err = readChoreography(reader, strings); err != nil {
			return nil, err
		}
	}

	return choreos, nil
}

func readChoreography(reader io.ReadSeeker, strings []string) (*choreography.Choreography, error) {
	choreo := choreography.Choreography{}
	var nameOffset uint32
	var blockOffset uint32
	var blockLength uint32
	var hasSounds uint32
	var blockReader io.ReadSeeker
	var err error

	if nameOffset, err = readOffset(reader); err != nil {
		return nil, fmt.Errorf("failed to read name offset: <%w>", err)
	}
	if blockOffset, err = readOffset(reader); err != nil {
		return nil, fmt.Errorf("failed to read block offset: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &blockLength); err != nil {
		return nil, fmt.Errorf("failed to read choreography length: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &choreo.Duration); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &choreo.SoundDuration); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &hasSounds); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	choreo.HasSounds = hasSounds == 1

	reader.Seek(int64(nameOffset), io.SeekStart)
	if choreo.Name, err = readNullTerminatedString(reader); err != nil {
		return nil, fmt.Errorf("failed to read choreography name: <%w>", err)
	}

	reader.Seek(int64(blockOffset), io.SeekStart)
	if blockReader, err = readBlock(reader, blockLength); err != nil {
		return nil, fmt.Errorf("failed to read choreography block: <%w>", err)
	}

	log.Println(blockReader)

	log.Println(nameOffset, blockOffset)

	return &choreo, nil
}

func readOffset(reader io.ReadSeeker) (uint32, error) {
	var offset uint32
	o, _ := reader.Seek(0, io.SeekCurrent)
	if err := binary.Read(reader, binary.LittleEndian, &offset); err != nil {
		return 0, fmt.Errorf("failed to read offset: <%w>", err)
	}
	offset += uint32(o)
	return offset, nil
}

func readBlock(reader io.ReadSeeker, length uint32) (io.ReadSeeker, error) {
	var lzmaMagic uint32
	var err error
	var r *lzma.Reader

	if err = binary.Read(reader, binary.LittleEndian, &lzmaMagic); err != nil {
		return nil, fmt.Errorf("failed to read choreography block magic: <%w>", err)
	}

	if lzmaMagic == 0x414D5A4C { //LZMA
		var uncompressedSize uint32
		var compressedSize uint32
		//var properties [5]byte

		if err = binary.Read(reader, binary.LittleEndian, &uncompressedSize); err != nil {
			return nil, fmt.Errorf("failed to read lzma uncompressedSize: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &compressedSize); err != nil {
			return nil, fmt.Errorf("failed to read lzma compressedSize: <%w>", err)
		}

		b := make([]byte, compressedSize+13)
		buf := bytes.NewBuffer(make([]byte, 0, 8))
		if err = binary.Read(reader, binary.LittleEndian, b[0:5]); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}
		if err = binary.Write(buf, binary.LittleEndian, uint64(uncompressedSize)); err != nil {
			return nil, fmt.Errorf("failed to write choreography uncompressed size: <%w>", err)
		}
		copy(b[5:], buf.Bytes())
		if err = binary.Read(reader, binary.LittleEndian, b[13:]); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}

		if r, err = lzma.NewReader(bytes.NewReader(b)); err != nil {
			return nil, fmt.Errorf("failed to create lzma reader: <%w>", err)
		}

		var uncompressed = make([]byte, uncompressedSize)
		if _, err = r.Read(uncompressed[:]); err != nil {
			return nil, err
		}

		return bytes.NewReader(uncompressed), nil
	} else {
		reader.Seek(-4, io.SeekCurrent)
		//TODO: optimize
		b := make([]byte, length)
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}
		return bytes.NewReader(b), nil
	}
}

func readNullTerminatedString(reader io.ReadSeeker) (string, error) {
	builder := strings.Builder{}
	var b byte

	for {
		if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
			return "", fmt.Errorf("unable to read null string: <%w>", err)
		}
		if b == 0 {
			break
		}
		builder.WriteString(string(b))
	}

	return builder.String(), nil
}
