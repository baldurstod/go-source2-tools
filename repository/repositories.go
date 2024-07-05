package repository

var repositories = func() map[string]string { return make(map[string]string) }()

func AddRepository(key string, path string) {
	repositories[key] = path
}

func GetRepository(key string) string {
	return repositories[key]
}
