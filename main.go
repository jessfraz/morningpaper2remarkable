package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/morningpaper2remarkable/remarkable"
	"github.com/jessfraz/morningpaper2remarkable/version"
	"github.com/sirupsen/logrus"
)

const (
	morningPaperRSSFeedURL = "https://blog.acolyer.org/feed/?paged=%d"

	defaultDir = "morningpaper"
)

var (
	debug   bool
	dataDir string

	interval time.Duration
	once     bool
	maxPages int

	rmAPI remarkable.Remarkable
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "morningpaper2remarkable"
	p.Description = "A bot to sync the morning paper to a remarkable tablet"
	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Build the list of available commands.
	p.Commands = []cli.Command{
		&downloadCommand{},
	}

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("morningpaper2remarkable", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	p.FlagSet.StringVar(&dataDir, "dir", defaultDir, "directory to store the downloaded papers in")

	p.FlagSet.DurationVar(&interval, "interval", 18*time.Hour, "update interval (ex. 5ms, 10s, 1m, 3h)")
	p.FlagSet.BoolVar(&once, "once", false, "run once and exit, do not run as a daemon")

	p.FlagSet.IntVar(&maxPages, "pages", 3, "number of pages of papers to download")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Authenticate with remarkable cloud.
		rmAPI = remarkable.New()

		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
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
		if err := getFiles(); err != nil {
			return err
		}

		if !once {
			logrus.Infof("Starting bot to update every %s", interval)
			for range ticker.C {
				// Parse the RSS feed.
				if err := getFiles(); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Run our program.
	p.Run()
}
