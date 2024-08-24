package buildz

import (
	"go/build"
	"path/filepath"
	"strings"
)

func IsInGOPATH(path string) (bool, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return false, ErrGOPATHNotSet
	}

	// Check if the current directory is in GOPATH
	for _, dir := range filepath.SplitList(gopath) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(path, absDir) {
			return true, nil
		}
	}

	return false, nil
}
