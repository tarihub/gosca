package checker

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/util/file"
	"os"
	"strconv"
)

func CheckGoModule(
	goDepsList []config.Dependence, vulDbIdxMap config.VulDbIdxMap, vulnDbMap map[string]config.VulnDb,
	workDir string, excludeList []string, outFormat string) ([]string, config.Imports,
) {
	var vulIdList []string

	for _, goDep := range goDepsList {
		if ids, ok := goDep.GetMatches(vulDbIdxMap); ok {
			vulIdList = append(vulIdList, ids...)
		}
	}

	projectImports := config.Imports{}

	if len(vulIdList) != 0 {
		packagePaths, err := file.PackagePaths(workDir, file.ExcludedDirsRegExp(excludeList))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}

		projectImports.PackagePaths = packagePaths
		projectImports.PackageLocationMap = make(map[string][]string)
		projectImports.GetImports()

		// Check if the project imports any of the std libs vulnerable packages
		// NOTICE: We cannot use `for _, vulId := range vulIdList` here
		// because remove element from slice while iterating will skip the removed element location
		for i := 0; i < len(vulIdList); i++ {
			if _, ok := vulnDbMap[vulIdList[i]]; ok && vulnDbMap[vulIdList[i]].Module == "std" {
				pkg := vulnDbMap[vulIdList[i]].Package
				if _, ok := projectImports.PackageLocationMap[pkg]; !ok {
					vulIdList = append(vulIdList[:i], vulIdList[i+1:]...)
					i--
				}

				// Add function location map, but it's complex to get it...
				//existSymbol := false
				//for _, symbol := range vulnDbMap[vulIdList[i]].Symbols {
				//	if _, ok := projectImports.FuncLocationMap[pkg+"."+symbol]; ok {
				//		existSymbol = true
				//	}
				//}
				//if !existSymbol {
				//	vulIdList = append(vulIdList[:i], vulIdList[i+1:]...)
				//	i--
				//}
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
		os.Exit(0)
	}

	return vulIdList, projectImports
}
