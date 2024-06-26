package source2_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/baldurstod/go-source2-tools/parser"
)

const varFolder = "./var/"

func Test(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b, _ := os.ReadFile(path.Join(varFolder, "pedestal_1.vmdl_c"))
	file, err := parser.Parse(b)
	if err != nil {
		log.Println(err)
	} else {
		//log.Println(file)
		log.Println(file.GetBlock("DATA"))
	}
}
