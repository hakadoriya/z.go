package buildz

import (
	"errors"
	"fmt"
	"go/build"
	"path/filepath"
)

var ErrReachedRootDirectory = errors.New("reached root directory without finding valid package import path")

func FindPackageImportPath(dir string) (string, error) {
	return findPackageImportPath(&pkg{FilepathAbsFunc: filepath.Abs}, dir)
}

func findPackageImportPath(iface pkgInterface, dir string) (string, error) {
	currentDir := dir
	var accumulatedPath string

	for {
		absDir, err := iface.FilepathAbs(currentDir)
		if err != nil {
			return "", fmt.Errorf("filepath.Abs: path=%s: %w", currentDir, err)
		}

		pkg, err := build.ImportDir(absDir, build.FindOnly)
		if err != nil {
			return "", fmt.Errorf("build.ImportDir: path=%s: %w", currentDir, err)
		}

		if pkg.ImportPath != "." {
			if accumulatedPath == "" {
				return pkg.ImportPath, nil
			}
			return filepath.Join(pkg.ImportPath, accumulatedPath), nil
		}

		// If ImportPath is ".", prepare to check the parent directory
		parentDir := filepath.Dir(absDir)

		// If we've reached the root directory, we can't go up any further
		if parentDir == absDir {
			return "", fmt.Errorf("path=%s: %w", dir, ErrReachedRootDirectory)
		}

		// Prepare for the next iteration
		if basename := filepath.Base(absDir); accumulatedPath == "" {
			accumulatedPath = basename
		} else {
			accumulatedPath = filepath.Join(basename, accumulatedPath)
		}

		currentDir = parentDir
	}
}

// pkgInterface is a entry point for mocking.
type pkgInterface interface {
	FilepathAbs(string) (string, error)
}

type pkg struct {
	FilepathAbsFunc func(string) (string, error)
}

func (s *pkg) FilepathAbs(path string) (string, error) {
	return s.FilepathAbsFunc(path)
}
