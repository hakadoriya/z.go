package buildz

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

func IsInGOPATH(path string) (bool, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return false, ErrGOPATHNotSet
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("os.Getwd: %w", err)
	}

	// Check if the current directory is in GOPATH
	for _, dir := range filepath.SplitList(gopath) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(currentDir, absDir) {
			return true, nil
		}
	}

	return false, nil
}
