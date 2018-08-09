package dragon

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func installUsing(f *os.File) error {
	current := filepath.Join(outPath(), "zstdlib.go")
	backup := strings.Replace(current, ".go", ".g_", 1)

	err := os.Rename(current, backup)
	if err != nil {
		return err
	}
	defer os.Rename(backup, current)

	err = os.Rename(f.Name(), current)
	if err != nil {
		return err
	}

	return install()
}

func install() error {
	cmd := exec.Command("go", "install", "-a", ".")
	cmd.Dir = cmdPath()
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
