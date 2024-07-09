//go:build dump

package source2_test

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/baldurstod/go-source2-tools/parser"
)

func TestDump(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initRepo()

	var filename string
	filename = "models/heroes/rubick/rubick.vmdl_c"
	//filename = "models/heroes/wisp/wisp.vmdl_c"
	filename = "models/heroes/tiny_01/tiny_01.vmdl_c"

	//b, _ := os.ReadFile(path.Join(varFolder, filename))
	file, err := parser.Parse("dota2", filename)
	if err != nil {
		t.Error(err)
	}

	//log.Println(file)

	buf := new(bytes.Buffer)
	for _, v := range file.BlocksArray {
		buf.WriteString(v.ResType + " {\n")
		buf.WriteString(v.String())
		buf.WriteString("}\n")
	}

	os.WriteFile(path.Join(varFolder, "tiny_01.dump.txt"), buf.Bytes(), 0666)
}
