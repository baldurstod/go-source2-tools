package source2_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/baldurstod/go-source2-tools/model"
	"github.com/baldurstod/go-source2-tools/parser"
	"github.com/baldurstod/go-source2-tools/repository"
	"github.com/baldurstod/go-source2-tools/vpk"
)

const varFolder = "./var/"

func initRepo() {
	repository.AddRepository("dota2", vpk.NewVpkFS("R:\\SteamLibrary\\steamapps\\common\\dota 2 beta\\game\\dota\\pak01_dir.vpk"))
}

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

	b, _ := os.ReadFile(path.Join(varFolder, filename))
	_, err := parser.Parse(b)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
}

func DisabledTestModel(t *testing.T) {
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

	b, _ := os.ReadFile(path.Join(varFolder, filename))
	file, err := parser.Parse(b)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	model := model.NewModel()
	model.SetFile(file)

	skel, _ := model.GetSkeleton()

	model.GetAnimationData(nil)

	seq := model.GetSequence("ACT_DOTA_IDLE", nil)

	log.Println(skel.GetBones(), seq)
}

func TestAnim(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initRepo()

	var filename string
	filename = "pedestal_1.vmdl_c"
	filename = "drow_base.vmdl_c"
	filename = "muerta_base.vmdl_c"
	filename = "snapfire.vmdl_c"
	filename = "primal_beast_base.vmdl_c"
	filename = "dragon_knight.vmdl_c"
	filename = "void_spirit.vmdl_c"
	filename = "wisp.vmdl_c"

	b, _ := os.ReadFile(path.Join(varFolder, filename))
	file, err := parser.Parse(b)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	model := model.NewModel()
	model.SetFile(file)

	seq := model.GetSequence("ACT_DOTA_IDLE", nil)

	log.Println(seq)
	//log.Println(seq.GetFps())
}

func TestRepo(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initRepo()
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
