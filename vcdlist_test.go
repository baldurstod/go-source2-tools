package source2_test

import (
	"log"
	"testing"

	"github.com/baldurstod/go-source2-tools/parser"
)

func TestTuskVcdList(t *testing.T) {

	filename := "scenes/tusk.vcdlist_c"

	//b, _ := os.ReadFile(path.Join(varFolder, filename))
	f, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	log.Println(f)
}
