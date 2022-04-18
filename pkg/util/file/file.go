package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ListSuffixFiles list directory's files with specify suffix
func ListSuffixFiles(dirPath string, suffixes []string) ([]string, error) {
	var files []string

	for _, suffix := range suffixes {
		globs, err := filepath.Glob(filepath.Join(dirPath, suffix))
		if err != nil {
			return nil, err
		}
		files = append(files, globs...)
	}

	return files, nil
}

// PackagePaths returns a slice with all packages path at given root directory
func PackagePaths(root string, excludes []*regexp.Regexp) ([]string, error) {
	paths := map[string]bool{}
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		// FIXME cannot identify soft link
		if filepath.Ext(path) == ".go" {
			path = filepath.Dir(path)
			if isExcluded(filepath.ToSlash(path), excludes) {
				return nil
			}
			paths[path] = true
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	var result []string
	for p := range paths {
		absPath, err := GetPkgAbsPath(p)
		if err != nil {
			fmt.Println("[skip] Path " + absPath + " doesn't exist.")
			continue
		}
		result = append(result, absPath)
	}
	return result, nil
}

// isExcluded checks if a string matches any of the exclusion regexps
func isExcluded(str string, excludes []*regexp.Regexp) bool {
	if excludes == nil {
		return false
	}
	for _, exclude := range excludes {
		if exclude != nil && exclude.MatchString(str) {
			return true
		}
	}
	return false
}

// ExcludedDirsRegExp builds the regexps for a list of excluded dirs provided as strings
func ExcludedDirsRegExp(excludedDirs []string) []*regexp.Regexp {
	var exps []*regexp.Regexp
	for _, excludedDir := range excludedDirs {
		str := fmt.Sprintf(`([\\/])?%s([\\/])?`, strings.ReplaceAll(filepath.ToSlash(excludedDir), "/", `\/`))
		r := regexp.MustCompile(str)
		exps = append(exps, r)
	}
	return exps
}

// GetPkgAbsPath returns the Go package absolute path derived from
// the given path
func GetPkgAbsPath(pkgPath string) (string, error) {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", errors.New("no project absolute path found")
	}
	return absPath, nil
}
