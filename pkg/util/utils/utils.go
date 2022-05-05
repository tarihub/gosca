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

// ExcelDiv get the correct column char from the excel column number
// https://blog.csdn.net/random_w/article/details/88431166
func ExcelDiv(num int) string {
	var (
		Str string = ""
		k   int
		// Save the value of each bit of converted data, and then match a-Z by index
		temp []int
	)
	// The characters a-z used to match the corresponding number
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	// Data larger than 26 needs to be split
	if num > 26 {
		for {
			k = num % 26
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			num = (num - k) / 26
			if num <= 26 {
				temp = append(temp, num)
				break
			}
		}
	} else {
		return Slice[num]
	}

	for _, value := range temp {
		// Since the data is stored in reverse order after being shred, Str comes after
		Str = Slice[value] + Str
	}
	return Str
}
