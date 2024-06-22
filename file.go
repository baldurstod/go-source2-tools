package source2

type File struct {
	FileLength uint32
	VersionMaj uint16
	VersionMin uint16
	FileBlocks []*FileBlock
}

func NewFile() *File {
	return &File{
		FileBlocks: make([]*FileBlock, 0),
	}
}

type FileBlock struct {
	*File
	ResType string
	Offset  uint32
	Length  uint32
}

func (f *File) AddFileBlock(ResType string, Offset uint32, Length uint32) {
	fb := &FileBlock{
		File:    f,
		ResType: ResType,
		Offset:  Offset,
		Length:  Length,
	}

	f.FileBlocks = append(f.FileBlocks, fb)
}
