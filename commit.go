package main

import (
	"fmt"
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
