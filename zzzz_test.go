package source2_test

import (
	"log"
	"testing"

	"github.com/baldurstod/go-source2-tools/model"
	"github.com/baldurstod/go-source2-tools/parser"
	"github.com/baldurstod/go-source2-tools/repository"
)

func DisabledTestFiles(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var filename string
	filename = "pedestal_1.vmdl_c"
	filename = "drow_base.vmdl_c"
	filename = "muerta_base.vmdl_c"
	filename = "snapfire.vmdl_c"
	filename = "primal_beast_base.vmdl_c"
	filename = "dragon_knight.vmdl_c"
	filename = "void_spirit.vmdl_c"

	//b, _ := os.ReadFile(path.Join(varFolder, filename))
	_, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
}

func TestModel(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var filename string
	filename = "pedestal_1.vmdl_c"
	filename = "drow_base.vmdl_c"
	filename = "muerta_base.vmdl_c"
	filename = "snapfire.vmdl_c"
	filename = "primal_beast_base.vmdl_c"
	filename = "dragon_knight.vmdl_c"
	filename = "void_spirit.vmdl_c"
	filename = "wisp.vmdl_c"
	filename = "models/items/dawnbreaker/dawnbreaker_astral_angel_armor/dawnbreaker_astral_angel_armor.vmdl_c"
	filename = "models/heroes/rattletrap/rattletrap_rocket.vmdl_c"

	file, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	m := model.NewModel()
	m.SetFile(file)

	skel, _ := m.GetSkeleton()

	m.GetAnimationData(nil)

	seq, err := m.GetSequence(model.NewActivity("ACT_DOTA_IDLE"))

	log.Println(skel.GetBones(), seq, err)
}

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

	m := model.NewModel()
	m.SetFile(file)

	skel, err := m.GetSkeleton()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(skel)

	flexes, err := m.GetFlexes()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(flexes)

	seq, err := m.GetSequence(model.NewActivity("ACT_DOTA_LOADOUT", "PostGameIdle"))
	if err != nil {
		t.Error(err)
		return
	}
	/*seq, err = model.GetSequenceByName("@rubick_run_haste_turns")
	if err != nil {
		t.Error(err)
		return
	}*/
	/*modifiers := make(map[string]struct{})
	modifiers["centaur_mount"] = struct{}{}
	seq, err = model.GetSequence("ACT_DOTA_CAST_ABILITY_5", modifiers)*/

	log.Println(seq, err)
	//model.PrintSequences()
	log.Println(seq.GetFps(), seq.GetFrameCount())
	frame, err := seq.GetFrame(30)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(frame)

}

func DisabledTestRepo(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	repo := repository.GetRepository("dota2")
	if repo == nil {
		t.Error("repo not found")
	}

	buf, err := repo.ReadFile("models/heroes/tiny_01/tiny_01.vmdl_c")
	log.Println(buf[0:200], err)
}

/*
func TestSkeleton(t *testing.T) {
	s := model.NewSkeleton(100)

	b0 := model.NewBone("bone 0")
	b1 := model.NewBone("bone 1")
	b2 := model.NewBone("bone 2")
	b3 := model.NewBone("bone 3")

	s.AddBone(b0)
	s.AddBone(b1)
	s.AddBone(b2)
	s.AddBone(b3)

	b1.ParentBone = b0

	if _, err := s.GetBoneById(2); err != nil {
		t.Error("GetBoneById returned qn error")
	}
	if _, err := s.GetBoneById(4); err == nil {
		t.Error("no error returned")
	}

	log.Println(s)
}*/

func TestAttachment(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	filename := "models/heroes/dark_willow/dark_willow.vmdl_c"

	file, err := parser.Parse("dota2", filename)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	model := model.NewModel()
	model.SetFile(file)

	atta, err := model.GetAttachement("attach_lantern")
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(atta)
}
