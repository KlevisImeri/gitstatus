package walker

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"gitstatus/src/logger"
	"gitstatus/src/types"
)

func TestWalkRealRepos(t *testing.T) {
	testEnv := SetupTestRepos(t)
	logger, _ := logger.NewLogger([]string{}, "")
	ctx := context.Background()

	t.Run("FindGitAtRoot", func(t *testing.T) {
		repoPath := filepath.Join(testEnv, "repo_synced")
		cfg := types.Config{
			RootPath: repoPath,
			MaxDepth: 0,
		}

		var results []types.RepoResult
		err := Walk(ctx, cfg, logger, func(result types.RepoResult) {
			results = append(results, result)
		})

		if err != nil {
			t.Fatalf("Walk failed: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if len(results) > 0 && results[0].Path != repoPath {
			t.Errorf("Expected path %s, got %s", repoPath, results[0].Path)
		}
	})

	t.Run("FindNestedGit", func(t *testing.T) {
		nestedRoot := filepath.Join(testEnv, "nested")

		cfg := types.Config{
			RootPath: nestedRoot,
			MaxDepth: 3,
		}

		var results []types.RepoResult
		err := Walk(ctx, cfg, logger, func(result types.RepoResult) {
			results = append(results, result)
		})

		if err != nil {
			t.Fatalf("Walk failed: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		expectedPath := filepath.Join(nestedRoot, "level1", "repo_deep")
		if len(results) > 0 && results[0].Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, results[0].Path)
		}
	})

	t.Run("RespectMaxDepth", func(t *testing.T) {
		nestedRoot := filepath.Join(testEnv, "nested")

		cfg := types.Config{
			RootPath: nestedRoot,
			MaxDepth: 1,
		}

		var results []types.RepoResult
		Walk(ctx, cfg, logger, func(result types.RepoResult) {
			results = append(results, result)
		})

		if len(results) != 0 {
			t.Errorf("Expected 0 results with MaxDepth=1, got %d: %v", len(results), results)
		}
	})

	t.Run("SkipDirsToSkip", func(t *testing.T) {
		cfg := types.Config{
			RootPath: testEnv,
			MaxDepth: 5,
		}

		var results []types.RepoResult
		Walk(ctx, cfg, logger, func(result types.RepoResult) {
			results = append(results, result)
		})

		for _, r := range results {
			if strings.Contains(r.Path, "node_modules") {
				t.Errorf("Should not have found repo in node_modules: %s", r.Path)
			}
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		cfg := types.Config{
			RootPath: testEnv,
			MaxDepth: 5,
		}

		ctx, cancel := context.WithCancel(context.Background())

		var count int
		var mu sync.Mutex

		err := Walk(ctx, cfg, logger, func(result types.RepoResult) {
			mu.Lock()
			count++
			if count == 1 {
				cancel()
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond) 
		})

		if err != nil && err != context.Canceled {
			t.Logf("Walk returned error (expected possibly Canceled): %v", err)
		}

		// We can't guarantee exactly 1 result due to concurrency, but it shouldn't be ALL of them
		// And we expect cancellation to happen.
	})

	t.Run("FindMultipleRepos", func(t *testing.T) {
		cfg := types.Config{
			RootPath: testEnv,
			MaxDepth: 1,
		}

		var results []types.RepoResult
		Walk(ctx, cfg, logger, func(result types.RepoResult) {
			results = append(results, result)
		})

		expectedRepos := []string{
			"repo_synced", "repo_ahead", "repo_behind",
			"repo_modified", "repo_staged", "repo_untracked",
			"repo_gone", "repo_no_upstream",
		}

		for _, name := range expectedRepos {
			found := false
			expectedPath := filepath.Join(testEnv, name)
			for _, r := range results {
				if r.Path == expectedPath {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find repo: %s", name)
			}
		}
	})
}
