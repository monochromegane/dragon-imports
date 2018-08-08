package dragon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type walkFunc func(info fileInfo) error

func concurrentWalk(root string, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 16)
	return walk(newFileInfo(root, info), walkFn, sem)
}

func walk(info fileInfo, walkFn walkFunc, sem chan struct{}) error {
	walkError := walkFn(info)
	if walkError != nil {
		if info.IsDir() && walkError == filepath.SkipDir {
			return nil
		}
		return walkError
	}

	if !info.IsDir() {
		return nil
	}

	path := info.path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, file := range files {
		f := newFileInfo(filepath.Join(path, file.Name()), file)
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(file fileInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(file, walkFn, sem)
			}(f, wg)
		default:
			walk(f, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
