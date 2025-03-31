package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Status() {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		color.Red("❌ Error getting current directory.")
		return
	}

	// Read head index file (current commit ID)
	bheadIndex, err := os.ReadFile(filepath.Join(cwd, headIndexFilePath))
	if err != nil {
		color.Red("❌ Error reading head index file.")
		return
	}
	headCommit := strings.TrimSpace(string(bheadIndex)) // Trim spaces/newlines

	// Read commit index file
	bcommitIndex, err := os.ReadFile(filepath.Join(cwd, commitIndexFile))
	if err != nil {
		color.Red("❌ Error reading commit index file.")
		return
	}

	// Split commits by ",,"
	commitEntries := strings.Split(strings.TrimSpace(string(bcommitIndex)), ",,")

	for _, entry := range commitEntries {
		entry = strings.TrimSpace(entry) // Remove extra spaces
		parts := strings.Split(entry, "\n")

		if len(parts) < 3 {
			//color.Red("⚠️ skipping malformed commit entry: %v", entry)
			continue
		}


		commit := strings.TrimSpace(parts[0])
		time := strings.TrimSpace(parts[1])
		message := strings.TrimSpace(parts[2])

		commitId := strings.TrimSpace(strings.Split(commit, ":")[1])

		// If commit ID matches HEAD, highlight it
		if commitId == headCommit {

			color.Set(color.FgBlue, color.Bold)
			fmt.Println("\nCurrent Commit:")
			color.Unset()

			// Display commit details
			color.Set(color.FgHiYellow)
			fmt.Printf("Commit ID: %s\n", commitId)
			fmt.Printf("Time: %s\n", time)
			fmt.Printf("Message: %s\n\n", message)
			color.Unset()

			return
		}

	}
}
