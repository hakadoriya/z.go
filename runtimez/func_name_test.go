package runtimez

import (
	"path"
	"testing"
)

func TestFuncName(t *testing.T) {
	t.Parallel()
	if expected, actual := "runtimez.TestFuncName", wrapFuncName(); actual != expected {
		t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
	}
}

func wrapFuncName() string {
	return FuncName(WithFuncNameOptionAddCallerSkip(1))
}

func TestFullFuncName(t *testing.T) {
	t.Parallel()
	if expected, actual := path.Join("github.com/hakadoriya/z.go", "runtimez.TestFullFuncName"), wrapFullFuncName(); actual != expected {
		t.Errorf("❌: expected(%q) != actual(%q)", expected, actual)
	}
}

func wrapFullFuncName() string {
	return FullFuncName(WithFuncNameOptionAddCallerSkip(1))
}
