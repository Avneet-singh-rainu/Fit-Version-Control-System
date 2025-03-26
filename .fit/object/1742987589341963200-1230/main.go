package main

import (
	"fmt"
	"os"
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
		fmt.Println("No arguments provided")
	}


	switch command {
	case "init":{
			Init()
			break
	}

	case "add":{
		if len(args) > 2 {
			filename = args[2]
			Add(filename)
		} else {
			fmt.Println("Filename argument is missing")
		}
		break
	}

	case "commit":{
		if len(args) > 3 {
			commitMessage = args[3]
			Commit(commitMessage)
		} else {
			fmt.Println("Filename argument is missing")
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
			fmt.Println("Filename argument is missing")
		}
		break
	}




	default :{
		fmt.Print("default case")
		break
	}


	}
}
