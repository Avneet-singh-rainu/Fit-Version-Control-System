package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)



func Commit(commitMessage string) {

	//create unique commit id
	//create a new folder in the "fit/object/unique-commit-id"
	//for each entry in the stash move the file into the "fit/object/unique-commit-id"
	//delete the stash dir
	//append the commit info into the "fit/commit/index.txt"
	//store the latest commit info in the "fit/HEAD/index.txt"

	// create a unique commit ID
	uniqueName := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(100000))
	commitPath := filepath.Join(blobFolder, uniqueName)

	// create a directory for the commit inside ".fit/object/"
	if err := os.MkdirAll(commitPath, 0750); err != nil {
		fmt.Println("Error creating commit directory:", err)
		return
	}

	// Move files from stage to commit directory
	if err := MoveDir(stageFolder, commitPath); err != nil {
		fmt.Println("Error moving staged files:", err)
		return
	}

	// Ensure commit folder exists
	if err := os.MkdirAll(commitFolder, 0750); err != nil {
		fmt.Println("Error ensuring commit folder exists:", err)
		return
	}

	// Open commit index file
	file, err := os.OpenFile(commitIndexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening commit index file:", err)
		return
	}
	defer file.Close()

	// Write commit details
	commitInfo := fmt.Sprintf("Commit: %s\nTime: %s\nMessage: %s\n %s \n", uniqueName, time.Now().Format(time.RFC1123), commitMessage,commitSeperator)
	_, err = file.WriteString(commitInfo)
	if err != nil {
		fmt.Println("Error writing to commit index file:", err)
	}

	fmt.Println("Commit successful:", uniqueName)
}


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
