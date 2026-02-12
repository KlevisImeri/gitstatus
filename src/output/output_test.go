package output

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitstatus/src/git"
	"gitstatus/src/logger"
	"gitstatus/src/types"
)

// Unit tests for formatting logic (independent of git)

func TestFormatBranchLineAheadOnly(t *testing.T) {
	b := types.BranchSyncStatus{
		Name:       "main",
		Current:    false,
		Ahead:      3,
		Behind:     0,
		Gone:       false,
		NoUpstream: false,
	}

	result := formatBranchLine("/repo", b, false)

	if !strings.Contains(result, "main") {
		t.Error("Expected result to contain branch name")
	}
	if !strings.Contains(result, "ahead 3") {
		t.Error("Expected result to contain 'ahead 3'")
	}
	if !strings.HasPrefix(result, ColorGreen) {
		t.Errorf("Expected green color for ahead only, got: %s", result)
	}
	if !strings.HasSuffix(result, ColorReset) {
		t.Error("Expected result to end with color reset")
	}
}

func TestFormatBranchLineBehindOnly(t *testing.T) {
	b := types.BranchSyncStatus{
		Name:       "feature",
		Current:    false,
		Ahead:      0,
		Behind:     5,
		Gone:       false,
		NoUpstream: false,
	}

	result := formatBranchLine("/repo", b, false)

	if !strings.Contains(result, "feature") {
		t.Error("Expected result to contain branch name")
	}
	if !strings.Contains(result, "behind 5") {
		t.Error("Expected result to contain 'behind 5'")
	}
	if !strings.HasPrefix(result, ColorRed) {
		t.Errorf("Expected red color for behind only, got: %s", result)
	}
}

func TestFormatBranchLineGone(t *testing.T) {
	b := types.BranchSyncStatus{
		Name:       "old-branch",
		Current:    false,
		Ahead:      0,
		Behind:     0,
		Gone:       true,
		NoUpstream: false,
	}

	result := formatBranchLine("/repo", b, false)

	if !strings.Contains(result, "old-branch") {
		t.Error("Expected result to contain branch name")
	}
	if !strings.Contains(result, "gone") {
		t.Error("Expected result to contain 'gone'")
	}
	if !strings.HasPrefix(result, ColorMagenta) {
		t.Errorf("Expected magenta color for gone branch, got: %s", result)
	}
}

// Integration tests using real repos

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintResultsReal(t *testing.T) {
	testEnv := SetupTestRepos(t)
	logger, _ := logger.NewLogger([]string{}, "")
	ctx := context.Background()

	t.Run("CleanRepo_ShowAllFalse", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_synced")
		result, _ := git.GetRepoStatus(ctx, repoPath, logger, types.Config{})

		cfg := types.Config{ShowAll: false, NoColor: true}

		output := captureOutput(func() {
			PrintResults([]types.RepoResult{*result}, cfg, logger)
		})

		if strings.Contains(output, "repo_synced") {
			t.Error("Should not show clean repo when ShowAll=false")
		}
	})

	t.Run("CleanRepo_ShowAllTrue", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_synced")
		result, _ := git.GetRepoStatus(ctx, repoPath, logger, types.Config{ShowAll: true})

		cfg := types.Config{ShowAll: true, NoColor: true}

		output := captureOutput(func() {
			PrintResults([]types.RepoResult{*result}, cfg, logger)
		})

		if !strings.Contains(output, "repo_synced") {
			t.Error("Should show clean repo when ShowAll=true")
		}
	})

	t.Run("RepoAhead", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_ahead")
		result, _ := git.GetRepoStatus(ctx, repoPath, logger, types.Config{})

		cfg := types.Config{ShowAll: false, NoColor: true}

		output := captureOutput(func() {
			PrintResults([]types.RepoResult{*result}, cfg, logger)
		})

		if !strings.Contains(output, "repo_ahead") {
			t.Error("Should show repo_ahead")
		}
		if !strings.Contains(output, "ahead 1") {
			t.Error("Should show 'ahead 1'")
		}
	})

	t.Run("RepoBehind", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_behind")
		result, _ := git.GetRepoStatus(ctx, repoPath, logger, types.Config{})

		cfg := types.Config{ShowAll: false, NoColor: true}

		output := captureOutput(func() {
			PrintResults([]types.RepoResult{*result}, cfg, logger)
		})

		if !strings.Contains(output, "repo_behind") {
			t.Error("Should show repo_behind")
		}
		if !strings.Contains(output, "behind 1") {
			t.Error("Should show 'behind 1'")
		}
	})

	t.Run("RepoNoUpstream", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_no_upstream")
		result, _ := git.GetRepoStatus(ctx, repoPath, logger, types.Config{})

		cfg := types.Config{ShowAll: false, NoColor: true}

		output := captureOutput(func() {
			PrintResults([]types.RepoResult{*result}, cfg, logger)
		})

		if !strings.Contains(output, "repo_no_upstream") {
			t.Error("Should show repo_no_upstream")
		}
		if !strings.Contains(output, "no upstream") {
			t.Error("Should show 'no upstream'")
		}
	})
}
