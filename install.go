package dragon

import (
	"os"
	"os/exec"
	"path/filepath"
)

func install() {
	cmd := exec.Command("go", "install", "-a", ".")
	cmd.Dir = cmdPath()
	cmd.Run()
}

func cmdPath() string {
	for _, src := range srcDirs() {
		cmdPath := filepath.Join(src, "golang.org/x/tools/cmd/goimports")
		if _, err := os.Stat(cmdPath); err == nil {
			return cmdPath
		}
	}
	return ""
}
