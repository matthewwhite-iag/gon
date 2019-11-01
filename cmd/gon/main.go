package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/gon/config"
	"github.com/mitchellh/gon/sign"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var logLevel string
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&logLevel, "log-level", "", "Log level to output. Defaults to no logging.")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	// Build a logger
	logOut := ioutil.Discard
	if logLevel != "" {
		logOut = os.Stderr
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.LevelFromString(logLevel),
		Output: logOut,
	})

	// We expect a configuration file
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, color.RedString("❗️ Path to configuration expected.\n\n"))
		printHelp(flags)
		return 1
	}

	// Parse the configuration
	cfg, err := config.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, color.RedString("❗️ Error loading configuration:\n\n%s\n", err))
		return 1
	}

	// Perform codesigning
	err = sign.Sign(context.Background(), &sign.Options{
		Files:    cfg.Source,
		Identity: cfg.Sign.ApplicationIdentity,
		Output:   os.Stdout,
		Logger:   logger.Named("sign"),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, color.RedString("❗️ Error signing files:\n\n%s\n", err))
	}

	return 0
}

func printHelp(fs *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, strings.TrimSpace(help)+"\n\n", os.Args[0])
	fs.PrintDefaults()
}

const help = `
gon signs, notarizes, and packages binaries for macOS.

Usage: %[1]s [flags] [CONFIG]

For full help text, see the README in the GitHub repository:
http://github.com/mitchellh/gon

Flags:
`
