package parser

import (
	"encoding/binary"
	"github.com/baldurstod/go-source2-tools"
	"io"
	_"log"
)

type parseContext struct {
	reader io.Reader
	file   *source2.File
}

func newParseContext(reader io.Reader) *parseContext {
	return &parseContext{
		reader: reader,
		file:   source2.NewFile(),
	}
}

func Parse(r io.Reader) *source2.File {

	context := newParseContext(r)

	parseHeader(context)
	return context.file
}

func parseHeader(context *parseContext) {
	//var test uint32
	binary.Read(context.reader, binary.LittleEndian, &context.file.FileLength)
	binary.Read(context.reader, binary.LittleEndian, &context.file.VersionMaj)
	binary.Read(context.reader, binary.LittleEndian, &context.file.VersionMin)

	//log.Println(test)
	/*

		file.fileLength = reader.getUint32();
		file.versionMaj = reader.getUint16();
		file.versionMin = reader.getUint16();*/

}
