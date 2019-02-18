package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/morningpaper2remarkable/version"
	"github.com/juruen/rmapi/api"
	"github.com/juruen/rmapi/log"
	"github.com/sirupsen/logrus"
)

const (
	morningPaperRSSFeedURL = "https://blog.acolyer.org/feed/?paged=%d"

	// Number of remarkable auth retries allowed.
	rmAuthRetries = 5

	// Number of pages to iterate over.
	maxPages = 3
)

var (
	debug   bool
	dataDir string

	interval time.Duration
	once     bool

	rmCtx *api.ApiCtx
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "morningpaper2remarkable"
	p.Description = "A bot to sync the morning paper to a remarkable tablet"
	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("morningpaper2remarkable", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	p.FlagSet.StringVar(&dataDir, "dir", "morningpaper", "directory to store the downloaded papers in")

	p.FlagSet.DurationVar(&interval, "interval", 18*time.Hour, "update interval (ex. 5ms, 10s, 1m, 3h)")
	p.FlagSet.BoolVar(&once, "once", false, "run once and exit, do not run as a daemon")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Authenticate with remarkable cloud.
		logrus.Info("authenticating with remarkable cloud")
		log.InitLog()
		for i := 0; i < rmAuthRetries; i++ {
			rmCtx = api.CreateApiCtx(api.AuthHttpCtx())

			if rmCtx.Filetree == nil && i < rmAuthRetries {
				logrus.Info("retrying remarkable auth...")
			}
		}

		if rmCtx.Filetree == nil {
			return errors.New("failed to build remarkable documents tree")
		}

		// Create the directory in remarkable cloud.
		if err := remarkableMkdir(dataDir); err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"dir": dataDir,
		}).Info("successfully created directory in remarkable cloud")

		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		// Create the directory if it does not exist.
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				return fmt.Errorf("creating directory %s failed: %v", dataDir, err)
			}
		}

		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				ticker.Stop()
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		// If the user passed the once flag, just do the run once and exit.
		if once {
			return getFiles()
		}

		logrus.Infof("Starting bot to update every %s", interval)
		for range ticker.C {
			// Parse the RSS feed.
			if err := getFiles(); err != nil {
				return err
			}
		}

		return nil
	}

	// Run our program.
	p.Run()
}
