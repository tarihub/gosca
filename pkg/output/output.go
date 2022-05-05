package output

import (
	"github.com/TARI0510/gosca/pkg/config"
)

type F struct {
	OutFormat string
	OutPath   string
	WorkDir   string
}

type OutBuilder func(vulIdList []string, vulnDbMap map[string]config.VulnDb, stdPkgLocationMap map[string][]string, f F)
