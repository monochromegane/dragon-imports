package dragon

import (
	"go/build"
	"path/filepath"
)

func srcDirs() []string {
	gopaths := filepath.SplitList(build.Default.GOPATH)

	srcDirs := make([]string, len(gopaths))
	for i, gopath := range gopaths {
		srcDirs[i] = filepath.Join(gopath, "src")
	}
	return srcDirs
}
