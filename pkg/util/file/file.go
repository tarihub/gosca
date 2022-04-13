package file

import (
	"io/ioutil"
	"path"
	"strings"
)

// ListSuffixFiles list directory's files with specify suffix
func ListSuffixFiles(dirPath string, suffixes []string) ([]string, error) {
	var files []string
	dirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	//
	for _, file := range dirs {
		for _, suffix := range suffixes {
			if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
				files = append(files, path.Join(dirPath, file.Name()))
			}
		}
	}
	return files, nil
}
