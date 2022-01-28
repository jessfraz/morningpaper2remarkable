package remarkable

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/juruen/rmapi/api"
	"github.com/juruen/rmapi/filetree"
	"github.com/juruen/rmapi/log"
	"github.com/juruen/rmapi/model"
	"github.com/juruen/rmapi/util"
	"github.com/sirupsen/logrus"
)

const (
	// Number of remarkable auth retries allowed.
	rmAuthRetries = 5
)

// Remarkable holds the data structure for interacting with the remarkable api.
type Remarkable struct {
	api *api.ApiCtx
}

// New creates a new instance that is authenticated to the remarkable API.
func New() (Remarkable, error) {
	var rm *api.ApiCtx

	// Initialize the logger for the api.
	log.InitLog()

	for i := 0; i < rmAuthRetries; i++ {
		rm = api.CreateApiCtx(api.AuthHttpCtx())

		if rm.Filetree == nil && i < rmAuthRetries {
			logrus.Debug("retrying remarkable auth...")
		}
	}

	if rm.Filetree == nil {
		return Remarkable{}, errors.New("failed to build remarkable documents tree")
	}

	return Remarkable{api: rm}, nil
}

// SyncFileAndRename syncs a file to the remarkable cloud and then renames it.
func (r Remarkable) SyncFileAndRename(file, title string) error {
	if len(file) < 1 || len(title) < 1 {
		return errors.New("file and title cannot be empty")
	}

	// Get the node for the directory.
	dir, err := r.api.Filetree.NodeByPath(filepath.Dir(file), r.api.Filetree.Root())
	if err != nil || dir.IsFile() {
		return fmt.Errorf("%s is a file not a directory", filepath.Dir(file))
	}

	logrus.WithFields(logrus.Fields{
		"file":  file,
		"title": title,
	}).Debug("uploading file to remarkable cloud")

	if _, err = r.api.Filetree.NodeByPath(title, dir); err == nil {
		// File already exists.
		return nil
	}

	// Extra walk because sometimes the above does not catch it....
	var found bool
	filetree.WalkTree(dir, filetree.FileTreeVistor{
		Visit: func(node *model.Node, path []string) bool {
			entryName := filepath.Join(strings.Join(path, "/"), node.Name())

			if strings.Contains(entryName, title) {
				found = true
			}

			return found
		},
	})
	if found {
		return nil
	}

	dstDir := dir.Id()
	document, err := r.api.UploadDocument(dstDir, file)
	if err != nil {
		return fmt.Errorf("uploading file %s failed: %v", file, err)
	}

	logrus.WithFields(logrus.Fields{
		"file":  file,
		"title": title,
	}).Info("successfully uploaded file to remarkable cloud")

	r.api.Filetree.AddDocument(*document)

	// Move the file to the title.
	docName := util.DocPathToName(filepath.Base(file))
	node, err := r.api.Filetree.NodeByPath(docName, dir)
	if err != nil {
		return fmt.Errorf("getting file %s failed: %v", file, err)
	}
	if _, err := r.api.MoveEntry(node, dir, title); err != nil {
		return fmt.Errorf("moving the file to a new name failed: %v", err)
	}

	return nil
}

// Mkdir creates a directory in the remarkable cloud.
func (r Remarkable) Mkdir(dir string) error {
	if len(dir) == 0 {
		return errors.New("directory cannot be empty")
	}

	if _, err := r.api.Filetree.NodeByPath(dir, r.api.Filetree.Root()); err == nil {
		// Directory already exists
		return nil
	}

	parentDir := filepath.Dir(dir)
	newDir := filepath.Base(dir)

	if newDir == "/" || newDir == "." {
		return fmt.Errorf("%s is an invalid directory name", newDir)
	}

	parentNode, err := r.api.Filetree.NodeByPath(parentDir, r.api.Filetree.Root())
	if err != nil || parentNode.IsFile() {
		return fmt.Errorf("%s is a file not a directory", parentNode.Name())
	}

	parentID := parentNode.Id()
	if parentNode.IsRoot() {
		parentID = ""
	}

	document, err := r.api.CreateDir(parentID, newDir)
	if err != nil {
		return fmt.Errorf("creating directory %s failed: %v", newDir, err)
	}

	r.api.Filetree.AddDocument(document)

	return nil
}
