package config

import "github.com/TARI0510/gosca/pkg/util/utils"

type Dependence struct {
	Module   string
	Versions []string
}

func (d Dependence) GetMatches(vulDbIdxMap VulDbIdxMap) (string, bool) {
	if _, ok := vulDbIdxMap[d.Module]; !ok {
		return "", false
	}

VersionLoop:
	for _, version := range vulDbIdxMap[d.Module].Versions {
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

		// FIXME go 标准库漏洞只返回一个..
		return vulDbIdxMap[d.Module].Id, true
	}
	return "", false
}
