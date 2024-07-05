package repository

import "io/fs"

var repositories = func() map[string]fs.ReadFileFS { return make(map[string]fs.ReadFileFS) }()

func AddRepository(key string, fs fs.ReadFileFS) {
	repositories[key] = fs
}

func GetRepository(key string) fs.ReadFileFS {
	return repositories[key]
}
