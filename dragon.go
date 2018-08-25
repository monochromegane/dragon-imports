package dragon

import (
	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/sync/errgroup"
)

// Imports generate zstdlib.go from api files and libs in GOPATH.
func Imports() error {
	if !existGoImports() {
		return errors.New("goimports command isn't installed")
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
		return gopathLibs(libChan)
	})
	err = eg.Wait()
	if err != nil {
		return err
	}
	close(libChan)

	err = <-done
	if err != nil {
		return err
	}

	return installUsing(tmp)
}

type lib struct {
	pkg    string
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
