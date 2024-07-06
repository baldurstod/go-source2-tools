package vpk

import (
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/NublyBR/go-vpk"
)

type VpkFS struct {
	path string
	vpk  vpk.VPK
}

func NewVpkFS(path string) *VpkFS {
	vpk := VpkFS{
		path: path,
	}

	vpk.init()

	return &vpk
}

func (fs *VpkFS) init() error {
	vpk, err := vpk.OpenDir(fs.path)

	if err != nil {
		return err
	}

	fs.vpk = vpk
	return nil
}

func (fs *VpkFS) Open(name string) (fs.File, error) {
	return nil, nil
}

func (fs *VpkFS) ReadFile(path string) ([]byte, error) {
	entry, ok := fs.vpk.Find(path)

	if !ok {
		return nil, errors.New("file not found")
	}

	fileReader, err := entry.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open file: <%w>", err)
	}

	buf, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: <%w>", err)
	}

	return buf, nil
}
