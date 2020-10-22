package util

import (
	"io"
	"os"
)

func SaveFile(data io.Reader, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, data); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}