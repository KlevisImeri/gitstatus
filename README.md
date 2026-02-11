# Git Repo Scanner

A powerful Go CLI tool that recursively scans directories to discover Git
repositories and provides comprehensive branch synchronization status. It
identifies repositories with unsynced branches and uncommitted changes,
displaying results in a clean, file-list format perfect for quick assessments
of many Git projects in a folder.

> "I simply needed a way to check all my projects at the end of the day to see
> if anything was uncommitted, unpushed, or needed to be pulled and merged."

## Features

- **File-List Output**: Displays each unsynced branch as a simple path like `/path/to/repo/branch-name (ahead 2, behind 1)`
- **Branch Status Detection**: Identifies branches that are:
  - **Ahead**: Have local commits not pushed to remote
  - **Behind**: Missing commits from remote
  - **Gone**: Remote branch has been deleted
- **Efficient Traversal**: Skips non-git directories, `node_modules`, `vendor`, and other common directories
- **Color-Coded Output**: 
  - Green: Ahead only
  - Red: Behind only
  - Yellow: Both ahead and behind
  - Magenta: Gone (remote deleted)
- **Zero Dependencies**: Uses only Go standard library

## Installation

```bash
go build
```

Or install to your PATH:

```bash
go install
```

## Usage

### Basic Usage

Scan current directory:
```bash
gitstatus
```

Scan specific directory:
```bash
gitstatus /path/to/projects
```

### Examples

**Show only unsynced branches in current directory:**
```bash
gitstatus
```

**Show only unsynced branches in specific directory:**
```bash
gitstatus ~/projects
```

**Show all branches including synced ones:**
```bash
gitstatus ~/projects -all
```

**Enable verbose logging:**
```bash
gitstatus ~/projects -v
```

**Limit depth to 2 levels:**
```bash
gitstatus ~/projects -depth 2
```

**Save logs to file:**
```bash
gitstatus ~/projects -v -log scan.log
```

## Example Output

```
/home/user/projects/backend-api/main [current] (ahead 3, behind 1)
/home/user/projects/frontend-app/develop (behind 5)
/home/user/projects/frontend-app/feature/auth [current] (ahead 2)
/home/user/projects/shared-lib/master (gone)
```

## How It Works

1. **Directory Traversal**: Walks the directory tree sequentially, checking for `.git` directories
2. **Git Analysis**: For each repository found, runs `git branch -vv` to get detailed branch information
3. **Status Parsing**: Parses the git output to extract:
   - Current branch (marked with `[current]`)
   - Ahead/behind counts from tracking information
   - "Gone" status for deleted remote branches
4. **Filtering**: Only shows branches that are ahead, behind, or gone (unless `-all` is used)
5. **Output**: Prints each branch as a simple path with status information


## Requirements

- Go 1.18 or later
- Git installed and accessible in PATH

## License

MIT License
