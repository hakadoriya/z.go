package buildz

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

func FindPackageImportPath(dir string) (string, error) {
	return findPackageImportPath(&pkg{
		BuildImportDirFunc: build.ImportDir,
		FilepathAbsFunc:    filepath.Abs,
		FilepathRelFunc:    filepath.Rel,
		OSReadFileFunc:     os.ReadFile,
		OSStatFunc:         os.Stat,
	}, dir)
}

func findPackageImportPath(iface pkgInterface, dir string) (string, error) {
	dirStat, err := iface.OSStat(dir)
	if err != nil {
		return "", fmt.Errorf("OStat: path=%s: %w", dir, err)
	}

	if !dirStat.IsDir() {
		return "", fmt.Errorf("path=%s: %w", dir, ErrPathIsNotDirectory)
	}

	path, err1 := findPackageImportPath1(iface, dir)
	if err1 == nil {
		return path, nil
	}

	path, err2 := findPackageImportPath2(iface, dir)
	if err2 == nil {
		return path, nil
	}

	return "", fmt.Errorf("findPackageImportPath: err1=%w, err2=%w", err1, err2)
}

func findPackageImportPath1(iface pkgInterface, dir string) (string, error) {
	currentDir := dir
	var accumulatedPath string

	for {
		absDir, err := iface.FilepathAbs(currentDir)
		if err != nil {
			return "", fmt.Errorf("filepath.Abs: path=%s: %w", currentDir, err)
		}

		pkg, err := iface.BuildImportDir(absDir, build.FindOnly)
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

// findPackageImportPath2 traverses up from the current directory,
// and if a go.mod file is found, returns the relative path from that directory.
// Returns an error if go.mod is not found.
func findPackageImportPath2(iface pkgInterface, dir string) (string, error) {
	absDir, err := iface.FilepathAbs(dir)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: path=%s: %w", dir, err)
	}

	for dir := absDir; ; {
		// Check if go.mod file exists
		goModPath := filepath.Join(dir, "go.mod")
		fileInfo, err := iface.OSStat(goModPath)
		if err == nil && !fileInfo.IsDir() {
			// Found go.mod file, calculate relative path and concat Go Module Name
			relPath, err := iface.FilepathRel(dir, absDir)
			if err != nil {
				return "", fmt.Errorf("filepath.Rel: %w", err)
			}

			goMod, err := iface.OSReadFile(goModPath)
			if err != nil {
				return "", fmt.Errorf("os.ReadFile: %w", err)
			}

			// read module name line
			modulePath, err := extractModuleName(goMod)
			if err != nil {
				return "", fmt.Errorf("extractModuleName: %w", err)
			}

			return filepath.Join(modulePath, relPath), nil
		}

		if !os.IsNotExist(err) {
			return "", fmt.Errorf("os.Stat: %w", err)
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory without finding go.mod
			break
		}
		dir = parent
	}

	return "", ErrGoModFileNotFound
}

func extractModuleName(goMod []byte) (string, error) {
	lines := strings.Split(string(goMod), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", ErrModulePathNotFound
}

// pkgInterface is a entry point for mocking.
type pkgInterface interface {
	FilepathAbs(path string) (string, error)
	FilepathRel(basepath, targpath string) (string, error)
	OSReadFile(path string) ([]byte, error)
	OSStat(path string) (os.FileInfo, error)
	BuildImportDir(dir string, mode build.ImportMode) (*build.Package, error)
}

type pkg struct {
	FilepathAbsFunc    func(path string) (string, error)
	FilepathRelFunc    func(basepath, targpath string) (string, error)
	OSReadFileFunc     func(path string) ([]byte, error)
	OSStatFunc         func(path string) (os.FileInfo, error)
	BuildImportDirFunc func(dir string, mode build.ImportMode) (*build.Package, error)
}

func (s *pkg) BuildImportDir(dir string, mode build.ImportMode) (*build.Package, error) {
	return s.BuildImportDirFunc(dir, mode)
}

func (s *pkg) FilepathAbs(path string) (string, error) {
	return s.FilepathAbsFunc(path)
}

func (s *pkg) FilepathRel(basepath, targpath string) (string, error) {
	return s.FilepathRelFunc(basepath, targpath)
}

func (s *pkg) OSReadFile(path string) ([]byte, error) {
	return s.OSReadFileFunc(path)
}

func (s *pkg) OSStat(path string) (os.FileInfo, error) {
	return s.OSStatFunc(path)
}
