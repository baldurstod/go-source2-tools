package parser

import (
	"bytes"
	_ "encoding/binary"
	_ "fmt"
	"github.com/baldurstod/go-source2-tools/kv3"
	"io"
	"log"
	"math"
)

type parseKv3Context struct {
	reader           io.ReadSeeker
	root             *kv3.Kv3Element
	stringDictionary []string
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
	if version >= 2 && blobCount != 0 {
		if compressedBlobReader != nil {
			uncompressedLength := blobCount * 4
			compressedEnd := blobOffset + uncompressedLength
			uncompressedBlobSizeReader := bytes.NewReader(b[blobOffset : blobOffset+uncompressedLength])
			compressedBlobSizeReader := bytes.NewReader(b[compressedEnd+4 : compressedEnd+4+blobCount*2])
			//uncompressedBlobSizeReader = new BinaryReader(reader, blobOffset, uncompressedLength);
			//compressedBlobSizeReader = new BinaryReader(reader, blobOffset + 4 + uncompressedLength, blobCount * 2);

			log.Println(uncompressedBlobSizeReader, compressedBlobSizeReader)

		} else {
			if uncompressedBlobReader != nil {
				//uncompressedBlobSizeReader = new BinaryReader(reader, reader.byteLength - blobCount * 4 - 4, blobCount * 4);
			}
		}
	}
	/*
		err := parseHeader(context)
		if err != nil {
			return nil, err
		}

		err = parseBlocks(context)
		if err != nil {
			return nil, err
		}
	*/
	log.Println("End parsing kv3")
	return context.root, nil
}
