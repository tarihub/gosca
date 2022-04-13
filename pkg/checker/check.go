package checker

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/output"
)

func CheckGoModule(goDepsList []config.Dependence, vulDbIdxMap config.VulDbIdxMap, vulnDbMap map[string]config.VulnDb) {
	var vulIdList []string

	for _, goDep := range goDepsList {
		if id, ok := goDep.GetMatches(vulDbIdxMap); ok {
			vulIdList = append(vulIdList, id)
		}
	}

	fmt.Printf("Scanning Go Modules Completed.\n\n")
	if len(vulIdList) > 0 {
		fmt.Println("[!] Vulnerability Found")
	}

	output.PrintVuln(vulIdList, vulnDbMap)

}
