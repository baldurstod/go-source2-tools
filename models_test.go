package source2_test

import (
	"log"
	"testing"

	"github.com/baldurstod/go-source2-tools/parser"
)

func TestModel1(t *testing.T) {
	// This model had an issue related to lz4 decompression
	_, err := parser.Parse("dota2", "models/items/puck/puck_ti10_immortal_wings/puck_ti10_immortal_wings.vmdl_c")
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
}
