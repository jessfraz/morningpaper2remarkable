package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func createDirectory(dir string) error {
	// Create the directory in remarkable cloud.
	if err := rmAPI.Mkdir(dir); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"dir": dir,
	}).Info("successfully created directory in remarkable cloud")

	// Create the directory if it does not exist.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s failed: %v", dir, err)
		}
	}

	return nil
}
