package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func MoveDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// copy directory
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(src)
}

func CopyFile(srcFile, destFile string) error {
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func CopyDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, 0750); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {

			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {

			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}


func GetFitignFiles()(ignoreFiles []string,ignoreDirs []string,err error){

	file, err := os.Open(".fitign")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // Ensure file is closed

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	rawEntries := strings.Split(string(content), "\n")

	for _, entry := range rawEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if strings.HasPrefix(entry, "/") {
			ignoreDirs = append(ignoreDirs, entry)
		} else {
			ignoreFiles = append(ignoreFiles, entry)
		}
	}


	fmt.Print(ignoreDirs)
	fmt.Println("")
	fmt.Print(ignoreFiles)


	return ignoreFiles, ignoreDirs, nil


}
