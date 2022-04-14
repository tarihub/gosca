package output

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"strings"
)

func PrintVuln(vulIdList []string, vulnDbMap map[string]config.VulnDb) {
	for _, vulId := range vulIdList {
		if _, ok := vulnDbMap[vulId]; ok {
			fmt.Println("\t" + strings.Repeat("=", 80))
			fmt.Printf("\t%s", vulId)
			var vul = []config.VulnDb{vulnDbMap[vulId]}
			printAddition(vul, "\t", 80)
		}
	}
}

func printAddition(additionalPacks []config.VulnDb, tab string, splitDashNum int) {
	for _, vul := range additionalPacks {
		dash := strings.Repeat("-", splitDashNum)
		fmt.Println("\n" + tab + dash)
		if vul.Module != "" {
			fmt.Printf("%s[Module] \n%s%s\n\n", tab, tab, vul.Module)
		}

		if vul.Package != "" {
			fmt.Printf("%s[Package] \n%s%s\n\n", tab, tab, vul.Package)
		}

		if len(vul.AdditionalPackages) > 0 {
			fmt.Printf("%s[Additional Packages] \n", tab)
			printAddition(vul.AdditionalPackages, tab+"\t", 20)
		}

		var versionsList []string
		for _, v := range vul.Versions {
			version := ""
			if _, ok := v["introduced"]; ok {
				version += ">= " + v["introduced"] + ", "
			}
			if _, ok := v["fixed"]; ok {
				version += "< " + v["fixed"]
			}

			versionsList = append(versionsList, version)
		}
		versions := strings.Join(versionsList, "\n"+tab)
		fmt.Printf("%s[Affects Versions] \n%s%s\n\n", tab, tab, versions)

		if vul.Description != "" {
			fmt.Printf("%s[Description] \n%s%s\n\n", tab, tab, strings.Replace(vul.Description, "\n", "\n"+tab, -1))
		}

		if vul.Published != "" {
			fmt.Printf("%s[Published] \n%s%s\n\n", tab, tab, vul.Published)
		}

		if len(vul.Cves) > 0 {
			cves := strings.Join(vul.Cves, "\n"+tab)
			fmt.Printf("%s[CVE] \n%s%s\n\n", tab, tab, cves)
		}

		if len(vul.Symbols) > 0 {
			symbols := strings.Join(vul.Symbols, "\n"+tab)
			fmt.Printf("%s[Symbols] \n%s%s\n\n", tab, tab, symbols)
		}

		if len(vul.DerivedSymbols) > 0 {
			derivedSymbols := strings.Join(vul.DerivedSymbols, "\n"+tab)
			fmt.Printf("%s[Derived Symbols] \n%s%s\n\n", tab, tab, derivedSymbols)
		}

		var linksList []string
		if vul.Links.Commit != "" {
			linksList = append(linksList, vul.Links.Commit)
		}
		if vul.Links.Pr != "" {
			linksList = append(linksList, vul.Links.Pr)
		}
		if vul.Links.Context != nil {
			linksList = append(linksList, vul.Links.Context...)
		}
		links := strings.Join(linksList, "\n"+tab)
		if links != "" {
			fmt.Printf("%s[References] \n%s%s\n\n", tab, tab, links)
		}
	}
}
