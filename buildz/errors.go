package buildz

import "errors"

var (
	ErrGOPATHNotSet         = errors.New("GOPATH is not set")
	ErrPathIsNotInGOPATH    = errors.New("path is not in GOPATH")
	ErrReachedRootDirectory = errors.New("reached root directory without finding valid package import path")
)
