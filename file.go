package source2

type File struct {
	FileLength uint32
	VersionMaj uint16
	VersionMin uint16
}

func NewFile() *File {
	return &File{}
}

type FileBlock struct {
	*File
	ResType string
	Offset  uint32
	Length  uint32
}

func NewFileBlock(file *File, ResType string, Offset uint32, Length uint32) *FileBlock {
	return &FileBlock{
		File:    file,
		ResType: ResType,
		Offset:  Offset,
		Length:  Length,
	}
}
