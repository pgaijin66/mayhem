package version

import "fmt"

var (
	// Version is the current version of the application
	// This will be set at build time using ldflags
	Version = "dev"

	// BuildTime is when the binary was built
	// This will be set at build time using ldflags
	BuildTime = "unknown"

	// GitCommit is the git commit hash
	// This will be set at build time using ldflags
	GitCommit = "unknown"
)

// Print displays version information
func Print() {
	fmt.Printf(`
	███╗   ███╗ █████╗ ██╗   ██╗██╗  ██╗███████╗███╗   ███╗
	████╗ ████║██╔══██╗╚██╗ ██╔╝██║  ██║██╔════╝████╗ ████║
	██╔████╔██║███████║ ╚████╔╝ ███████║█████╗  ██╔████╔██║
	██║╚██╔╝██║██╔══██║  ╚██╔╝  ██╔══██║██╔══╝  ██║╚██╔╝██║
	██║ ╚═╝ ██║██║  ██║   ██║   ██║  ██║███████╗██║ ╚═╝ ██║
	╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝

🔥 Chaos Engineering Tool 🔥

Version:    %s
Build Time: %s
Git Commit: %s
`, Version, BuildTime, GitCommit)
}

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetBuildTime returns the build time
func GetBuildTime() string {
	return BuildTime
}

// GetGitCommit returns the git commit hash
func GetGitCommit() string {
	return GitCommit
}
