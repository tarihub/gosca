package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/TARI0510/gosca/assets"
	"github.com/TARI0510/gosca/pkg/checker"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/output"
	"github.com/TARI0510/gosca/pkg/parser"
	"github.com/TARI0510/gosca/pkg/util/file"
	"os"
	"path"
	"strings"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, " ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var (
	workDir     string
	vulnDbDir   string
	version     bool
	excludeDirs arrayFlags
	outFormat   string
	outPath     string
)

func init() {
	flag.StringVar(&workDir, "w", ".", "work directory which contains the go.mod(sum) file")
	flag.StringVar(&vulnDbDir, "vulndb-dir", "", "vulndb path, yaml from https://github.com/golang/vulndb/tree/master/reports")
	flag.StringVar(&outFormat, "f", "screen", "output format: screen, excel")
	flag.StringVar(&outPath, "o", "", "output path")
	// TODO Add include and exclude vulndb-dir flag
	flag.BoolVar(&version, "v", false, "show gosca version")
	flag.Var(&excludeDirs, "e", "Exclude folder from go std libs scan (can be specified multiple times)")
	err := flag.Set("e", "vendor")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: failed to exclude the %q directory from scan", "vendor")
	}
	err = flag.Set("e", "\\.git")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: failed to exclude the %q directory from scan", ".git")
	}

	flag.Parse()

	if version {
		fmt.Println(config.Version)
		os.Exit(0)
	}
}

func main() {
	outputMap := map[string]output.OutBuilder{
		"screen": output.PrintVuln,
		"excel":  output.ExcelVuln,
	}

	if _, ok := outputMap[outFormat]; !ok {
		fmt.Fprintf(os.Stderr, "Unsupport output format\n")
		os.Exit(1)
	}

	goModPath := path.Join(workDir, "go.mod")
	goSumPath := path.Join(workDir, "go.sum")

	_, goModErr := os.Stat(goModPath)
	_, goSumErr := os.Stat(goSumPath)

	if goSumErr != nil && goModErr != nil {
		fmt.Println(goModPath + " and " + goSumPath + " not found")
		os.Exit(1)
	}

	var vulnDBs map[string]string
	var err = error(nil)
	if vulnDbDir == "" {
		vulnDBs = assets.VulnDBs
	} else {
		vulnDBs, err = file.ReadSuffixFiles(vulnDbDir, []string{"*.yaml", "*.yml"})
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	if len(vulnDBs) == 0 {
		fmt.Fprintf(os.Stderr, "No vulnerability rule load from %q directory\n", vulnDbDir)
		os.Exit(1)
	}

	// Read and parse all vulndb yaml files
	vulnDbMap := make(config.VulnDbMap)
	vulDbIdxMap := vulnDbMap.ReadVulnDbYaml(vulnDBs)

	// Parse work_dir/go.mod and work_dir/sum.mod
	goDepsList := make([]config.Dependence, 0)

	e1, e2 := error(nil), error(nil)

	if goModErr == nil {
		e1 = parser.ParseGoModOrSum(&goDepsList, goModPath, `(\S*)\s+(v[\d\w\-+.]*)[\s\n]`)
		if e1 != nil {
			fmt.Fprintf(os.Stderr, "[skip] "+goModPath+", "+e1.Error()+"\n")
		}
	}

	if goModErr == nil {
		e2 = parser.ParseGoModOrSum(&goDepsList, goSumPath, `(\S*)\s+(v[\d\w\-+.]*)/go.mod[\s\n]`)
		if e2 != nil {
			fmt.Fprintf(os.Stderr, "[skip] "+goSumPath+", "+e2.Error()+"\n")
		}
	}

	if e1 != nil && e2 != nil {
		os.Exit(1)
	}

	fmt.Println("[*] Working on " + workDir)
	vulIdList, projectImports := checker.CheckGoModule(goDepsList, vulDbIdxMap, vulnDbMap, workDir, excludeDirs, outFormat)

	fmt.Println("[*] Handle output ... ")

	f := output.F{WorkDir: workDir, OutFormat: outFormat, OutPath: outPath}

	outputMap[outFormat](vulIdList, vulnDbMap, projectImports.PackageLocationMap, f)
}
