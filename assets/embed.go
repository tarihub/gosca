package assets

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"strings"
)

//go:embed vulndb
var f embed.FS

var VulnDBs = make(map[string]string)

func init() {
	dirEntries, err := f.ReadDir("vulndb")
	if err != nil {
		panic(err)
	}

	subFs, err := fs.Sub(f, "vulndb")
	if err != nil {
		panic(err)
	}

	for _, dir := range dirEntries {
		if !strings.HasSuffix(dir.Name(), ".yml") && !strings.HasSuffix(dir.Name(), ".yaml") {
			continue
		}

		file, err := subFs.Open(dir.Name())
		if err != nil {
			return
		}

		c, err := ioutil.ReadAll(file)
		if err != nil {
			return
		}

		VulnDBs[dir.Name()] = string(c)
	}
}
