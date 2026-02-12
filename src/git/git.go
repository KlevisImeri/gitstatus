package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitstatus/src/defaults"
	"gitstatus/src/logger"
	"gitstatus/src/types"
)

// Regex to parse: * main a1b2c3d [origin/main: ahead 2, behind 1] Commit message
// or: * main a1b2c3d [origin/main] Commit message
// Groups: 1=Current(* or space), 2=BranchName, 3=RemoteInfo
var branchLineRegex = regexp.MustCompile(`^([\* ])\s+(\S+)\s+\w+\s+\[([^\]]+)\]`)

// Regex to parse remote info: origin/main: ahead 2, behind 1
var aheadRegex = regexp.MustCompile(`ahead (\d+)`)
var behindRegex = regexp.MustCompile(`behind (\d+)`)

func GetRepoStatus(ctx context.Context, path string, logger *logger.Logger, cfg types.Config) (*types.RepoResult, error) {
	logger.Debug("Analyzing branches in repo: %s", path)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(defaults.DefaultGitCommandTimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "branch", "-vv")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to execute git command in %s. Error: %v. Output: %s", path, err, string(output))
		return nil, fmt.Errorf("git command failed: %w", err)
	}

	logger.Debug("Git command output for %s:\n%s", path, string(output))

	branches, err := parseGitOutput(string(output), logger)
	if err != nil {
		return nil, err
	}

	result := &types.RepoResult{
		Path:     path,
		Branches: []types.BranchSyncStatus{},
	}

	for _, b := range branches {
		if cfg.ShowAll || b.Ahead > 0 || b.Behind > 0 || b.Gone || b.NoUpstream {
			if b.Ahead > 0 || b.Behind > 0 || b.Gone || b.NoUpstream {
				result.HasUnsynced = true
			}
			result.Branches = append(result.Branches, b)
		}
	}

	logger.Debug("Repo %s: %d branches found (%d unsynced)",
		path, len(result.Branches), len(result.Branches))
	return result, nil
}

func parseGitOutput(output string, logger *logger.Logger) ([]types.BranchSyncStatus, error) {
	var branches []types.BranchSyncStatus
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "(HEAD detached at") {
			logger.Debug("Skipping detached HEAD line: %s", line)
			continue
		}

		matches := branchLineRegex.FindStringSubmatch(line)
		if matches == nil {
			branchNoUpstreamRegex := regexp.MustCompile(`^([\* ])\s+(\S+)\s+\w+\s+`)
			noUpstreamMatches := branchNoUpstreamRegex.FindStringSubmatch(line)
			if noUpstreamMatches != nil {
				isCurrent := noUpstreamMatches[1] == "*"
				name := noUpstreamMatches[2]
				b := types.BranchSyncStatus{
					Name:       name,
					Current:    isCurrent,
					NoUpstream: true,
				}
				logger.Debug("Parsed branch (no upstream): %s (Current: %v)", b.Name, b.Current)
				branches = append(branches, b)
			} else {
				logger.Debug("Skipping line (format mismatch): %s", line)
			}
			continue
		}

		isCurrent := matches[1] == "*"
		name := matches[2]
		remoteInfo := matches[3]

		b := types.BranchSyncStatus{
			Name:    name,
			Current: isCurrent,
		}

		if strings.Contains(remoteInfo, ": gone]") || strings.HasSuffix(remoteInfo, ": gone") {
			b.Gone = true
		}

		aheadMatch := aheadRegex.FindStringSubmatch(remoteInfo)
		if len(aheadMatch) > 1 {
			val, err := strconv.Atoi(aheadMatch[1])
			if err == nil {
				b.Ahead = val
			} else {
				logger.Error("Failed to parse ahead count in '%s': %v", remoteInfo, err)
			}
		}

		behindMatch := behindRegex.FindStringSubmatch(remoteInfo)
		if len(behindMatch) > 1 {
			val, err := strconv.Atoi(behindMatch[1])
			if err == nil {
				b.Behind = val
			} else {
				logger.Error("Failed to parse behind count in '%s': %v", remoteInfo, err)
			}
		}

		logger.Debug("Parsed branch: %s (Current: %v, Ahead: %d, Behind: %d, Gone: %v)",
			b.Name, b.Current, b.Ahead, b.Behind, b.Gone)

		branches = append(branches, b)
	}

	return branches, nil
}
