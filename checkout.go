package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Checkout(commitHash string) {
	// i need to check if fit is initiated or not?

	err := FitExists()
	if(err!=nil){
		color.Red("---Please initiate the fit first---")
		return
	}

	// before checking out i need to ensure that there are some previous commits
	// if there are no commits then i will just return
	// no previous versions found

	err = CommitsExists()

	if err!=nil{
		color.Red(err.Error())
	}




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
				fmt.Println("successfully Checkout done...🚀🚀")
				//os.RemoveAll(dirName)

			}
		}
	}
}

func RevertChanges(srcDir, destDir string) error {
    fmt.Println("Starting RevertChanges from", srcDir, "to", destDir)

    // Read all files from the srcDir
    files, err := os.ReadDir(srcDir)
    if err != nil {
        return fmt.Errorf("error reading source directory: %v", err)
    }

    for _, file := range files {
        srcFilePath := filepath.Join(srcDir, file.Name())
        destFilePath := filepath.Join(destDir, file.Name())

		if file.IsDir(){
			CopyDirAndDecompress(srcFilePath,destFilePath)
			continue
		}

        // Ensure we are reading compressed files correctly
        if !strings.HasSuffix(file.Name(), ".gz"){
            fmt.Println("Skipping non-gzipped file:", srcFilePath)
            continue
        }

        err := CopyFileAndDecompress(srcFilePath, destFilePath)
        if err != nil {
            fmt.Println("Error decompressing file:", err)
            return err
        }
    }

    fmt.Println("✅ Successfully reverted changes")
    return nil
}
