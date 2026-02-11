package output

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitstatus/src/logger"
	"gitstatus/src/types"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorCyan    = "\033[36m"
	ColorMagenta = "\033[35m"
)

func PrintResults(results []types.RepoResult, cfg types.Config, logger *logger.Logger) {
	hasIssues := false
	for _, res := range results {
		if res.HasUnsynced || res.HasUncommitted {
			hasIssues = true
			break
		}
	}

	if !hasIssues && !cfg.ShowAll {
		fmt.Println("No git repositories with unsynced status or uncommitted changes found.")
		return
	}

	for _, res := range results {
		if res.Error != nil {
			logger.Error("Error in repository %s: %v", res.Path, res.Error)
			continue
		}

		if !res.HasUnsynced && !res.HasUncommitted && !cfg.ShowAll {
			continue
		}

		for _, b := range res.Branches {
			line := formatBranchLine(res.Path, b, cfg.NoColor)
			fmt.Println(line)
		}

		if res.HasUncommitted {
			line := formatWorkdirLine(res.Path, res.Uncommitted, cfg.NoColor)
			fmt.Println(line)
		}

		if cfg.ShowAll && !res.HasUnsynced && !res.HasUncommitted {
			line := formatCleanRepoLine(res.Path, cfg.NoColor)
			fmt.Println(line)
		}
	}
}

func formatBranchLine(repoPath string, b types.BranchSyncStatus, noColor bool) string {
	branchPath := filepath.Join(repoPath, b.Name)

	parts := []string{branchPath}

	if b.Current {
		parts = append(parts, "[current]")
	}

	details := []string{}
	if b.NoUpstream {
		details = append(details, "no upstream")
	} else if b.Gone {
		details = append(details, "gone")
	} else {
		if b.Ahead > 0 {
			details = append(details, fmt.Sprintf("ahead %d", b.Ahead))
		}
		if b.Behind > 0 {
			details = append(details, fmt.Sprintf("behind %d", b.Behind))
		}
	}

	if len(details) > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", strings.Join(details, ", ")))
	}

	line := strings.Join(parts, " ")

	if noColor {
		return line
	}

	if b.NoUpstream {
		return ColorCyan + line + ColorReset
	}
	if b.Gone {
		return ColorMagenta + line + ColorReset
	}
	if b.Ahead > 0 && b.Behind > 0 {
		return ColorYellow + line + ColorReset
	}
	if b.Ahead > 0 {
		return ColorGreen + line + ColorReset
	}
	if b.Behind > 0 {
		return ColorRed + line + ColorReset
	}

	return line
}

func formatWorkdirLine(repoPath string, w types.WorkdirStatus, noColor bool) string {
	parts := []string{repoPath}

	details := []string{}
	if w.Modified > 0 {
		details = append(details, fmt.Sprintf("modified %d", w.Modified))
	}
	if w.Staged > 0 {
		details = append(details, fmt.Sprintf("staged %d", w.Staged))
	}
	if w.Untracked > 0 {
		details = append(details, fmt.Sprintf("untracked %d", w.Untracked))
	}

	if len(details) > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", strings.Join(details, ", ")))
	}

	line := strings.Join(parts, " ")

	if noColor {
		return line
	}

	return ColorYellow + line + ColorReset
}

func formatCleanRepoLine(repoPath string, noColor bool) string {
	line := repoPath + " (clean)"
	if noColor {
		return line
	}
	return ColorGreen + line + ColorReset
}
