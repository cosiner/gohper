package server

import (
	"os"
	"path/filepath"
)

type Translation struct {
	defLocale map[string]string
	locales   map[string]map[string]string
}

func Load(p string, tr *Translation) *Translation {

}

func (tr *Translation) Load(path string) {
	os.Lstat(path)
	filepath.Walk(root, walkFn)
}
