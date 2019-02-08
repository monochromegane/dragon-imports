package dragon

import (
	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/sync/errgroup"
)

// Imports generate zstdlib.go from api files and libs in GOPATH.
func Imports(restore, gomodules bool) error {
	if !existGoImports() {
		return errors.New("goimports command isn't installed")
	}
	if restore {
		return install()
	}

	libChan := make(chan lib, 1000)
	done := make(chan error)
	tmp, err := ioutil.TempFile(outPath(), "dragon-imports")
	if err != nil {
		return err
	}
	defer func(fname string) {
		tmp.Close()
		os.Remove(fname)
	}(tmp.Name())

	go func() {
		done <- out(libChan, tmp)
	}()

	eg := &errgroup.Group{}
	eg.Go(func() error {
		return stdLibs(libChan)
	})
	eg.Go(func() error {
		if gomodules {
			return gomoduleLibs(libChan)
		} else {
			return gopathLibs(libChan)
		}
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	close(libChan)

	err = <-done
	if err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return installUsing(tmp.Name())
}

type lib struct {
	object string
	path   string
}

func existGoImports() bool {
	for _, path := range [...]string{outPath(), cmdPath()} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}
	}
	return true
}
