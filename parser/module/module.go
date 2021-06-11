package module

import (
	"github.com/parvez3019/go-swagger3/parser/model"
	"os"
	"path/filepath"
	"strings"
)

type Parser interface {
	Parse() error
}

type parser struct {
	model.Utils
}

func NewParser(utils model.Utils) Parser {
	return &parser{
		Utils: utils,
	}
}

func (p *parser) Parse() error {
	walker := func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			if strings.HasPrefix(strings.Trim(strings.TrimPrefix(path, p.ModulePath), "/"), ".git") {
				return nil
			}
			fns, err := filepath.Glob(filepath.Join(path, "*.go"))
			if len(fns) == 0 || err != nil {
				return nil
			}
			// p.debug(path)
			name := filepath.Join(p.ModuleName, strings.TrimPrefix(path, p.ModulePath))
			name = filepath.ToSlash(name)
			p.KnownPkgs = append(p.KnownPkgs, model.Pkg{
				Name: name,
				Path: path,
			})
			p.KnownNamePkg[name] = &p.KnownPkgs[len(p.KnownPkgs)-1]
			p.KnownPathPkg[path] = &p.KnownPkgs[len(p.KnownPkgs)-1]
		}
		return nil
	}
	return filepath.Walk(p.ModulePath, walker)
}
