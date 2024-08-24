package buildinfoz

import (
	"runtime/debug"
)

// Other than the value to embed with -ldflags
//
//nolint:gochecknoglobals
var (
	buildCGOEnabled string
	debugBuildInfo  *debug.BuildInfo
)

// CGOEnabled returns the value of CGO_ENABLED set at build time.
func CGOEnabled() string { return buildCGOEnabled }

// GoVersion returns the Go version used for the build.
func GoVersion() string { return debugBuildInfo.GoVersion }
