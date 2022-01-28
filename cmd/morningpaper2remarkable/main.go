package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/pseudo-su/morningpaper2remarkable/internal"
	"github.com/pseudo-su/morningpaper2remarkable/internal/remarkable"
	"github.com/pseudo-su/morningpaper2remarkable/internal/version"
	"github.com/sirupsen/logrus"
)
func main() {
	cfg := internal.NewAppConfig()

	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "morningpaper2remarkable"
	p.Description = "A bot to sync the morning paper to a remarkable tablet"
	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Build the list of available commands.
	p.Commands = []cli.Command{
		&internal.DownloadCommand{},
	}

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("morningpaper2remarkable", flag.ExitOnError)
	p.FlagSet.BoolVar(&cfg.Debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")

	p.FlagSet.StringVar(&cfg.DataDir, "dir", cfg.DefaultDir, "directory to store the downloaded papers in")

	p.FlagSet.DurationVar(&cfg.Interval, "interval", 18*time.Hour, "update interval (ex. 5ms, 10s, 1m, 3h)")
	p.FlagSet.BoolVar(&cfg.Once, "once", false, "run once and exit, do not run as a daemon")

	p.FlagSet.IntVar(&cfg.MaxPages, "pages", 3, "number of pages of papers to download")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if cfg.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Authenticate with remarkable cloud.
		logrus.Info("authenticating with remarkable cloud")
		var err error
		cfg.RemarkableAPI, err = remarkable.New()
		if err != nil {
			return err
		}

		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		if err := internal.CreateDataDirectory(cfg); err != nil {
			return err
		}

		ticker := time.NewTicker(cfg.Interval)

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
		if err := internal.GetFiles(cfg); err != nil {
			return err
		}

		if !cfg.Once {
			logrus.Infof("Starting bot to update every %s", cfg.Interval)
			for range ticker.C {
				// Parse the RSS feed.
				if err := internal.GetFiles(cfg); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Run our program.
	p.Run()
}
