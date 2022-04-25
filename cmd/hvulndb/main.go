// Handle assets/vulndb
// Add CVSS score and Severity to assets/vulndb/*.yaml
// according https://nvd.nist.gov/vuln/detail/[CVE_ID]
// Skip if not set CVE in *.yaml
package main

import (
	"flag"
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/util/file"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

var (
	vulnDbDir string
)

func getCVEScore(cve string) ([]string, error) {
	doc, err := htmlquery.LoadURL("https://nvd.nist.gov/vuln/detail/" + cve)
	if err != nil {
		return nil, err
	}

	retry := 0
	for {
		var cvss2 []*html.Node
		cvss3, _ := htmlquery.QueryAll(doc, "//*[@id=\"Cvss3NistCalculatorAnchor\"]")
		if cvss3 == nil {
			cvss3, _ = htmlquery.QueryAll(doc, "//*[@id=\"Cvss3CnaCalculatorAnchor\"]")
		}
		if cvss3 == nil {
			cvss2, _ = htmlquery.QueryAll(doc, "//*[@id=\"Cvss2CalculatorAnchor\"]")
		}
		if cvss2 != nil {
			cvss3 = cvss2
		}
		if cvss3 == nil {
			return nil, fmt.Errorf(cve + " cannot get score")
		}
		score := strings.Split(cvss3[0].FirstChild.Data, " ")
		if score == nil {
			retry++
			if retry >= 3 {
				return nil, fmt.Errorf(cve + " cannot get score")
			}
			continue
		}
		return score, nil
	}

}

func main() {
	flag.StringVar(&vulnDbDir, "vulndb-dir", "", "vulndb path, yaml from https://github.com/golang/vulndb/tree/master/reports")
	flag.Parse()

	if vulnDbDir == "" {
		fmt.Fprintf(os.Stderr, "`-vulndb-dir` param is require\n")
		os.Exit(1)
	}
	//var vdb config.VulnDb
	vulnDBs, err := file.ReadSuffixFiles(vulnDbDir, []string{"*.yaml", "*.yml"})
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for yamlPath, yamlContent := range vulnDBs {
		vuln := config.VulnDb{}
		err := yaml.Unmarshal([]byte(yamlContent), &vuln)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[skip] Error while parse "+yamlPath+": "+err.Error()+"\n")
			continue
		}

		vuln.Cvss3 = make([]string, 0)
		vuln.Severities = make([]string, 0)

		if len(vuln.Cves) == 0 {
			fmt.Fprintf(os.Stderr, "[skip] Error while get "+yamlPath+", it hasn't CVE \n")
			continue
		}

		for _, cve := range vuln.Cves {
			score, err := getCVEScore(cve)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[skip] Error while get "+cve+": "+err.Error()+"\n")
				continue
			}

			vuln.Cvss3 = append(vuln.Cvss3, score[0])
			vuln.Severities = append(vuln.Severities, score[1])
		}
		data, err := yaml.Marshal(vuln)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[skip] Error while marshal "+yamlPath+": "+err.Error()+"\n")
			continue
		}
		err = ioutil.WriteFile(yamlPath, data, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[skip] Error while write "+yamlPath+": "+err.Error()+"\n")
			continue
		}
	}

}
