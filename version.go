package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
	tag       = "unknown"
	GoVersion = runtime.Version()
)

// showVersion prints version information
func showVersion() {
	fmt.Printf(`mayhem %s
Git Commit: %s
Build Date: %s
Go Version: %s
Tag : %s
Platform:   %s/%s
`, Version, GitCommit, BuildDate, GoVersion, runtime.GOOS, tag, runtime.GOARCH)
}
