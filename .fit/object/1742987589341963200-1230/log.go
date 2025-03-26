package main

import (
	"fmt"
	"os"
	"strings"
)

func Log() {
	file, err := os.Open(commitIndexFile)
	if err != nil {
		fmt.Println("Could not open commit index file:", err)
		return
	}
	defer file.Close()
	content,e := os.ReadFile(file.Name())
	if e !=nil{
		fmt.Print("Error reading the commit index file...",e)
		return
	}
	contentStr := string(content)
	commitArray := strings.Split(contentStr, ",,")

	for _ ,commit := range commitArray{
		fmt.Print(commit)
	}

}
