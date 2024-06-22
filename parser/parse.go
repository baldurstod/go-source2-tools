package parser

import (
	"encoding/binary"
	"errors"
	"github.com/baldurstod/go-source2-tools"
	"io"
	"log"
)

type parseContext struct {
	reader io.ReadSeeker
	file   *source2.File
}

func newParseContext(reader io.ReadSeeker) *parseContext {
	return &parseContext{
		reader: reader,
		file:   source2.NewFile(),
	}
}

func Parse(r io.ReadSeeker) (*source2.File, error) {
	context := newParseContext(r)

	log.Println("Start parsing file")
	err := parseHeader(context)
	if err != nil {
		return nil, err
	}

	err = parseBlocks(context)
	if err != nil {
		return nil, err
	}

	log.Println("End parsing file")
	return context.file, nil
}

func parseHeader(context *parseContext) error {
	//var test uint32
	var headerOffset uint32
	var resCount uint32
	reader := context.reader
	binary.Read(reader, binary.LittleEndian, &context.file.FileLength)
	binary.Read(reader, binary.LittleEndian, &context.file.VersionMaj)
	binary.Read(reader, binary.LittleEndian, &context.file.VersionMin)

	currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
	binary.Read(reader, binary.LittleEndian, &headerOffset)
	headerOffset += uint32(currentPos)
	binary.Read(reader, binary.LittleEndian, &resCount)

	reader.Seek(int64(headerOffset), io.SeekStart)

	resType := make([]byte, 4)
	var resOffset uint32
	var resLength uint32
	for i := uint32(0); i < resCount; i++ {
		read, _ := reader.Read(resType)
		if read != 4 {
			return errors.New("Failed to read block type")
		}
		binary.Read(reader, binary.LittleEndian, &resOffset)
		binary.Read(reader, binary.LittleEndian, &resLength)

		context.file.AddBlock(string(resType), resOffset, resLength)
	}
	return nil
}

func parseBlocks(context *parseContext) error {
	for _, block := range context.file.BlocksArray {
		err := parseBlock(context, block)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseBlock(context *parseContext, block *source2.FileBlock) error {

	log.Println(context, block)

	return nil
}
