package dragon

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func updateZstdlib(libChan chan lib) error {
	f, err := os.OpenFile(filepath.Join(outPath(), "zstdlib.go"), os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return out(libChan, f)
}

func out(libChan chan lib, w io.Writer) error {
	libs := map[string]lib{}
	ambiguous := map[string]bool{}
	var keys []string

	for lib := range libChan {
		full := lib.path
		key := lib.pkg + "." + lib.object
		if exist, ok := libs[key]; ok {
			if exist.path != full {
				ambiguous[key] = true
			}
		} else {
			libs[key] = lib
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	outf := func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}
	outf("// AUTO-GENERATED BY dragon-imports\n\n")
	outf("package imports\n")
	outf("var stdlib = map[string]string{\n")

	for _, key := range keys {
		if ambiguous[key] {
			outf("\t// %q is ambiguous\n", key)
		} else {
			outf("\t%q: %q,\n", key, libs[key].path)
		}
	}
	outf("}\n")
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
