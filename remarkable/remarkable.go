package remarkable

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/juruen/rmapi/auth"
	"github.com/juruen/rmapi/cloud"
	"github.com/sirupsen/logrus"
)

// Remarkable holds the data structure for interacting with the remarkable api.
type Remarkable struct {
	api *cloud.Client
}

// New creates a new instance that is authenticated to the remarkable API.
func New() Remarkable {
	return Remarkable{api: cloud.NewClient(auth.New().Client())}
}

// SyncFileAndRename syncs a file to the remarkable cloud and then renames it.
func (r Remarkable) SyncFileAndRename(dir string, file []byte, title string) error {
	// Check if the file already exists and exit early.
	exists, err := r.checkExists(title)
	if err != nil {
		return err
	}
	if exists {
		// File exists exit early.
		logrus.WithFields(logrus.Fields{
			"title": title,
		}).Warn("file exists in cloud, skipping")
		return nil
	}

	// Make the directory or get the UUID of the directory.
	dirID, err := r.Mkdir(dir)
	if err != nil {
		return err
	}

	d := cloud.Document{
		ID:      uuid.New().String(),
		Type:    "pdf",
		Name:    title,
		Parent:  dirID,
		Version: 1,
	}

	// Zip the file.
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	f, err := w.Create("paper.pdf")
	if err != nil {
		return fmt.Errorf("creating zip file failed: %v", err)
	}
	if _, err = f.Write(file); err != nil {
		return fmt.Errorf("writing zip file failed: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("closing zip file failed: %v", err)
	}

	// Upload the file.
	if err := r.api.UploadDocument(d, bytes.NewReader(buf.Bytes())); err != nil {
		return fmt.Errorf("uploading the file %s failed: %v", title, err)
	}

	logrus.WithFields(logrus.Fields{
		"title":  title,
		"parent": dirID,
	}).Info("successfully uploaded file to remarkable cloud")

	return nil
}

// Mkdir creates a directory in the remarkable cloud.
func (r Remarkable) Mkdir(dir string) (string, error) {
	dir, err := r.api.CreateFolder(dir, "")
	if err != nil {
		return "", fmt.Errorf("creating directory %s in remarkable cloud failed: %v", dir, err)
	}
	return dir, nil
}

func (r Remarkable) checkExists(title string) (bool, error) {
	files, err := r.api.List()
	if err != nil {
		return false, fmt.Errorf("listing files failed: %v", err)
	}

	// Iterate over the files.
	for _, d := range files {
		if d.Name == title {
			// File exists.
			return true, nil
		}
	}

	return false, nil
}
