package buildz

import "errors"

var (
	ErrPathIsNotDirectory   = errors.New("path is not a directory")
	ErrReachedRootDirectory = errors.New("reached root directory without finding valid package import path")
	ErrModulePathNotFound   = errors.New("module path not found")
	ErrGoModFileNotFound    = errors.New("go.mod file not found")
)
