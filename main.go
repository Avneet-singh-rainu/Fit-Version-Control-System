package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// Paths
const blobFolder = ".fit/object/"
const commitFolder = ".fit/commit/"
const commitIndexFile = commitFolder + "index.txt"
const stageFolder = ".fit/stage/"
const indexFile = stageFolder + "index.txt"
const commitSeperator = ",,"
const stageIndexFile = ".fit/stage/index.txt"
//var previousCommitId string = ""
var headIndexFilePath  string =".fit/HEAD/index.txt"
var terminalWidth int = 80;

func main() {
	var err error
	terminalWidth, _, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Could not determine terminal size:", err)
	}

	args := os.Args
	var command string
	var filename string
	var commitMessage string
	var commitHash string

	if len(args) > 1 {
		command = args[1]
	} else {
		// Display welcome message with version information
		headerStyle := color.New(color.FgHiBlack, color.Bold).Add(color.BgBlue)
		headerStyle.Printf("  Version ----> %s  \n", Version)

		color.New(color.FgCyan).Println("Hi, I am 'fit' - your version control system.")
		color.New(color.FgCyan).Println("Created by Avneet Singh.")
		color.New(color.FgCyan).Println("Built with Go, which makes me efficient and elegant.")
		return
	}


	switch command {

	case "status":{
		Status()
		break
	}

	case "init":{
			Init()
			break
	}

	case "help":{
		Help()
		break
	}

	case "add":{
		if len(args) > 2 {
			filename = args[2]
			Add(filename)
		} else {
			color.Red("Filename argument is missing")
		}
		break
	}

	case "commit":{
		if len(args) > 3 {
			commitMessage = args[3]
			Commit(commitMessage)
		} else {
			color.Red("Filename argument is missing")
		}
		break
	}

	case "log":{
		Log()
	}


	case "cto":{
		if len(args) > 1 {
			commitHash = args[2]
			Checkout(commitHash)
		} else {
			color.Red("Filename argument is missing")
		}
		break
	}


	default :{
		color.Red("Please provide a valid command")
		color.Green("fit [help]")
		break
	}
	}
}
