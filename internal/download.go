package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const downloadHelp = `Download a paper and upload to remarkable cloud.`

func (cmd *DownloadCommand) Name() string      { return "download" }
func (cmd *DownloadCommand) Args() string      { return "[OPTIONS] URL TITLE" }
func (cmd *DownloadCommand) ShortHelp() string { return downloadHelp }
func (cmd *DownloadCommand) LongHelp() string  { return downloadHelp }
func (cmd *DownloadCommand) Hidden() bool      { return false }

func (cmd *DownloadCommand) Register(fs *flag.FlagSet) {
}

type DownloadCommand struct {
	cfg *AppConfig
}

func (cmd *DownloadCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a url for a pdf")
	}
	if len(args) < 2 {
		return errors.New("must pass a title for the pdf")
	}

	// We want the default dir for this command to be papers.
	if cmd.cfg.DataDir == cmd.cfg.DefaultDir {
		cmd.cfg.DataDir = "papers"
	}

	if err := CreateDataDirectory(cmd.cfg); err != nil {
		return err
	}

	// Download the pdf.
	logrus.WithFields(logrus.Fields{
		"link": args[0],
	}).Debug("downloading paper")

	file := filepath.Join(cmd.cfg.DataDir, args[1]+".pdf")
	if err := downloadPaper(args[0], file); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"link": args[0],
		"file": file,
	}).Info("downloaded paper to file")

	// Sync the file with remarkable cloud.
	if err := cmd.cfg.RemarkableAPI.SyncFileAndRename(file, args[1]); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s and renamed to %s", args[0], args[1])
	return nil
}
