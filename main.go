package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"gitstatus/src/logger"
	"gitstatus/src/output"
	"gitstatus/src/types"
	"gitstatus/src/walker"
)

func parseLogTypes(logStr string) []string {
	if logStr == "" {
		return nil
	}
	return strings.Split(logStr, ",")
}

func main() {
	depth := flag.Int("depth", 0, "Maximum directory depth (0 = unlimited)")
	logLevels := flag.String("log", "", "Log levels (comma-separated: DEBUG, INFO, WARNING, ERROR)")
	showAll := flag.Bool("all", false, "Show all repositories including clean ones")
	noColor := flag.Bool("no-color", false, "Disable colored output")
	logFile := flag.String("logfile", "", "Log file path (optional)")
	flag.Parse()

	rootPath := "."
	if len(flag.Args()) > 0 {
		rootPath = flag.Args()[0]
	}

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	cfg := types.Config{
		RootPath: absRoot,
		MaxDepth: *depth,
		LogTypes: parseLogTypes(*logLevels),
		ShowAll:  *showAll,
		NoColor:  *noColor,
		LogFile:  *logFile,
	}

	logger, err := logger.NewLogger(cfg.LogTypes, cfg.LogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Received interrupt signal, stopping...")
		cancel()
	}()

	logger.Info("Starting git status scan in: %s", cfg.RootPath)

	var results []types.RepoResult

	err = walker.Walk(ctx, cfg, logger, func(res types.RepoResult) {
		results = append(results, res)
	})

	if err != nil && err != context.Canceled {
		logger.Error("Walk failed: %v", err)
	}

	logger.Info("Scan complete. Found %d repositories.", len(results))

	output.PrintResults(results, cfg, logger)
}
