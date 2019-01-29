package dragon

import (
	"bufio"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"
)

var sym = regexp.MustCompile(`^pkg (\S+).*?, (?:var|func|type|const) ([A-Z]\w*)`)

func stdLibs(libChan chan lib) error {
	apiFiles, err := filepath.Glob(filepath.Join(build.Default.GOROOT, "api", "go1.*txt"))
	if err != nil {
		return err
	}
	eg := &errgroup.Group{}
	for _, f := range apiFiles {
		r, err := os.Open(f)
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
