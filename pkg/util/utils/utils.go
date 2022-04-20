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
		// '-' is one of the go version delimiters. It will be followed by a false version number.
		//Skip it if you encounter it; otherwise, misjudgment will occur
		for ; i < n && version1[i] != '.' && version1[i] != '-'; i++ {
			x = x*10 + int(version1[i]-'0')
		}
		i++ // Skip the dot and dash
		y := 0
		for ; j < m && version2[j] != '.' && version2[j] != '-'; j++ {
			y = y*10 + int(version2[j]-'0')
		}
		j++ // Skip the dot and dash
		if x > y {
			return 1
		}
		if x < y {
			return -1
		}
	}
	return 0
}

// RemoveElemDisorder remove element from slice disorderly.
//func RemoveElemDisorder(s []string, i int) []string {
//	if i == len(s) {
//		return s[:len(s)-1]
//	}
//	s[i] = s[len(s)-1]
//	return s[:len(s)-1]
//}

//RemoveElem remove element from slice.
func RemoveElem(s []string, e string) []string {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
