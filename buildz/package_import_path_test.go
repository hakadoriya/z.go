package buildz

import (
	"errors"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hakadoriya/z.go/internal/consts"
)

func TestFindPackageImportPath(t *testing.T) {
	t.Parallel()

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

	t.Run("success,NotInGOPATH", func(t *testing.T) {
		t.Parallel()

		tempDir, err := os.MkdirTemp("", "buildz-test")
		if err != nil {
			t.Fatalf("❌: err != nil: %+v", err)
		}
		defer os.RemoveAll(tempDir)

		testGoModFile, err := os.Create(filepath.Join(tempDir, "go.mod"))
		if err != nil {
			t.Fatalf("❌: err != nil: %+v", err)
		}
		defer testGoModFile.Close()

		if _, err := testGoModFile.WriteString("module testdata/testdata"); err != nil {
			t.Fatalf("❌: err != nil: %+v", err)
		}

		actual, err := FindPackageImportPath(tempDir)
		if err != nil {
			t.Errorf("❌: err != nil: %+v", err)
		}

		if expected := "testdata/testdata"; expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,no-such-file-or-directory", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("no-such-file-or-directory")
		if expected := "no such file or directory"; err == nil || !strings.Contains(err.Error(), expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}

		expected := ""
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,ErrPathIsNotDirectory", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("testdata/testdata/testdata.go")
		if expected := ErrPathIsNotDirectory; !errors.Is(err, expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}

		expected := ""
		if expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,/", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("/")
		if !errors.Is(err, ErrReachedRootDirectory) {
			t.Errorf("❌: !errors.Is(err, ErrReachedRootDirectory): %+v", err)
		}

		if expected := ""; expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})

	t.Run("error,/tmp", func(t *testing.T) {
		t.Parallel()

		actual, err := FindPackageImportPath("/tmp")
		if !errors.Is(err, ErrReachedRootDirectory) {
			t.Errorf("❌: !errors.Is(err, ErrReachedRootDirectory): %+v", err)
		}

		if expected := ""; expected != actual {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
		}
	})
}

func Test_findPackageImportPath1(t *testing.T) {
	t.Parallel()

	t.Run("error,filepath.Abs", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: func(s string) (string, error) {
				return "", os.ErrInvalid
			},
			OSStatFunc: os.Stat,
		}

		_, err := findPackageImportPath1(iface, "testdata")
		if expected := os.ErrInvalid; !errors.Is(err, expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}
	})

	t.Run("error,build.ImportDir", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			BuildImportDirFunc: func(s string, mode build.ImportMode) (*build.Package, error) {
				return nil, os.ErrInvalid
			},
			FilepathAbsFunc: filepath.Abs,
			FilepathRelFunc: filepath.Rel,
			OSReadFileFunc:  os.ReadFile,
			OSStatFunc:      os.Stat,
		}

		_, err := findPackageImportPath1(iface, "testdata")
		if expected := os.ErrInvalid; !errors.Is(err, expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}
	})
}

func Test_findPackageImportPath2(t *testing.T) {
	t.Parallel()

	t.Run("error,filepath.Abs", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: func(s string) (string, error) {
				return "", os.ErrInvalid
			},
			OSStatFunc: os.Stat,
		}

		_, err := findPackageImportPath2(iface, "testdata")
		if expected := os.ErrInvalid; !errors.Is(err, expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}
	})

	t.Run("error,os.Stat", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: filepath.Abs,
			OSStatFunc:      func(s string) (os.FileInfo, error) { return nil, os.ErrInvalid },
		}

		_, err := findPackageImportPath2(iface, "testdata")
		if expected := os.ErrInvalid; !errors.Is(err, expected) {
			t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
		}
	})

	t.Run("error,filepath.Rel", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: filepath.Abs,
			FilepathRelFunc: func(s string, s2 string) (string, error) { return "", os.ErrInvalid },
			OSReadFileFunc:  os.ReadFile,
			OSStatFunc:      os.Stat,
		}

		_, err := findPackageImportPath2(iface, "testdata")
		if !errors.Is(err, os.ErrInvalid) {
			t.Errorf("❌: !errors.Is(err, os.ErrInvalid): %+v", err)
		}
	})

	t.Run("error,os.ReadFile", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: filepath.Abs,
			FilepathRelFunc: filepath.Rel,
			OSReadFileFunc:  func(s string) ([]byte, error) { return nil, os.ErrInvalid },
			OSStatFunc:      os.Stat,
		}

		_, err := findPackageImportPath2(iface, "testdata")
		if !errors.Is(err, os.ErrInvalid) {
			t.Errorf("❌: !errors.Is(err, os.ErrInvalid): %+v", err)
		}
	})

	t.Run("success,empty_go.mod", func(t *testing.T) {
		t.Parallel()

		iface := &pkg{
			FilepathAbsFunc: filepath.Abs,
			FilepathRelFunc: filepath.Rel,
			OSReadFileFunc:  os.ReadFile,
			OSStatFunc:      os.Stat,
		}

		tempDir, err := os.MkdirTemp("", "buildz-test")
		if err != nil {
			t.Fatalf("❌: err != nil: %+v", err)
		}
		defer os.RemoveAll(tempDir)

		testGoModFile, err := os.Create(filepath.Join(tempDir, "go.mod"))
		if err != nil {
			t.Fatalf("❌: err != nil: %+v", err)
		}
		defer testGoModFile.Close()

		{
			_, err := findPackageImportPath2(iface, tempDir)
			if expected := ErrModulePathNotFound; !errors.Is(err, expected) {
				t.Errorf("❌: expected(%q) != actual(%q)", expected, err)
			}
		}
	})
}
