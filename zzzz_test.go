package source2_test

import (
	"bytes"
	"github.com/baldurstod/go-source2-tools/parser"
	"log"
	"os"
	"path"
	"testing"
)

const varFolder = "./var/"

func Test(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b, _ := os.ReadFile(path.Join(varFolder, "pedestal_1.vmdl_c"))
	file, _ := parser.Parse(bytes.NewReader(b))
	log.Println(file)
}
