package checker

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/output"
	"github.com/TARI0510/gosca/pkg/util/file"
	"github.com/TARI0510/gosca/pkg/util/utils"
	"os"
	"strconv"
)

func CheckGoModule(goDepsList []config.Dependence, vulDbIdxMap config.VulDbIdxMap, vulnDbMap map[string]config.VulnDb, workDir string, excludeList []string) {
	var vulIdList []string

	for _, goDep := range goDepsList {
		if ids, ok := goDep.GetMatches(vulDbIdxMap); ok {
			vulIdList = append(vulIdList, ids...)
		}
	}

	if len(vulIdList) != 0 {
		packagePaths, err := file.PackagePaths(workDir, file.ExcludedDirsRegExp(excludeList))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}

		projectImports := config.Imports{
			PackagePaths:       packagePaths,
			PackageLocationMap: make(map[string][]string),
		}

		// Compatible running directory is not in Go project
		cacheCwd, _ := os.Getwd()
		_ = os.Chdir(workDir)
		projectImports.GetImports()
		_ = os.Chdir(cacheCwd)

		// Check if the project imports any of the std libs vulnerable packages
		for _, vulId := range vulIdList {
			if _, ok := vulnDbMap[vulId]; ok {
				if vulnDbMap[vulId].Module == "std" {
					if _, ok := projectImports.PackageLocationMap[vulnDbMap[vulId].Package]; !ok {
						vulIdList = utils.RemoveElem(vulIdList, vulId)
					}
				}
			}
		}
	}

	fmt.Printf("Scanning Go Modules Completed.\n\n")
	vulNum := len(vulIdList)
	if vulNum > 1 {
		fmt.Println("[!] " + strconv.Itoa(vulNum) + " Vulnerabilities Found")
	} else if vulNum == 1 {
		fmt.Println("[!] " + strconv.Itoa(vulNum) + " Vulnerability Found")
	} else {
		fmt.Println("[+] No Vulnerability Found")
		return
	}

	// Output result
	output.PrintVuln(vulIdList, vulnDbMap)

}
