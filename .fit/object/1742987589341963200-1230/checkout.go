package main

import (
	"fmt"
	"os"
)

func Checkout(commitHash string) {


	if dirs, e := os.ReadDir(blobFolder); e != nil {
		fmt.Print("Error reading the blob folder during Checkout", e)
		return
	} else {
		for _, dir := range dirs {
			dirName := dir.Name()
			if dirName==commitHash{
				RevertChanges(dirName)
			}
		}
	}


}

func RevertChanges(dirName string) {
	rootDir, e := os.Getwd()

	if e != nil {
		fmt.Println("Error getting the current working directory:", e)
		return
	}

	fmt.Print(rootDir)







}
