package parser

import (
	"encoding/binary"
	"fmt"
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
		_, err := reader.Read(resType)
		if err != nil {
			return fmt.Errorf("Failed to read block type in parseHeader: <%w>", err)
		}
		currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
		binary.Read(reader, binary.LittleEndian, &resOffset)
		binary.Read(reader, binary.LittleEndian, &resLength)

		context.file.AddBlock(string(resType), resOffset+uint32(currentPos), resLength)
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

	switch block.ResType {
	case "RERL":
		err := parseRERL(context, block)
		if err != nil {
			return fmt.Errorf("An error occured in parseBlock: <%w>", err)
		}
	}

	return nil
}

func parseRERL(context *parseContext, block *source2.FileBlock) error {
	log.Println(context, block)

	var resOffset uint32
	var resCount uint32
	var strOffset int32
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)

	binary.Read(reader, binary.LittleEndian, &resOffset)
	binary.Read(reader, binary.LittleEndian, &resCount)

	log.Println(block.Offset, resOffset, resCount)

	fileBlockRERL := source2.NewFileBlockRERL(block)

	handle := make([]byte, 8)
	for i := uint32(0); i < resCount; i++ {
		context.reader.Seek(int64(block.Offset+resOffset+16*i), io.SeekStart)
		log.Println(block.Offset + resOffset + 16*i)

		_, err := reader.Read(handle)
		if err != nil {
			return fmt.Errorf("Failed to read handle in parseRERL: <%w>", err)
		}
		binary.Read(reader, binary.LittleEndian, &strOffset)
		context.reader.Seek(int64(strOffset-4), io.SeekCurrent)
		filename, err := readNullString(reader)
		if err != nil {
			return fmt.Errorf("Failed to read filename in parseRERL: <%w>", err)
		}

		fileBlockRERL.AddExternalFile(string(handle[:]), filename)

		//readHandle(reader)
	}
	log.Println(block.Offset, resOffset, resCount, fileBlockRERL)
	/*
		function loadRerl(reader, block) {
			reader.seek(block.offset);
			var resOffset = reader.getInt32();// Seems to be always 0x00000008
			var resCount = reader.getInt32();
			block.externalFiles = {};
			block.externalFiles2 = [];

			reader.seek(block.offset + resOffset);


			for (var resIndex = 0; resIndex < resCount; resIndex++) {
				reader.seek(block.offset + resOffset + 16 * resIndex);
				var handle = readHandle(reader);//reader.getUint64(fieldOffset);
				var strOffset = reader.getInt32();
				reader.skip(strOffset - 4);
				var s = reader.getNullString();
				block.externalFiles[handle] = s;
				block.externalFiles2[resIndex] = s;
			}
		}
	*/

	return nil
}

func readNullString(r io.ReadSeeker) (string, error) {
	// Probably not the fastest way to do that
	s := ""
	c := make([]byte, 1)

	for {
		_, err := r.Read(c)
		if err != nil {
			return s, fmt.Errorf("An error occured while reading null string: <%w>", err)
		}
		if c[0] == 0 {
			break
		} else {
			s += string(c[0])
		}
	}

	log.Println(s)
	return s, nil
}
