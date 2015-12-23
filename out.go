package dragon

import "fmt"

func out(libChan chan lib, done chan struct{}) {
	for lib := range libChan {
		fmt.Printf("%s.%s: %s\n", lib.pkg, lib.object, lib.path)
	}
	done <- struct{}{}
}
