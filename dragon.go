package dragon

import "sync"

func Imports() error {

	libChan := make(chan lib, 1000)
	done := make(chan struct{})

	go out(libChan, done)

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
