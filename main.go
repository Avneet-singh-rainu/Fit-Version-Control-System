package main

import (
	"os"

	"github.com/fatih/color"
)

// Paths
const blobFolder = ".fit/object/"
const commitFolder = ".fit/commit/"
const commitIndexFile = commitFolder + "index.txt"
const stageFolder = ".fit/stage/"
const indexFile = stageFolder + "index.txt"
const commitSeperator = ",,"


func main() {
	args := os.Args
	var command string
	var filename string
	var commitMessage string
	var commitHash string

	if len(args) > 1 {
		command = args[1]
	} else {
		color.Yellow("Version ----> %s", Version)
		color.Cyan("Hi i am 'fit' your version control system.")
		color.Cyan("Avneet Singh gave me birth.")
		color.Cyan("I am built in Golang which makes me pretty.")
		return
	}


	switch command {

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
