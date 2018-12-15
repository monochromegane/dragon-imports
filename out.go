package dragon

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
)

func out(libChan chan lib, w io.Writer) error {
	libs := map[string]lib{}
	ambiguous := map[string]bool{}
	for lib := range libChan {
		full := lib.path
		key := lib.pkg + "." + lib.object
		if exist, ok := libs[key]; ok {
			if exist.path != full {
				ambiguous[key] = true
			}
		} else {
			libs[key] = lib
		}
	}

	stdlib := map[string]map[string]bool{
		"unsafe": map[string]bool{
			"Alignof":       true,
			"ArbitraryType": true,
			"Offsetof":      true,
			"Pointer":       true,
			"Sizeof":        true,
		},
	}
	for _, lib := range libs {
		if ambiguous[lib.pkg+"."+lib.object] {
			continue
		}
		objMap, ok := stdlib[lib.path]
		if !ok {
			objMap = make(map[string]bool)
		}
		objMap[lib.object] = true
		stdlib[lib.path] = objMap
	}

	var buf bytes.Buffer
	outf := func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}
	outf("// AUTO-GENERATED BY dragon-imports\n\n")
	outf("package imports\n")
	outf("var stdlib = %#v\n", stdlib)

	fmtbuf, _ := format.Source(buf.Bytes())
	_, err := w.Write(fmtbuf)
	return err
}

func outPath() string {
	for _, src := range srcDirs() {
		outPath := filepath.Join(src, "golang.org/x/tools/imports")
		if _, err := os.Stat(outPath); err == nil {
			return outPath
		}
	}
	return ""
}
