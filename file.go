package source2

import (
	"errors"
	"strings"

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
	Content FileBlockContent
}

func (fb *FileBlock) String() string {
	if fb.Content != nil {
		return fb.Content.String()
	}

	return ""
}

func (fb *FileBlock) GetBlockStruct(path []string) (any, error) {
	if fb.Content != nil {
		return fb.Content.GetBlockStruct(path)
	}

	return nil, errors.New("block content is empty")
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

func (f *File) GetBlockStruct(path string) (any, error) {
	v := strings.Split(path, ".")

	if len(v) < 1 {
		return nil, errors.New("path is too short: " + path)
	}

	block := f.GetBlock(v[0])

	if block == nil {
		return nil, errors.New("block not found: " + v[0])
	}

	return block.GetBlockStruct(v[1:])
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

func (rerl *FileBlockRERL) String() string {
	panic("TODO")
}

func (rerl *FileBlockRERL) GetBlockStruct(path []string) (any, error) {
	panic("TODO")
}

type FileBlockContent interface {
	String() string
	GetBlockStruct(path []string) (any, error)
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

func (data *FileBlockDATA) GetBlockStruct(path []string) (any, error) {
	//panic("TODO")
	if len(path) < 1 {
		return nil, nil
	}

	current := data.KeyValue
	if current == nil {
		return nil, errors.New("data block don't have key value")
	}

	var v any
	var ok bool

	v = current.GetAttribute(path[0])

	for _, s := range path[1:] {

		current, ok = v.(*kv3.Kv3Element)
		if !ok {
			return nil, errors.New("can't convert to Kv3Element")
		}
		//element = current.(*Kv3Element)
		v = current.GetAttribute(s)

		if v == nil {
			return nil, nil
		}
		//GetAttribute
		//ret += "\n\t" + valueToString(v2) + ","
	}

	return v, nil

	/*

		var arr = path.split('.');
		var data = this.blocks;
		if (!data) {
			return null;
		}

		var sub;
		for (var i = 0; i < arr.length; i++) {
			sub = data[arr[i]];
			if (!sub) {
				return null;
			}
			data = sub;
		}

		return data;
	*/
}

type FileBlockMBUF struct {
}

func NewFileBlockMBUF() *FileBlockMBUF {
	return &FileBlockMBUF{}
}

func (mbuf *FileBlockMBUF) String() string {
	panic("TODO")
}

func (mbuf *FileBlockMBUF) GetBlockStruct(path []string) (any, error) {
	panic("TODO")
}
