package version

import "fmt"

// These variables are given values at build time with ldflags parameters
// to the go compiler. See the goreleaser config (/.goreleaser.yml) for more
// detail.

// Version field is a SemVer that should indicate the baked-in version
// of the CLI
var Version = "unset"

// Commit is the commit hash of the source version used to build this binary
var Commit = "unset"

// BuildNumber field is the particular CI build number for this particular
// version.
var BuildNumber = "unset"

// FullVersionName is the user-visible aggregation of version and tag
// of this codebase and the build number that produced it.
var FullVersionName = fmt.Sprintf("%s-%s (Build %s)", Version, Commit, BuildNumber)
