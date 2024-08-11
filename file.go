package source2

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/baldurstod/go-source2-tools/choreography"
	"github.com/baldurstod/go-source2-tools/kv3"
)

type File struct {
	repository  string
	filename    string
	FileType    string
	FileLength  uint32
	VersionMaj  uint16
	VersionMin  uint16
	Blocks      map[string]*FileBlock // Blocks stores the first block of a particular type
	BlocksArray []*FileBlock
}

func NewFile(repo string, filename string) *File {
	f := File{
		repository:  repo,
		filename:    filename,
		Blocks:      make(map[string]*FileBlock),
		BlocksArray: make([]*FileBlock, 0),
	}

	f.FileType, _ = strings.CutSuffix(filepath.Ext(filename), "_c")
	f.FileType, _ = strings.CutPrefix(f.FileType, ".")

	return &f
}

type FileBlock struct {
	*File
	ResType  string
	Offset   uint32
	Length   uint32
	Content  FileBlockContent
	FileType string
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

func (f *File) AddBlock(fileType string, resType string, offset uint32, length uint32) {
	fb := &FileBlock{
		FileType: fileType,
		File:     f,
		ResType:  resType,
		Offset:   offset,
		Length:   length,
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

func (f *File) GetBlockStructAsKv3Element(path string) (*kv3.Kv3Element, error) {
	res, err := f.GetBlockStruct(path)
	if err != nil {
		return nil, err
	}

	elem, ok := res.(*kv3.Kv3Element)
	if !ok {
		return nil, errors.New("unable to convert value to Kv3Element")
	}
	return elem, nil
}

func (f *File) GetBlockStructAsKv3ValueArray(path string) ([]kv3.Kv3Value, error) {
	res, err := f.GetBlockStruct(path)
	if err != nil {
		return nil, err
	}

	elem, ok := res.([]kv3.Kv3Value)
	if !ok {
		return nil, errors.New("unable to convert value to array of Kv3Value")
	}
	return elem, nil
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
	return "Code String() for FileBlockRERL\n"
}

func (rerl *FileBlockRERL) GetBlockStruct(path []string) (any, error) {
	panic("FileBlockRERL GetBlockStruct\n")
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
		return fb.KeyValue.StringIndent(1)
	} else {
		return ""
	}
}

func (data *FileBlockDATA) GetBlockStruct(path []string) (any, error) {
	//panic("TODO")
	if len(path) < 1 {
		return data.KeyValue, nil
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
	panic("FileBlockMBUF String")
}

func (mbuf *FileBlockMBUF) GetBlockStruct(path []string) (any, error) {
	panic("FileBlockMBUF GetBlockStruct")
}

type FileBlockVcdList struct {
	Version        uint32
	Strings        []string
	Choreographies []*choreography.Choreography
}

func NewFileBlockVcdList() *FileBlockVcdList {
	return &FileBlockVcdList{}
}

func (mbuf *FileBlockVcdList) String() string {
	panic("FileBlockVcdList String")
}

func (mbuf *FileBlockVcdList) GetBlockStruct(path []string) (any, error) {
	panic("FileBlockVcdList GetBlockStruct")
}
