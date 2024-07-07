package source2_test

import (
	"github.com/baldurstod/go-source2-tools/repository"
	"github.com/baldurstod/go-source2-tools/vpk"
)

const varFolder = "./var/"

func initRepo() {
	repository.AddRepository("dota2", vpk.NewVpkFS("R:\\SteamLibrary\\steamapps\\common\\dota 2 beta\\game\\dota\\pak01_dir.vpk"))
}
