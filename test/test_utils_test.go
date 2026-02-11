package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func SetupTestRepos(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	testEnvPath := filepath.Join(cwd, "test_env")

	if _, err := os.Stat(testEnvPath); err == nil {
		return testEnvPath
	}

	if err := os.MkdirAll(testEnvPath, 0755); err != nil {
		t.Fatalf("Failed to create test_env directory: %v", err)
	}

	runCmd := func(dir string, name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Command '%s %v' failed in %s: %v\nOutput: %s", name, args, dir, err, output)
		}
	}

	git := func(dir string, args ...string) {
		runCmd(dir, "git", args...)
	}

	remoteRepoPath := filepath.Join(testEnvPath, "remote_repo.git")
	if err := os.MkdirAll(remoteRepoPath, 0755); err != nil {
		t.Fatalf("Failed to create remote_repo.git: %v", err)
	}
	git(remoteRepoPath, "init", "--bare")

	tempSetupPath := filepath.Join(testEnvPath, "temp_setup")
	git(testEnvPath, "clone", remoteRepoPath, "temp_setup")
	git(tempSetupPath, "config", "user.email", "test@example.com")
	git(tempSetupPath, "config", "user.name", "Test User")
	runCmd(tempSetupPath, "touch", "initial_file")
	git(tempSetupPath, "add", "initial_file")
	git(tempSetupPath, "commit", "-m", "Initial commit")
	git(tempSetupPath, "push", "origin", "master")
	os.RemoveAll(tempSetupPath)

	cloneRepo := func(name string) string {
		path := filepath.Join(testEnvPath, name)
		git(testEnvPath, "clone", remoteRepoPath, name)
		git(path, "config", "user.email", "test@example.com")
		git(path, "config", "user.name", "Test User")
		return path
	}

	cloneRepo("repo_synced")
	pathAhead := cloneRepo("repo_ahead")
	runCmd(pathAhead, "touch", "ahead_file")
	git(pathAhead, "add", "ahead_file")
	git(pathAhead, "commit", "-m", "Ahead commit")

	pathBehindSetup := cloneRepo("temp_behind_setup")
	runCmd(pathBehindSetup, "touch", "behind_file")
	git(pathBehindSetup, "add", "behind_file")
	git(pathBehindSetup, "commit", "-m", "Behind commit")
	git(pathBehindSetup, "push")
	os.RemoveAll(pathBehindSetup)

	pathBehind := cloneRepo("repo_behind")
	git(pathBehind, "reset", "--hard", "HEAD~1")

	pathModified := cloneRepo("repo_modified")

	f, err := os.OpenFile(filepath.Join(pathModified, "initial_file"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("modified content"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	pathStaged := cloneRepo("repo_staged")
	runCmd(pathStaged, "touch", "staged_file")
	git(pathStaged, "add", "staged_file")

	pathUntracked := cloneRepo("repo_untracked")
	runCmd(pathUntracked, "touch", "untracked_file")

	pathNoUpstream := filepath.Join(testEnvPath, "repo_no_upstream")
	if err := os.MkdirAll(pathNoUpstream, 0755); err != nil {
		t.Fatal(err)
	}
	git(pathNoUpstream, "init")
	git(pathNoUpstream, "config", "user.email", "test@example.com")
	git(pathNoUpstream, "config", "user.name", "Test User")
	runCmd(pathNoUpstream, "touch", "local_file")
	git(pathNoUpstream, "add", "local_file")
	git(pathNoUpstream, "commit", "-m", "Local commit")

	pathGone := cloneRepo("repo_gone")
	git(pathGone, "checkout", "-b", "feature-gone")
	git(pathGone, "push", "-u", "origin", "feature-gone")

	pathGoneSetup := cloneRepo("temp_gone_setup")
	git(pathGoneSetup, "push", "origin", "--delete", "feature-gone")
	os.RemoveAll(pathGoneSetup)

	git(pathGone, "fetch", "-p")

	if err := os.MkdirAll(filepath.Join(testEnvPath, "not_a_repo"), 0755); err != nil {
		t.Fatal(err)
	}

	pathNested := filepath.Join(testEnvPath, "nested", "level1", "repo_deep")
	if err := os.MkdirAll(pathNested, 0755); err != nil {
		t.Fatal(err)
	}
	git(pathNested, "init")
	git(pathNested, "config", "user.email", "test@example.com")
	git(pathNested, "config", "user.name", "Test User")
	runCmd(pathNested, "touch", "deep_file")
	git(pathNested, "add", "deep_file")
	git(pathNested, "commit", "-m", "Deep commit")

	pathIgnored := filepath.Join(testEnvPath, "node_modules", "repo_ignored")
	if err := os.MkdirAll(pathIgnored, 0755); err != nil {
		t.Fatal(err)
	}
	git(pathIgnored, "init")

	fmt.Println("Test environment created at:", testEnvPath)
	return testEnvPath
}
