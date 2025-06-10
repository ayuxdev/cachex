package cmd

import (
	"fmt"

	"github.com/ayuxdev/cachex/pkg/config"
)

func buildHelpMessage(cfg *config.Config) string {
	fullTimeout := cfg.ScannerConfig.Client.DialTimeout +
		cfg.ScannerConfig.Client.HandshakeTimeout +
		cfg.ScannerConfig.Client.ResponseTimeout

	helpMessage := fmt.Sprintf(`cachex - Tool to detect cache poisoning vulnerabilities

USAGE:
  cachex [flags]

FLAGS:

INPUT:
  -u, --url                      URL to scan
  -l, --list                     Path to a file containing a list of URLs to scan

GENERAL:
  -t, --threads                  Number of threads to use (default: %d)
  -m, --scan-mode                Scan mode: single or multi (default: %s)

HTTP CLIENT:
  -timeout, --request-timeout    Request timeout in seconds (default: %.1f)
  -proxy, --proxy-url            Proxy URL to use for requests (default: %s)

PERSISTENCE CHECKER:
  -np, --no-chk-prst         	 Disable persistence checker or real time poisoning check (default: %v)
  -pr, --prst-requests           Number of requests to send for poisoning the cache (default: %d)
  -pt, --prst-threads            Number of concurrent threads to use while poisoning (default: %d)

OUTPUT:
  -j, --json                     Write JSONLines output
  -o, --output                   Path to output file (default: stdout)

PAYLOADS:
  -pcf, --payload-config-file    Path to payload config YAML file (default: %s)
`,
		cfg.ScannerConfig.Threads,
		cfg.ScannerConfig.ScanMode,
		fullTimeout,
		"None",
		!cfg.ScannerConfig.PersistenceCheckerArgs.Enabled,
		cfg.ScannerConfig.PersistenceCheckerArgs.NumRequestsToSend,
		cfg.ScannerConfig.PersistenceCheckerArgs.Threads,
		config.DefaultPayloadHeadersPath,
	)

	return helpMessage
}
