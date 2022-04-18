package checker

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/output"
	"github.com/TARI0510/gosca/pkg/util/utils"
	"strconv"
)

func CheckGoModule(goDepsList []config.Dependence, vulDbIdxMap config.VulDbIdxMap, vulnDbMap map[string]config.VulnDb, packagePaths []string) {
	var vulIdList []string

	for _, goDep := range goDepsList {
		if ids, ok := goDep.GetMatches(vulDbIdxMap); ok {
			vulIdList = append(vulIdList, ids...)
		}
	}

	projectImports := config.Imports{
		PackagePaths:       packagePaths,
		PackageLocationMap: make(map[string][]string),
	}
	if len(vulIdList) != 0 {
		projectImports.GetImports()

		// Check if the project imports any of the std libs vulnerable packages
		for idx, vulId := range vulIdList {
			if _, ok := vulnDbMap[vulId]; ok {
				if vulnDbMap[vulId].Module == "std" {
					if _, ok := projectImports.PackageLocationMap[vulnDbMap[vulId].Package]; !ok {
						vulIdList = utils.RemoveElemDisorder(vulIdList, idx)
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

	output.PrintVuln(vulIdList, vulnDbMap)

}
