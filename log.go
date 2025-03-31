package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// CommitEntry represents a single commit record
type CommitEntry struct {
	ID      string
	Time    time.Time
	Message string
}

// Log displays the commit history from the commit index file in an appealing format
func Log() {
	// Get terminal width for formatting
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Could not determine terminal size:", err)
		width = 80 // Default width if size can't be determined
	}

	// Open and read commit index file
	content, err := readCommitFile(commitIndexFile)
	if err != nil {
		return
	}

	// Parse commit entries
	commits := parseCommits(content)

	// Display header
	headerStyle := color.New(color.FgHiWhite, color.Bold).Add(color.BgBlue)
	divider := strings.Repeat("─", width-1)

	headerStyle.Println(" COMMIT HISTORY ")
	fmt.Println(divider)

	// Display commits with formatting
	for i, commit := range commits {
		// Format commit ID
		idStyle := color.New(color.FgHiYellow, color.Bold)
		idStyle.Printf("Commit: ")
		fmt.Printf("%s\n", commit.ID)

		// Format timestamp
		timeStyle := color.New(color.FgHiCyan)
		timeStyle.Printf("Time:    ")
		fmt.Printf("%s\n", commit.Time.Format("Mon, 02 Jan 2006 15:04:05 MST"))

		// Format message
		msgStyle := color.New(color.FgHiGreen)
		msgStyle.Printf("Message: ")
		fmt.Printf("%s\n", commit.Message)

		// Add divider between entries
		if i < len(commits)-1 {
			color.New(color.FgHiBlue).Println(" ")
		}
	}

	// Final divider
	color.New(color.FgHiWhite).Println(divider)
	fmt.Printf("\nTotal commits: %d\n", len(commits))
}

// readCommitFile reads the commit index file and returns its contents
func readCommitFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Could not open commit index file:", err)
		return "", err
	}
	defer file.Close()

	content, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("Error reading the commit index file:", err)
		return "", err
	}

	return string(content), nil
}

// parseCommits parses the raw commit file content into structured CommitEntry objects
func parseCommits(content string) []CommitEntry {
	commitArray := strings.Split(content, ",,")
	var commits []CommitEntry

	for _, entry := range commitArray {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		lines := strings.Split(entry, "\n")
		if len(lines) < 3 {
			continue
		}

		// Extract commit ID, timestamp and message
		idLine := strings.TrimPrefix(lines[0], "Commit: ")
		timeLine := strings.TrimPrefix(lines[1], "Time: ")
		msgLine := strings.TrimPrefix(lines[2], "Message: ")

		// Parse time
		t, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", timeLine)

		commits = append(commits, CommitEntry{
			ID:      idLine,
			Time:    t,
			Message: msgLine,
		})
	}

	return commits
}
