package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ayuxdev/cachex/internal/pkg/logger"
	"github.com/ayuxdev/cachex/pkg/cachex"
	"github.com/ayuxdev/cachex/pkg/config"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:  "cachex",
		Usage: "Tool to detect cache poisoning",
		Flags: BuildFlags(),
		Action: func(c *cli.Context) error {
			// Load any extra config like headers
			if err := ProcessPayloadConfigFile(c, config.Cfg); err != nil {
				logger.Errorf("failed to process payload config file: %v", err)
				return nil
			}
			ProcessRequestTimeout(requestTimeout, config.Cfg)
			if jsonOutput {
				ProcessJSONOutput(config.Cfg)
			}
			ProcessCfg(config.Cfg)
			return Run(config.Cfg)
		},
		CustomAppHelpTemplate: buildHelpMessage(config.Cfg),
	}
}

func Run(cfg *config.Config) error {
	PrintBanner()
	var urls []string
	if url != "" {
		urls = append(urls, url)
	} else if urlsFilePath != "" {
		var err error
		urls, err = fileToSlice(urlsFilePath)
		if err != nil {
			logger.Errorf("failed to read URLs from file: %v", err)
			return nil
		}
	} else {
		// Check for piped input
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// We have piped input
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					urls = append(urls, line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read from stdin: %v", err)
			}
		}
	}

	if len(urls) == 0 {
		logger.Errorf("No URLs provided")
		return nil
	}

	// Initialize and run the scanner
	scanner := cachex.Scanner{
		ScannerConfig: &cfg.ScannerConfig,
		PayloadConfig: &cfg.PayloadConfig,
		OutputFile:    outputFile,
		URLs:          urls,
	}
	scanner.Run()
	return nil
}
