package dragon

import (
	"bufio"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	apiFiles = [...]string{
		"go1.txt",
		"go1.1.txt",
		"go1.2.txt",
		"go1.3.txt",
		"go1.4.txt",
		"go1.5.txt",
		"go1.6.txt",
		"go1.7.txt",
		"go1.8.txt",
		"go1.9.txt",
		"go1.10.txt",
	}
	sym = regexp.MustCompile(`^pkg (\S+).*?, (?:var|func|type|const) ([A-Z]\w*)`)
)

func stdLibs(libChan chan lib) error {
	apiDir := filepath.Join(build.Default.GOROOT, "api")
	eg := &errgroup.Group{}
	for _, f := range apiFiles {
		r, err := os.Open(filepath.Join(apiDir, f))
		if err != nil {
			return err
		}
		eg.Go(func() error {
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
			return sc.Err()
		})
	}
	return eg.Wait()
}
