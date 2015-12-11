package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// extractSecret extracts a secret based on given variables
func saveTargetFile(targetPath, content string) error {
	if targetPath == "" {
		return maskAny(fmt.Errorf("target not set"))
	}

	folder := filepath.Dir(targetPath)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return maskAny(err)
	}

	if err := ioutil.WriteFile(targetPath, []byte(content), 0400); err != nil {
		return maskAny(err)
	}

	return nil
}
