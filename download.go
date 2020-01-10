package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

const downloadHelp = `Download a paper and upload to remarkable cloud.`

func (cmd *downloadCommand) Name() string      { return "download" }
func (cmd *downloadCommand) Args() string      { return "[OPTIONS] URL TITLE" }
func (cmd *downloadCommand) ShortHelp() string { return downloadHelp }
func (cmd *downloadCommand) LongHelp() string  { return downloadHelp }
func (cmd *downloadCommand) Hidden() bool      { return false }

func (cmd *downloadCommand) Register(fs *flag.FlagSet) {
}

type downloadCommand struct {
}

func (cmd *downloadCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a url for a pdf")
	}
	if len(args) < 2 {
		return errors.New("must pass a title for the pdf")
	}

	// We want the default dir for this command to be papers.
	if dataDir == defaultDir {
		dataDir = "papers"
	}

	// Download the pdf.
	logrus.WithFields(logrus.Fields{
		"link": args[0],
	}).Debug("downloading paper")

	file, err := downloadPaper(args[0])
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"link": args[0],
	}).Info("downloaded paper")

	// Sync the file with remarkable cloud.
	if err := rmAPI.SyncFileAndRename(dataDir, file, args[1]); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s and named %s", args[0], args[1])
	return nil
}
