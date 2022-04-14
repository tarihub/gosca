package config

type VulnDbIdx struct {
	Id       string
	Versions []map[string]string
}

type VulDbIdxMap map[string][]VulnDbIdx
