package main

import (
	"fmt"
	"os"
	"path"
)

func Checkout(commitHash string) {

	// find the commithash named dir in the "fit/object/" dir.
	//after finding the dir i need to move the content of that dir into the cwd or base dir
	//after reverting i can optionally delete the commithash folder
	if dirs, e := os.ReadDir(blobFolder); e != nil {
		fmt.Print("Error reading the blob folder during Checkout", e)
		return
	} else {
		for _, dir := range dirs {
			dirName := dir.Name()
			if dirName==commitHash{

				rootDirName, e := os.Getwd()
				if e != nil {
					fmt.Println("Error getting the current working directory:")
					return
				}
				srcDir := path.Join(rootDirName,blobFolder,dirName)
				// after finding the commithash folder i will revert the changes
				err := RevertChanges(srcDir,rootDirName)

				if err != nil{
					fmt.Println(err)
					return
				}

				// if successfully moved the content then delete the folder
				// BUT TO STAY SAFE I WILL NOT DELETE AS I MIGHT DELETE SOMETHING ELSE

				fmt.Println(dirName,"-> is to be deleted...")
				//os.RemoveAll(dirName)

			}
		}
	}
}

func RevertChanges(srcDir,destDir string) error {

	CopyDir(srcDir,destDir)


	fmt.Print(srcDir,"----------->",destDir)
	return nil
}
