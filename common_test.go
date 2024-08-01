package source2_test

import (
	"log"

	"github.com/baldurstod/go-source2-tools/repository"
	"github.com/baldurstod/go-source2-tools/vpk"
)

const varFolder = "./var/"

var _ = func() bool {
	repository.AddRepository("dota2", vpk.NewVpkFS("R:\\SteamLibrary\\steamapps\\common\\dota 2 beta\\game\\dota\\pak01_dir.vpk"))
	return true
}()

var _ = func() bool { log.SetFlags(log.LstdFlags | log.Lshortfile); return true }()
