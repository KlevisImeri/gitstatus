package walker

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitstatus/src/defaults"
	"gitstatus/src/git"
	"gitstatus/src/logger"
	"gitstatus/src/types"
)

func Walk(
	ctx context.Context,
	cfg types.Config,
	logger *logger.Logger,
	callback func(types.RepoResult),
) error {
	logger.Info("Starting scan from: %s", cfg.RootPath)
	err := filepath.WalkDir(cfg.RootPath, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			logger.Error("Error accessing path %s: %v", path, err)
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		if cfg.MaxDepth > 0 {
			rel, relErr := filepath.Rel(cfg.RootPath, path)
			if relErr != nil {
				logger.Warn("Could not calculate relative path for %s: %v", path, relErr)
				return nil
			}
			depth := strings.Count(rel, string(os.PathSeparator))
			if depth >= cfg.MaxDepth {
				return filepath.SkipDir
			}
		}

		for _, skipDir := range defaults.DefaultIgnoredDirs {
			if d.Name() == skipDir {
				logger.Debug("Skipping directory: %s", path)
				return filepath.SkipDir
			}
		}

		logger.Debug("Checking directory: %s", path)

		gitDir := filepath.Join(path, ".git")
		fileInfo, statErr := os.Stat(gitDir)
		if statErr == nil {
			if fileInfo.IsDir() {
				logger.Debug("Found git repo: %s", path)

				result, repoErr := git.GetRepoStatus(ctx, path, logger, cfg)
				if repoErr != nil {
					logger.Error("Error getting repo status for %s: %v", path, repoErr)
					callback(types.RepoResult{Path: path, Error: repoErr})
				} else {
					callback(*result)
				}
			}
		} else {
			if !os.IsNotExist(statErr) {
				logger.Debug("Error checking for .git in %s: %v", path, statErr)
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("WalkDir returned error: %v", err)
	}

	return err
}
