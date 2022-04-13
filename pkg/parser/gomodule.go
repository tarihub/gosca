package parser

import (
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/util/utils"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// ParseGoModOrSum parses a go.mod file or a go.sum file, result store in goDepsMap
func ParseGoModOrSum(goDepsList *[]config.Dependence, filePath string, regex string) error {
	tmpGoDepsMap := make(map[string][]string)

	fc, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	for _, match := range regexp.MustCompile(regex).FindAllStringSubmatch(string(fc), -1) {
		if len(match) != 3 {
			continue
		}
		name := strings.Replace(strings.Trim(match[1], `'"`), "/go.mod", "", -1)
		ver := match[2]
		tmpGoDepsMap[name] = append(tmpGoDepsMap[name], ver)
	}

	// Get Go version
	if path.Base(filePath) == "go.mod" {
		verRegex := `(go)\s+(\d+(.\d+){0,2})`
		goVersion := regexp.MustCompile(verRegex).FindString(string(fc))
		tmpGoDepsMap["std"] = append(tmpGoDepsMap["std"], strings.Replace(goVersion, " ", "", -1))
	}

	// duplicate result
	for k, v := range tmpGoDepsMap {
		tmpGoDepsMap[k] = utils.UniqueSlice(v)
	}

	for module, versions := range tmpGoDepsMap {
		*goDepsList = append(*goDepsList, config.Dependence{Module: module, Versions: versions})
	}

	return nil
}
