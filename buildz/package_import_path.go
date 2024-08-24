package buildz

import (
	"errors"
	"fmt"
	"go/build"
	"path/filepath"
)

//unc PackageImportPath(dir string) (string, error) {
//	absDir, err := filepath.Abs(dir)
//	if err != nil {
//		return "", fmt.Errorf("filepath.Abs: path=%s: %w", dir, err)
//	}
//
//	pkg, err := build.ImportDir(absDir, build.FindOnly)
//	if err != nil {
//		return "", fmt.Errorf("build.ImportDir: path=%s: %w", dir, err)
//	}
//
//	if pkg.ImportPath == "." {
//		// If ImportPath is ".", find the parent package recursively.
//		basename := filepath.Base(absDir)
//		parentDir := filepath.Dir(absDir)
//		parent, err := PackageImportPath(parentDir)
//		if err != nil {
//			return "", fmt.Errorf("DetectPackageImportPath: %w", err)
//		}
//
//		return filepath.Join(parent, basename), nil
//	}
//
//	return pkg.ImportPath, nil
//

var ErrReachedRootDirectory = errors.New("reached root directory without finding valid package import path")

func FindPackageImportPath(dir string) (string, error) {
	return findPackageImportPath(&pkg{FilepathAbsFunc: filepath.Abs}, dir)
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
