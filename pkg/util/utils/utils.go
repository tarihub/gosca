package utils

// UniqueSlice returns a new slice with all duplicate elements removed.
func UniqueSlice(s []string) []string {
	var unique []string
	m := make(map[string]bool)

	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

// CompareVersion compares two version strings.
func CompareVersion(version1, version2 string) int {
	n, m := len(version1), len(version2)
	i, j := 0, 0
	for i < n || j < m {
		x := 0
		// `-` 是 go 版本分隔符之一, 比如后面会接伪版本号等, 遇到要跳过否则会出现误判
		for ; i < n && version1[i] != '.' && version1[i] != '-'; i++ {
			x = x*10 + int(version1[i]-'0')
		}
		i++ // 跳过点号
		y := 0
		for ; j < m && version2[j] != '.' && version2[j] != '-'; j++ {
			y = y*10 + int(version2[j]-'0')
		}
		j++ // 跳过点号
		if x > y {
			return 1
		}
		if x < y {
			return -1
		}
	}
	return 0
}
