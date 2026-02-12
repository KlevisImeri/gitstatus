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
		if res.HasUnsynced {
			hasIssues = true
			break
		}
	}

	if !hasIssues && !cfg.ShowAll {
		fmt.Println("No git repositories with unsynced status found.")
		return
	}

	for _, res := range results {
		if res.Error != nil {
			logger.Error("Error in repository %s: %v", res.Path, res.Error)
			continue
		}

		// Skip if no unsynced branches and not showing all
		if !res.HasUnsynced && !cfg.ShowAll {
			continue
		}

		// Show branches if:
		// - ShowAll is true (show all branches), OR
		// - HasUnsynced is true (show only unsynced branches)
		if cfg.ShowAll || res.HasUnsynced {
			for _, b := range res.Branches {
				line := formatBranchLine(res.Path, b, cfg.NoColor)
				fmt.Println(line)
			}
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
