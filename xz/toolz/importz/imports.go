package importz

import (
	"fmt"
	"go/format"

	"golang.org/x/tools/imports"
)

func Format(src []byte) ([]byte, error) {
	return internalFormat(&pkg{
		FormatSourceFunc:   format.Source,
		ImportsProcessFunc: imports.Process,
	}, src)
}

func internalFormat(iface pkgInterface, src []byte) ([]byte, error) {
	// Format the contents
	formatted, err := iface.FormatSource(src)
	if err != nil {
		return nil, fmt.Errorf("format.Source: %w", err)
	}

	const filename = ""
	imported, err := iface.ImportsProcess(filename, formatted, (*imports.Options)(nil))
	if err != nil {
		return nil, fmt.Errorf("imports.Process: %w", err)
	}

	return imported, nil
}

type pkgInterface interface {
	FormatSource(src []byte) ([]byte, error)
	ImportsProcess(filename string, src []byte, opt *imports.Options) ([]byte, error)
}

type pkg struct {
	FormatSourceFunc   func(src []byte) ([]byte, error)
	ImportsProcessFunc func(filename string, src []byte, opt *imports.Options) ([]byte, error)
}

func (s *pkg) FormatSource(src []byte) ([]byte, error) {
	return s.FormatSourceFunc(src)
}

func (s *pkg) ImportsProcess(filename string, src []byte, opt *imports.Options) ([]byte, error) {
	return s.ImportsProcessFunc(filename, src, opt)
}
