package dragon

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

type mod struct {
	version string
	lib     lib
}

var regPath = regexp.MustCompile("@v[0-9]+\\.[0-9]+\\.[0-9].*?(/|$)")

func extractImportPathAndVersion(path string) (string, string) {
	version := regPath.FindString(path)
	importPath := regPath.ReplaceAllString(path, "/")
	return strings.TrimSuffix(importPath, "/"), strings.Trim(version, "@/")
}

type versions struct {
	libByVersion map[string][]lib
}

func (v *versions) append(version string, l lib) {
	v.libByVersion[version] = append(v.libByVersion[version], l)
}

func (v *versions) latest() []lib {
	var versions semver.Versions
	for version, _ := range v.libByVersion {
		sv, err := semver.ParseTolerant(version)
		if err != nil {
			continue
		}
		versions = append(versions, sv)
	}
	semver.Sort(versions)

	latest := versions[len(versions)-1]
	if libs, ok := v.libByVersion["v"+latest.String()]; ok {
		return libs
	}
	return []lib{}
}

func gomoduleLibs(libChan chan lib) error {
	modChan := make(chan mod, 1000)
	done := make(chan bool)

	go func() {
		modules := map[string]*versions{}
		for mod := range modChan {
			vs, ok := modules[mod.lib.path]
			if !ok {
				vs = &versions{libByVersion: map[string][]lib{}}
				modules[mod.lib.path] = vs
			}
			vs.append(mod.version, mod.lib)
		}
		for _, v := range modules {
			for _, lib := range v.latest() {
				libChan <- lib
			}
		}
		done <- true
	}()

	for _, modDir := range modDirs() {
		ignoreDirs, err := fetchGoImportsIgnore(modDir)
		if err != nil {
			return err
		}
		err = concurrentWalk(modDir, func(info fileInfo) error {
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

			importPath, err := filepath.Rel(modDir, filepath.Dir(path))
			if err != nil {
				return nil
			}

			importPath, version := extractImportPathAndVersion(importPath)
			for _, v := range f.Scope.Objects {
				if ast.IsExported(v.Name) {
					modChan <- mod{
						version: version,
						lib: lib{
							object: v.Name,
							path:   importPath,
						},
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	close(modChan)
	<-done
	return nil
}
