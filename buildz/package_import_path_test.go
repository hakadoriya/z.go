package buildz

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/internal/consts"
)

func TestFindPackageImportPath(t *testing.T) {
	t.Parallel()

	if isInGOPATH, err := IsInGOPATH("."); err != nil {
		t.Skipf("🚫: SKIP: This is a workaround to skip this test when not in GOPATH (mainly for GitHub Actions): %v", err)
	} else if !isInGOPATH {
		t.Skip("🚫: SKIP: This is a workaround to skip this test when not in GOPATH (mainly for GitHub Actions)")
	}

	t.Run("success,testdata/testdata", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("testdata/testdata")
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		expected := path.Join(consts.ModuleName, "buildz", "testdata/testdata")
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("success,.", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath(".")
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		expected := path.Join(consts.ModuleName, "buildz")
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,no-such-file-or-directory", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("no-such-file-or-directory")
		if err == nil || !strings.Contains(err.Error(), `build.ImportDir: path=no-such-file-or-directory: cannot find package "`) {
			t.Errorf("❌: !errors.Is(err, ErrReachedRootDirectory): %+v", err)
		}

		expected := ""
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,/", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("/")
		if !errors.Is(err, ErrPathIsNotInGOPATH) {
			t.Errorf("❌: !errors.Is(err, ErrPathIsNotInGOPATH): %+v", err)
		}

		expected := ""
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("success,findPackageImportPath,testdata", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{FilepathAbsFunc: func(s string) (string, error) {
			return "", os.ErrInvalid
		}}

		actual, err := findPackageImportPath(iface, "testdata")
		if !errors.Is(err, os.ErrInvalid) {
			t.Errorf("❌: !errors.Is(err, os.ErrInvalid): %+v", err)
		}

		expected := ""
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})
}
