package parser

import (
	_ "encoding/binary"
	_ "fmt"
	"github.com/baldurstod/go-source2-tools/kv3"
	"io"
	_ "log"
)

type parseKv3Context struct {
	reader io.ReadSeeker
	root   *kv3.Kv3Element
}

func newParseKv3Context(reader io.ReadSeeker) *parseKv3Context {
	return &parseKv3Context{
		reader: reader,
		root:   kv3.NewKv3Element(),
	}
}

func ParseKv3(r io.ReadSeeker) (*kv3.Kv3Element, error) {
	context := newParseKv3Context(r)
	/*
		log.Println("Start parsing file")
		err := parseHeader(context)
		if err != nil {
			return nil, err
		}

		err = parseBlocks(context)
		if err != nil {
			return nil, err
		}

		log.Println("End parsing file")*/
	return context.root, nil
}

/*
func parseHeader(context *parseKv3Context) error {
}
*/
