package dragon

import (
	"bufio"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	apiFiles = [...]string{"go1.txt", "go1.1.txt", "go1.2.txt", "go1.3.txt", "go1.4.txt", "go1.5.txt"}
	sym      = regexp.MustCompile(`^pkg (\S+).*?, (?:var|func|type|const) ([A-Z]\w*)`)
)

func stdLibs(libChan chan lib) {
	apiDir := filepath.Join(build.Default.GOROOT, "api")

	var wg sync.WaitGroup
	for _, f := range apiFiles {
		r, err := os.Open(filepath.Join(apiDir, f))
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				l := sc.Text()
				has := func(v string) bool { return strings.Contains(l, v) }
				if has("struct, ") || has("interface, ") || has(", method (") {
					continue
				}

				if m := sym.FindStringSubmatch(l); m != nil {
					libChan <- lib{
						pkg:    path.Base(m[1]),
						object: m[2],
						path:   m[1],
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
