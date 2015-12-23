package dragon

import (
	"errors"
	"os"
	"sync"
)

func Imports() error {

	if !existGoImports() {
		return errors.New("goimports command isn't installed.")
	}

	libChan := make(chan lib, 1000)
	done := make(chan struct{})

	go updateZstdlib(libChan, done)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		stdLibs(libChan)
		wg.Done()
	}()

	go func() {
		gopathLibs(libChan)
		wg.Done()
	}()

	wg.Wait()
	close(libChan)
	<-done

	install()

	return nil
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
