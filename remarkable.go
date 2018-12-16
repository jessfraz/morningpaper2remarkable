package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/juruen/rmapi/util"
	"github.com/sirupsen/logrus"
)

func syncFileWithRemarkableCloud(file, title string) error {
	if len(file) < 1 || len(title) < 1 {
		return errors.New("file and title cannot be empty")
	}

	// Get the node for the directory.
	dir, err := rmCtx.Filetree.NodeByPath(filepath.Dir(file), rmCtx.Filetree.Root())
	if err != nil || dir.IsFile() {
		return fmt.Errorf("%s is a file not a directory", filepath.Dir(file))
	}

	logrus.WithFields(logrus.Fields{
		"file":  file,
		"title": title,
	}).Debug("uploading file to remarkable cloud")

	if _, err := rmCtx.Filetree.NodeByPath(title, dir); err == nil {
		// File already exists.
		return nil
	}

	dstDir := dir.Id()
	document, err := rmCtx.UploadDocument(dstDir, file)
	if err != nil {
		return fmt.Errorf("uploading file %s failed: %v", file, err)
	}

	logrus.WithFields(logrus.Fields{
		"file":  file,
		"title": title,
	}).Info("successfully uploaded file to remarkable cloud")

	rmCtx.Filetree.AddDocument(*document)

	// Move the file to the title.
	docName := util.DocPathToName(filepath.Base(file))
	node, err := rmCtx.Filetree.NodeByPath(docName, dir)
	if err != nil {
		return fmt.Errorf("getting file %s failed: %v", file, err)
	}
	if _, err := rmCtx.MoveEntry(node, dir, title); err != nil {
		return fmt.Errorf("moving the file to a new name failed: %v", err)
	}

	return nil
}

func remarkableMkdir(dir string) error {
	if len(dir) == 0 {
		return errors.New("directory cannot be empty")
	}

	if _, err := rmCtx.Filetree.NodeByPath(dir, rmCtx.Filetree.Root()); err == nil {
		// Directory already exists
		return nil
	}

	parentDir := filepath.Dir(dir)
	newDir := filepath.Base(dir)

	if newDir == "/" || newDir == "." {
		return fmt.Errorf("%s is an invalid directory name", newDir)
	}

	parentNode, err := rmCtx.Filetree.NodeByPath(parentDir, rmCtx.Filetree.Root())
	if err != nil || parentNode.IsFile() {
		return fmt.Errorf("%s is a file not a directory", parentNode.Name())
	}

	parentID := parentNode.Id()
	if parentNode.IsRoot() {
		parentID = ""
	}

	document, err := rmCtx.CreateDir(parentID, newDir)
	if err != nil {
		return fmt.Errorf("creating directory %s failed: %v", newDir, err)
	}

	rmCtx.Filetree.AddDocument(document)

	return nil
}
