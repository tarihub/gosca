package checker

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/output"
	"strconv"
)

func CheckGoModule(goDepsList []config.Dependence, vulDbIdxMap config.VulDbIdxMap, vulnDbMap map[string]config.VulnDb) {
	var vulIdList []string

	for _, goDep := range goDepsList {
		if ids, ok := goDep.GetMatches(vulDbIdxMap); ok {
			vulIdList = append(vulIdList, ids...)
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

	}

	output.PrintVuln(vulIdList, vulnDbMap)

}
