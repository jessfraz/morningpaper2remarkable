package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const downloadHelp = `Download a paper and upload to remarkable cloud.`

func (cmd *downloadCommand) Name() string      { return "download" }
func (cmd *downloadCommand) Args() string      { return "[OPTIONS] URL TITLE" }
func (cmd *downloadCommand) ShortHelp() string { return downloadHelp }
func (cmd *downloadCommand) LongHelp() string  { return downloadHelp }
func (cmd *downloadCommand) Hidden() bool      { return false }

func (cmd *downloadCommand) Register(fs *flag.FlagSet) {
	fs.StringVar(&cmd.dataDir, "dataDir", "papers", "directory to download the file to")
}

type downloadCommand struct {
	dataDir string
}

func (cmd *downloadCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a url for a pdf")
	}
	if len(args) < 2 {
		return errors.New("must pass a title for the pdf")
	}

	// Create the directory in remarkable cloud.
	if err := rmAPI.Mkdir(cmd.dataDir); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"dir": cmd.dataDir,
	}).Info("successfully created directory in remarkable cloud")

	// Create the directory if it does not exist.
	if _, err := os.Stat(cmd.dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cmd.dataDir, 0755); err != nil {
			return fmt.Errorf("creating directory %s failed: %v", cmd.dataDir, err)
		}
	}

	// Download the pdf.
	logrus.WithFields(logrus.Fields{
		"link": args[0],
	}).Debug("downloading paper")

	file := filepath.Join(cmd.dataDir, args[1]+".pdf")
	if err := downloadPaper(args[0], file); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"link": args[0],
		"file": file,
	}).Info("downloaded paper to file")

	// Sync the file with remarkable cloud.
	if err := rmAPI.SyncFileAndRename(file, args[1]); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s and renamed to %s", args[0], args[1])
	return nil
}
