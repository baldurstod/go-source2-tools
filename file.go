package source2

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type File struct {
	FileLength  uint32
	VersionMaj  uint16
	VersionMin  uint16
	Blocks      map[string]*FileBlock // Blocks stores the first block of a particular type
	BlocksArray []*FileBlock
}

func NewFile() *File {
	return &File{
		Blocks:      make(map[string]*FileBlock),
		BlocksArray: make([]*FileBlock, 0),
	}
}

type FileBlock struct {
	*File
	ResType string
	Offset  uint32
	Length  uint32
	Data    any
}

func (fb *FileBlock) String() string {
	switch fb.Data.(type) {
	case *FileBlockDATA:
		return fb.Data.(*FileBlockDATA).String()

	}

	return ""
}

func (f *File) AddBlock(resType string, offset uint32, length uint32) {
	fb := &FileBlock{
		File:    f,
		ResType: resType,
		Offset:  offset,
		Length:  length,
	}

	_, exist := f.Blocks[resType]
	// Add the first block of this type to the map
	if !exist {
		f.Blocks[resType] = fb
	}

	f.BlocksArray = append(f.BlocksArray, fb)
}

func (f *File) GetBlock(resType string) *FileBlock {
	return f.Blocks[resType]
}

type FileBlockRERL struct {
	ExternalFilesByIndex  []string
	ExternalFilesByHandle map[string]string
}

func NewFileBlockRERL() *FileBlockRERL {
	return &FileBlockRERL{
		ExternalFilesByIndex:  make([]string, 0),
		ExternalFilesByHandle: make(map[string]string),
	}
}

func (fb *FileBlockRERL) AddExternalFile(handle string, filename string) {
	fb.ExternalFilesByIndex = append(fb.ExternalFilesByIndex, filename)
	fb.ExternalFilesByHandle[handle] = filename
}

type FileBlockDATA struct {
	KeyValue *kv3.Kv3Element
}

func NewFileBlockDATA() *FileBlockDATA {
	return &FileBlockDATA{}
}

func (fb *FileBlockDATA) String() string {
	if fb.KeyValue != nil {
		return fb.KeyValue.String()
	} else {
		return ""
	}
}
