package source2_test

import (
	"log"
	"testing"

	"github.com/baldurstod/go-source2-tools/parser"
	"github.com/baldurstod/go-source2-tools/particles"
)

func TestPsFiles(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	filename := "particles/units/heroes/hero_dark_willow/dark_willow_lantern_ambient_fairy.vpcf_c"

	//b, _ := os.ReadFile(path.Join(varFolder, filename))
	f, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
	log.Println(f)
}

func TestSystem(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	filename := "particles/units/heroes/hero_dark_willow/dark_willow_lantern_ambient_fairy.vpcf_c"

	file, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	system := particles.NewParticleSystem()
	system.SetFile(file)

	config, err := system.GetControlPointConfiguration("preview")
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(config)
}

/*
func TestAnim(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var filename string
	filename = "pedestal_1.vmdl_c"
	filename = "drow_base.vmdl_c"
	filename = "muerta_base.vmdl_c"
	filename = "snapfire.vmdl_c"
	filename = "primal_beast_base.vmdl_c"
	filename = "dragon_knight.vmdl_c"
	filename = "void_spirit.vmdl_c"
	filename = "models/heroes/wisp/wisp.vmdl_c"
	filename = "models/heroes/rubick/rubick.vmdl_c"
	filename = "models/heroes/terrorblade/terrorblade.vmdl_c"
	filename = "models/heroes/dawnbreaker/dawnbreaker.vmdl_c"

	//b, _ := os.ReadFile(path.Join(varFolder, filename))
	file, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	model := model.NewModel()
	model.SetFile(file)

	skel, err := model.GetSkeleton()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(skel)

	flexes, err := model.GetFlexes()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(flexes)

	modifiers := make(map[string]struct{})
	modifiers["PostGameIdle"] = struct{}{}
	seq, err := model.GetSequence("ACT_DOTA_LOADOUT", modifiers)
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(seq, err)
	//model.PrintSequences()
	log.Println(seq.GetFps(), seq.GetFrameCount())
	frame, err := seq.GetFrame(30)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(frame)

}*/
