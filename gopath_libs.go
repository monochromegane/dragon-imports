package dragon

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func fetchGoImportsIgnore(src string) (map[string]bool, error) {
	dirs := make(map[string]bool)
	f, err := os.Open(filepath.Join(src, ".goimportsignore"))
	if err != nil {
		return dirs, nil
	}
	scr := bufio.NewScanner(f)
	for scr.Scan() {
		dir := scr.Text()
		if !strings.HasPrefix(dir, "#") {
			dirs[filepath.Join(src, dir)] = true
		}
	}
	return dirs, scr.Err()
}

func isSkipDir(fi fileInfo, ignoreDirs map[string]bool) bool {
	name := fi.Name()
	switch name {
	case "", "testdata", "vendor":
		return true
	}
	switch name[0] {
	case '.', '_':
		return true
	}
	return ignoreDirs[filepath.Join(fi.path, fi.Name())]
}

func gopathLibs(libChan chan lib) error {
	for _, srcDir := range srcDirs() {
		ignoreDirs, err := fetchGoImportsIgnore(srcDir)
		if err != nil {
			return err
		}
		err = concurrentWalk(srcDir, func(info fileInfo) error {
			if info.isDir(false) {
				if isSkipDir(info, ignoreDirs) {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(info.Name(), "_test.go") {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".go") {
				return nil
			}

			path := info.path
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return nil
			}

			pkg := f.Name.Name
			if pkg == "main" {
				return nil
			}

			importPath, err := filepath.Rel(srcDir, filepath.Dir(path))
			if err != nil {
				return nil
			}

			for _, v := range f.Scope.Objects {
				if ast.IsExported(v.Name) {
					libChan <- lib{
						pkg:    pkg,
						object: v.Name,
						path:   importPath,
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
