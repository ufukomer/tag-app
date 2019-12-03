package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func appendTags(root, tags, suffix string) error {
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() || !strings.HasSuffix(path, suffix) {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDWR, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		firstLine := fmt.Sprintf("// +build %v", tags)
		lines := strings.Split(string(b), "\n")
		if len(lines) > 0 {
			if strings.Contains(lines[0], "+build") {
				// Replace existing build tags.
				lines[0] = firstLine
			} else {
				lines = append([]string{firstLine}, lines...)
			}
		}

		if _, err = file.Seek(0, io.SeekStart); err != nil {
			return err
		}

		_, err = io.Copy(file, strings.NewReader(strings.Join(lines, "\n")))
		return err
	})

	return err
}
