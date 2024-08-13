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
		t.Error(err)
		return
	}

	log.Println(f)
	/*
		var b []byte
		if b, err = json.MarshalIndent(f.Blocks["DATA"].Content.(*source2.FileBlockVcdList).Choreographies, "", "\t"); err != nil {
			t.Error(err)
			return
		}
		os.WriteFile(path.Join(varFolder, "tusk_vcd.json"), b, 0666)
	*/
}
