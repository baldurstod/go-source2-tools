package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/baldurstod/go-source2-tools/kv3"
	"io"
	"log"
	"math"
)

type parseKv3Context struct {
	reader           io.ReadSeeker
	root             *kv3.Kv3Element
	stringDictionary []string
	decompressOffset int
}

func newParseKv3Context(reader io.ReadSeeker) *parseKv3Context {
	return &parseKv3Context{
		reader:           reader,
		root:             kv3.NewKv3Element(),
		stringDictionary: make([]string, 0),
	}
}

func ParseKv3(b []byte, version int, singleByteCount uint32, quadByteCount uint32, eightByteCount uint32, dictionaryTypeLength uint32,
	blobCount uint32, totalUncompressedBlobSize uint32, compressedBlobReader io.ReadSeeker, uncompressedBlobReader io.ReadSeeker, compressionFrameSize uint16) (*kv3.Kv3Element, error) {
	context := newParseKv3Context(bytes.NewReader(b))
	log.Println("Start parsing kv3")

	quadCursor := math.Ceil(float64(singleByteCount)/4) * 4
	eightCursor := math.Ceil((quadCursor+float64(quadByteCount)*4)/8) * 8

	dictionaryOffset := uint32(eightCursor) + eightByteCount*8
	blobOffset := dictionaryOffset + dictionaryTypeLength
	blobEnd := uint32(len(b) - 4)

	var uncompressedBlobSizeReader, compressedBlobSizeReader io.ReadSeeker
	if version >= 2 && blobCount != 0 {
		if compressedBlobReader != nil {
			uncompressedLength := blobCount * 4
			compressedEnd := blobOffset + uncompressedLength
			uncompressedBlobSizeReader = bytes.NewReader(b[blobOffset : blobOffset+uncompressedLength])
			compressedBlobSizeReader = bytes.NewReader(b[compressedEnd+4 : compressedEnd+4+blobCount*2])

			log.Println(uncompressedBlobSizeReader, compressedBlobSizeReader)

		} else {
			if uncompressedBlobReader != nil {
				uncompressedBlobSizeReader = bytes.NewReader(b[blobEnd-blobCount*4 : blobEnd])
				log.Println(uncompressedBlobSizeReader)
			}
		}
	}

	var offset uint32
	if version == 1 {
		offset = blobEnd
	} else if version >= 2 {
		offset = blobOffset
	}

	// First compute the size
	size := 0
	s := offset
	for {
		s--
		if s <= 0 {
			break
		}
		if b[s] != 0 {
			size++
		} else {
			break
		}
	}

	typeArray := b[s+1 : offset]
	valueArray := make([]any, s)
	//let valueArray = [];
	//byteReader := bytes.NewReader(b)
	quadReader := bytes.NewReader(b)
	eightReader := bytes.NewReader(b)

	quadReader.Seek(int64(quadCursor), io.SeekStart)
	eightReader.Seek(int64(eightCursor), io.SeekStart)

	var stringCount uint32
	err := binary.Read(quadReader, binary.LittleEndian, &stringCount)
	if err != nil {
		return nil, fmt.Errorf("Failed to read stringCount in ParseKv3: <%w>", err)
	}

	context.reader.Seek(int64(dictionaryOffset), io.SeekStart)
	stringDictionary := make([]string, stringCount)
	readStringDictionary(context, stringDictionary, stringCount)

	var decompressBlobBuffer []byte
	//let decompressBlobArray;

	if compressedBlobReader != nil { //if a compressed reader is provided, we have to uncompress the blobs
		decompressBlobBuffer = make([]byte, totalUncompressedBlobSize) //new ArrayBuffer(totalUncompressedBlobSize);
		/*decompressBlobArray = new Uint8Array(decompressBlobBuffer);
		decompressBlobArray.decompressOffset = 0;*/
		context.decompressOffset = 0
	}

	rootElement := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, nil, compressedBlobReader, uncompressedBlobReader, typeArray, valueArray, -1, false, compressionFrameSize)

	log.Println("End parsing kv3", size, typeArray, stringCount, rootElement)
	return context.root, nil
}

func readStringDictionary(context *parseKv3Context, stringDictionary []string, stringCount uint32) error {
	for i := uint32(0); i < stringCount; i++ {
		s, err := readNullString(context.reader)
		if err != nil {
			return fmt.Errorf("Failed to read string in readStringDictionary: <%w>", err)
		}
		stringDictionary = append(stringDictionary, s)
	}
	return nil
}

func parseBinaryKv3Element(context *parseKv3Context, quadReader *bytes.Reader, eightReader *bytes.Reader, uncompressedBlobSizeReader io.ReadSeeker, compressedBlobSizeReader io.ReadSeeker, blobCount uint32,
	decompressBlobBuffer []byte, decompressBlobArray any, compressedBlobReader io.ReadSeeker, uncompressedBlobReader io.ReadSeeker, typeArray []byte, valueArray []any, elementType int, isArray bool, compressionFrameSize uint16) error {
	return nil
}
