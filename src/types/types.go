package types

// BranchSyncStatus represents a branch's sync state with origin
type BranchSyncStatus struct {
	Name       string
	Current    bool // is checked out?
	Ahead      int  // commits ahead of origin
	Behind     int  // commits behind origin
	Gone       bool // remote branch is gone
	NoUpstream bool // no upstream configured
}

// WorkdirStatus represents uncommitted changes in the working directory
type WorkdirStatus struct {
	Modified  int // files modified but not staged
	Staged    int // files staged (added to index)
	Untracked int // untracked files
}

// RepoResult holds info about a git repository
type RepoResult struct {
	Path           string
	Branches       []BranchSyncStatus // branches relevant to status (unsynced or all depending on config)
	HasUnsynced    bool               // true if any branch is ahead/behind/gone
	Uncommitted    WorkdirStatus      // uncommitted changes in working directory
	HasUncommitted bool               // true if there are uncommitted changes
	Error          error              // any error encountered
}

// Config holds CLI configuration
type Config struct {
	RootPath string
	MaxDepth int
	LogTypes []string
	ShowAll  bool
	NoColor  bool
	LogFile  string
}
