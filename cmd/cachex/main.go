package main

import (
	"os"

	"github.com/ayuxdev/cachex/internal/app/cachex/cmd"
	"github.com/ayuxdev/cachex/internal/pkg/logger"
)

func main() {
	if err := cmd.App().Run(os.Args); err != nil {
		logger.Errorf(err.Error())
	}
}
