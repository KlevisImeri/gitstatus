package defaults

// DirsToSkip contains directories that should be skipped during traversal.
// You can add or remove directories here to customize the scanning process.
var DefaultIgnoredDirs = []string{
	".git",
	"node_modules",
	"vendor",
	".idea",
	".vscode",
	"dist",
	"build",
	"target",
	"__pycache__",
	".sass-cache",
}

// DefaultLogFile is the default name for the log file
const DefaultLogFile = "gitstatus.log"

// DefaultMaxDepth is the default maximum depth for directory traversal (0 = unlimited)
const DefaultMaxDepth = 0

// DefaultGitCommandTimeoutSeconds is the timeout in seconds for git commands
const DefaultGitCommandTimeoutSeconds = 5
