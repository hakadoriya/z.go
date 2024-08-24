package buildinfoz

import (
	"runtime"
	"runtime/debug"
)

// Exceptionally, initialize the value in init considering the coexistence with -ldflags
//
//nolint:gochecknoinits,cyclop
func init() {
	var ok bool
	debugBuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		//nolint:exhaustruct
		debugBuildInfo = &debug.BuildInfo{GoVersion: runtime.Version()}
	}

	// If the value is not set by -ldflags, set the value obtained from BuildInfo as a fallback
	for _, s := range debugBuildInfo.Settings {
		switch s.Key {
		case "vcs.revision":
			if buildVersion == "" {
				buildVersion = s.Value
			}
			if buildRevision == "" {
				buildRevision = s.Value
			}
		case "vcs.branch":
			if buildBranch == "" {
				buildBranch = s.Value
			}
		case "vcs.time":
			if buildTimestamp == "" {
				buildTimestamp = s.Value
			}
		case "CGO_ENABLED":
			buildCGOEnabled = s.Value
		}
	}
}
