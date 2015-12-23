package dragon

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func gopathLibs(libChan chan lib) {

	for _, srcDir := range srcDirs() {
		concurrentWalk(srcDir, func(path string, info fileInfo) error {

			if info.isDir(false) {
				name := info.Name()
				if name == "" || name[0] == '.' || name[0] == '_' || name == "testdata" {
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
	}
}
