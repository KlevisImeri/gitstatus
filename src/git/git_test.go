package git

import (
	"context"
	"path/filepath"
	"testing"

	"gitstatus/src/logger"
)

func TestGetRepoStatusReal(t *testing.T) {
	testEnv := setupTestRepos(t)
	logger, _ := logger.NewLogger([]string{}, "")
	ctx := context.Background()

	tests := []struct {
		repoName        string
		wantHasUnsynced bool
		wantAhead       int
		wantBehind      int
		wantGone        bool
		wantNoUpstream  bool
		wantUncommitted bool
	}{
		{
			repoName:        "repo_synced",
			wantHasUnsynced: false,
			wantAhead:       0,
			wantBehind:      0,
			wantUncommitted: false,
		},
		{
			repoName:        "repo_ahead",
			wantHasUnsynced: true,
			wantAhead:       1,
			wantBehind:      0,
			wantUncommitted: false,
		},
		{
			repoName:        "repo_behind",
			wantHasUnsynced: true,
			wantAhead:       0,
			wantBehind:      1,
			wantUncommitted: false,
		},
		{
			repoName:        "repo_gone",
			wantHasUnsynced: true,
			wantGone:        true,
			wantUncommitted: false,
		},
		{
			repoName:        "repo_no_upstream",
			wantHasUnsynced: true,
			wantNoUpstream:  true,
			wantUncommitted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.repoName, func(t *testing.T) {
			repoPath := filepath.Join(testEnv, tt.repoName)
			result, err := GetRepoStatus(ctx, repoPath, logger)
			if err != nil {
				t.Fatalf("GetRepoStatus failed: %v", err)
			}

			if result.HasUnsynced != tt.wantHasUnsynced {
				t.Errorf("HasUnsynced = %v, want %v", result.HasUnsynced, tt.wantHasUnsynced)
			}

			if tt.wantHasUnsynced {
				foundBranch := false
				for _, b := range result.Branches {
					if b.Current {
						foundBranch = true
						if b.Ahead != tt.wantAhead {
							t.Errorf("Ahead = %d, want %d", b.Ahead, tt.wantAhead)
						}
						if b.Behind != tt.wantBehind {
							t.Errorf("Behind = %d, want %d", b.Behind, tt.wantBehind)
						}
						if b.Gone != tt.wantGone {
							t.Errorf("Gone = %v, want %v", b.Gone, tt.wantGone)
						}
						if b.NoUpstream != tt.wantNoUpstream {
							t.Errorf("NoUpstream = %v, want %v", b.NoUpstream, tt.wantNoUpstream)
						}
					}
				}
				if !foundBranch {
					t.Error("Expected to find a current branch")
				}
			}
		})
	}
}

func TestGetWorkdirStatusReal(t *testing.T) {
	testEnv := setupTestRepos(t)
	logger, _ := logger.NewLogger([]string{}, "")
	ctx := context.Background()

	tests := []struct {
		repoName      string
		wantModified  int
		wantStaged    int
		wantUntracked int
	}{
		{
			repoName:      "repo_synced",
			wantModified:  0,
			wantStaged:    0,
			wantUntracked: 0,
		},
		{
			repoName:      "repo_modified",
			wantModified:  1,
			wantStaged:    0,
			wantUntracked: 0,
		},
		{
			repoName:      "repo_staged",
			wantModified:  0,
			wantStaged:    1,
			wantUntracked: 0,
		},
		{
			repoName:      "repo_untracked",
			wantModified:  0,
			wantStaged:    0,
			wantUntracked: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.repoName, func(t *testing.T) {
			repoPath := filepath.Join(testEnv, tt.repoName)
			status, err := GetWorkdirStatus(ctx, repoPath, logger)
			if err != nil {
				t.Fatalf("GetWorkdirStatus failed: %v", err)
			}

			if status.Modified != tt.wantModified {
				t.Errorf("Modified = %d, want %d", status.Modified, tt.wantModified)
			}
			if status.Staged != tt.wantStaged {
				t.Errorf("Staged = %d, want %d", status.Staged, tt.wantStaged)
			}
			if status.Untracked != tt.wantUntracked {
				t.Errorf("Untracked = %d, want %d", status.Untracked, tt.wantUntracked)
			}
		})
	}
}
