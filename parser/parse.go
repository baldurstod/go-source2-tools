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

		res := source2.NewFileBlock(context.file, string(resType), resOffset, resLength)

		log.Println(res, read)

		/*
			resType = reader.getString(4);
			resOffset = reader.tell() + reader.getUint32();
			resLength = reader.getUint32();

			file.maxBlockOffset = Math.max(file.maxBlockOffset, resOffset + resLength);

			block = new Source2FileBlock(file, resType, resOffset, resLength);
			file.addBlock(block);
		*/
	}

	log.Println(headerOffset, resCount)
	/*

		file.fileLength = reader.getUint32();
		file.versionMaj = reader.getUint16();
		file.versionMin = reader.getUint16();*/
	return nil

}
