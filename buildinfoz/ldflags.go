package buildinfoz

// Values to embed with -ldflags
//
//nolint:gochecknoglobals
var (
	buildVersion   string
	buildRevision  string
	buildBranch    string
	buildTimestamp string
)

// BuildVersion returns the value embedded with -ldflags at build time.
//
// Expected -ldflags value:
//
//	-X github.com/hakadoriya/z.go/buildinfoz.buildVersion=`git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD`
func BuildVersion() string { return buildVersion }

// BuildRevision returns the value embedded with -ldflags at build time.
//
// Expected -ldflags value:
//
//	-X github.com/hakadoriya/z.go/buildinfoz.buildRevision=`git rev-parse HEAD`
func BuildRevision() string { return buildRevision }

// BuildBranch returns the value embedded with -ldflags at build time.
//
// Expected -ldflags value:
//
//	-X github.com/hakadoriya/z.go/buildinfoz.buildBranch=`git rev-parse --abbrev-ref HEAD | tr / -`
func BuildBranch() string { return buildBranch }

// BuildTimestamp returns the value embedded with -ldflags at build time.
//
// Expected -ldflags value:
//
//	-X github.com/hakadoriya/z.go/buildinfoz.buildTimestamp=`git log -n 1 --format='%cI'`"
func BuildTimestamp() string { return buildTimestamp }
