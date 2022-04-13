package main

import (
	"flag"
	"fmt"
	"github.com/TARI0510/gosca/pkg/checker"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/parser"
	"github.com/TARI0510/gosca/pkg/util/file"
	"os"
	"path"
	"strings"
)

var (
	workDir   string
	vulnDbDir string
	version   bool
)

func init() {
	flag.StringVar(&workDir, "work_dir", "", "home directory which contains the go.mod(sum) file")
	flag.StringVar(&vulnDbDir, "vulndb_dir", "vulndb", "vulndb path, yaml from https://github.com/golang/vulndb/tree/master/reports")
	flag.BoolVar(&version, "v", false, "show gosca version")
	flag.Parse()

	if version {
		fmt.Println(config.Version)
		os.Exit(0)
	}
}

func main() {
	var baseDir string
	cwd, _ := os.Getwd()

	if !strings.HasPrefix(workDir, "/") {
		workDir = path.Join(cwd, workDir)
	}

	goModPath := path.Join(workDir, "go.mod")
	goSumPath := path.Join(workDir, "go.sum")

	_, goModErr := os.Stat(goModPath)
	_, goSumErr := os.Stat(goSumPath)

	if goSumErr != nil && goModErr != nil {
		fmt.Println(goModPath + " and " + goSumPath + " not found")
		os.Exit(1)
	}

	if !strings.HasPrefix(vulnDbDir, "/") {
		vulnDbDir = path.Join(baseDir, vulnDbDir)
	}

	yamlPaths, err := file.ListSuffixFiles(vulnDbDir, []string{".yaml", ".yml"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read and parse all vulndb yaml files
	vulnDbMap := make(config.VulnDbMap)
	vulDbIdxMap := vulnDbMap.ReadVulnDbYaml(yamlPaths)

	// parse work_dir/go.mod and work_dir/sum.mod
	goDepsList := make([]config.Dependence, 0)

	e1, e2 := error(nil), error(nil)

	if goModErr == nil {
		e1 = parser.ParseGoModOrSum(&goDepsList, goModPath, `(\S*)\s+(v[\d\w\-+.]*)[\s\n]`)
		if e1 != nil {
			fmt.Println("[skip] " + goModPath + ", " + e1.Error())
		}
	}

	if goModErr == nil {
		e2 = parser.ParseGoModOrSum(&goDepsList, goSumPath, `(\S*)\s+(v[\d\w\-+.]*)/go.mod[\s\n]`)
		if e2 != nil {
			fmt.Println("[skip] " + goSumPath + ", " + e2.Error())
		}
	}

	if e1 != nil && e2 != nil {
		os.Exit(1)
	}

	checker.CheckGoModule(goDepsList, vulDbIdxMap, vulnDbMap)
}
