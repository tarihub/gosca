package config

import "github.com/TARI0510/gosca/pkg/util/utils"

type Dependence struct {
	Module   string
	Versions []string
}

func (d Dependence) GetMatches(vulDbIdxMap VulDbIdxMap) ([]string, bool) {
	var matches []string
	var matchesMap = make(map[string]bool)

	for _, vul := range vulDbIdxMap[d.Module] {
	VersionLoop:
		for _, version := range vul.Versions {
			if introVersion, ok := version["introduced"]; ok {
				for _, depsVersion := range d.Versions {
					if utils.CompareVersion(depsVersion, introVersion) < 0 {
						continue VersionLoop
					}
				}
			}

			if fixedVersion, ok := version["fixed"]; ok {
				for _, depsVersion := range d.Versions {
					if utils.CompareVersion(depsVersion, fixedVersion) > -1 {
						continue VersionLoop
					}
				}
			}

			if matchesMap[vul.Id] {
				break
			}
			matchesMap[vul.Id] = true
			matches = append(matches, vul.Id)
		}
	}

	if matches == nil {
		return nil, false
	}

	return matches, true
}
