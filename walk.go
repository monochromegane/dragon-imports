package dragon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type walkFunc func(path string, info fileInfo) error

func concurrentWalk(root string, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 16)
	return walk(root, newFileInfo(root, info), walkFn, sem)
}

func walk(path string, info fileInfo, walkFn walkFunc, sem chan struct{}) error {
	walkError := walkFn(path, info)
	if walkError != nil {
		if info.IsDir() && walkError == filepath.SkipDir {
			return nil
		}
		return walkError
	}

	if !info.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, file := range files {
		f := newFileInfo(path, file)
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file fileInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(path, file, walkFn, sem)
			}(filepath.Join(path, file.Name()), f, wg)
		default:
			walk(filepath.Join(path, file.Name()), f, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
