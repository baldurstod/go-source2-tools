package source2

type File struct {
	FileLength uint32
	VersionMaj uint16
	VersionMin uint16
}

func NewFile() *File {
	return &File{}
}
