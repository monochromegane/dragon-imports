package dragon

import "testing"

func TestVersionsLatest(t *testing.T) {

	expects := map[string]map[string][]lib{
		"111": map[string][]lib{
			"v1.0.1": []lib{lib{path: "p", object: "101"}, lib{path: "p", object: "101"}},
			"v1.1.1": []lib{lib{path: "p", object: "111"}, lib{path: "p", object: "111"}},
			"v1.0.0": []lib{lib{path: "p", object: "100"}, lib{path: "p", object: "100"}},
		},
		"2": map[string][]lib{
			"v0.0.0-1": []lib{lib{path: "p", object: "1"}, lib{path: "p", object: "1"}},
			"v0.0.0-2": []lib{lib{path: "p", object: "2"}, lib{path: "p", object: "2"}},
			"v0.0.0-0": []lib{lib{path: "p", object: "0"}, lib{path: "p", object: "0"}},
		},
	}

	for expect, v := range expects {
		versions := &versions{libByVersion: map[string][]lib{}}
		for version, libs := range v {
			for _, lib := range libs {
				versions.append(version, lib)
			}
		}
		libs := versions.latest()
		for _, lib := range libs {
			if lib.object != expect {
				t.Errorf("versions.latest should return latest libs")
			}
		}
	}
}

func TestExtractImportPathAndVersion(t *testing.T) {
	expects := map[string][]string{
		"golang.org/x/net@v0.0.0-20190125091013-d26f9f9a57f3/nettest": []string{
			"golang.org/x/net/nettest",
			"v0.0.0-20190125091013-d26f9f9a57f3",
		},
		"golang.org/x/tools@v0.0.0-20190201231825-51e363b66d25/godoc/dl": []string{
			"golang.org/x/tools/godoc/dl",
			"v0.0.0-20190201231825-51e363b66d25",
		},
		"github.com/donvito/hellomod/v2@v2.0.0": []string{
			"github.com/donvito/hellomod/v2",
			"v2.0.0",
		},
		"github.com/donvito/hellomod/v2@v2.0.0-alpha.1.beta": []string{
			"github.com/donvito/hellomod/v2",
			"v2.0.0-alpha.1.beta",
		},
	}

	for path, expect := range expects {
		importPath, version := extractImportPathAndVersion(path)
		if importPath != expect[0] {
			t.Errorf("extractImportPathAndVersion should return import path %s, but %s", expect[0], importPath)
		}
		if version != expect[1] {
			t.Errorf("extractImportPathAndVersion should return version %s, but %s", expect[1], version)
		}
	}
}
