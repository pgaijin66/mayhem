package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
	Tag       = "unknown"
	GoVersion = runtime.Version()
)

// showVersion prints version information
func showVersion() {
	fmt.Printf(`phailure %s
Git Commit: %s
Build Date: %s
Go Version: %s
Tag:        %s
Platform:   %s/%s
`, Version, GitCommit, BuildDate, GoVersion, Tag, runtime.GOOS, runtime.GOARCH)
}
