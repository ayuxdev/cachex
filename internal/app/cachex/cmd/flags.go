package cmd

import (
	"fmt"
	"os"

	"github.com/ayuxdev/cachex/internal/pkg/logger"
	"github.com/ayuxdev/cachex/pkg/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	urlsFilePath   string
	outputFile     string
	url            string
	requestTimeout float64
	jsonOutput     bool
	disablePrst    bool
)

func BuildFlags() []cli.Flag {
	if err := config.LoadConfig(); err != nil {
		logger.Errorf("Failed to load config: %v", err)
		return nil
	}

	cfg := config.Cfg // pointer to global config

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Aliases:     []string{"u"},
			Usage:       "URL to scan",
			Destination: &url,
		},
		&cli.StringFlag{
			Name:        "list",
			Aliases:     []string{"l"},
			Usage:       "Path to a file containing a list of URLs",
			Destination: &urlsFilePath,
		},
		&cli.IntFlag{
			Name:        "threads",
			Aliases:     []string{"t"},
			Usage:       "Number of threads to use",
			Destination: &cfg.ScannerConfig.Threads,
			Value:       cfg.ScannerConfig.Threads,
		},
		&cli.StringFlag{
			Name:        "scan-mode",
			Aliases:     []string{"m"},
			Usage:       "Scan mode: single or multi",
			Destination: &cfg.ScannerConfig.ScanMode,
			Value:       cfg.ScannerConfig.ScanMode,
		},
		// Client config flags
		&cli.Float64Flag{
			Name:        "request-timeout",
			Aliases:     []string{"timeout"},
			Destination: &requestTimeout,
			Value: cfg.ScannerConfig.Client.ResponseTimeout +
				cfg.ScannerConfig.Client.DialTimeout +
				cfg.ScannerConfig.Client.HandshakeTimeout,
		},
		&cli.StringFlag{
			Name:        "proxy-url",
			Aliases:     []string{"proxy"},
			Destination: &cfg.ScannerConfig.Client.ProxyURL,
			Value:       cfg.ScannerConfig.Client.ProxyURL,
		},
		// Persistence checker flags
		&cli.BoolFlag{
			Name:        "no-chk-prst",
			Aliases:     []string{"np"},
			Destination: &disablePrst,
			Value:       !cfg.ScannerConfig.PersistenceCheckerArgs.Enabled,
		},
		&cli.IntFlag{
			Name:        "prst-requests",
			Aliases:     []string{"pr"},
			Destination: &cfg.ScannerConfig.PersistenceCheckerArgs.NumRequestsToSend,
			Value:       cfg.ScannerConfig.PersistenceCheckerArgs.NumRequestsToSend,
		},
		&cli.IntFlag{
			Name:        "prst-threads",
			Aliases:     []string{"pt"},
			Destination: &cfg.ScannerConfig.PersistenceCheckerArgs.Threads,
			Value:       cfg.ScannerConfig.PersistenceCheckerArgs.Threads,
		},
		&cli.BoolFlag{
			Name:        "json",
			Aliases:     []string{"j"},
			Usage:       "Write JSONLines output",
			Destination: &jsonOutput,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "Write output to file",
			Destination: &outputFile,
		},
		&cli.StringFlag{
			Name:    "payload-config-file",
			Aliases: []string{"pcf"},
			Usage:   "Path to payload config YAML file",
		},
	}
}

func ProcessPayloadConfigFile(ctx *cli.Context, cfg *config.Config) error {
	payloadFile := ctx.String("payload-config-file")
	if payloadFile == "" {
		return nil
	}
	data, err := os.ReadFile(payloadFile)
	if err != nil {
		return fmt.Errorf("failed to read payload config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg.PayloadConfig); err != nil {
		return fmt.Errorf("failed to parse payload config: %w", err)
	}
	return nil
}

func ProcessRequestTimeout(requestTimeout float64, cfg *config.Config) {
	if requestTimeout > 0 {
		meanTimeout := requestTimeout / 3
		cfg.ScannerConfig.Client.DialTimeout = meanTimeout
		cfg.ScannerConfig.Client.HandshakeTimeout = meanTimeout
		cfg.ScannerConfig.Client.ResponseTimeout = meanTimeout
	}
}

func ProcessJSONOutput(cfg *config.Config) {
	cfg.ScannerConfig.LoggerConfig.LogMode = "json"
}

func ProcessCfg(cfg *config.Config) {
	if disablePrst {
		cfg.ScannerConfig.PersistenceCheckerArgs.Enabled = false
	}
}
