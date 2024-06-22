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
	file, err := parser.Parse(bytes.NewReader(b))
	if err != nil {
		log.Println(err)
	} else {
		log.Println(file)
		log.Println(file.GetBlock("AGRP"))
	}

}
