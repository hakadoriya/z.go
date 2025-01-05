package buildz

import (
	"fmt"
	"go/build"
	"path/filepath"
	"strings"
)

func IsInGOPATH(path string) (bool, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return false, ErrGOPATHNotSet
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("filepath.Abs: path=%s: %w", path, err)
	}

	// Check if the current directory is in GOPATH
	for _, dir := range filepath.SplitList(gopath) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absDir) {
			return true, nil
		}
	}

	return false, nil
}
