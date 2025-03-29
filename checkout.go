package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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
				cwd, e := os.Getwd()
				fmt.Println("ced,,,",cwd)
				if e != nil {
					fmt.Println("Error getting the current working directory:")
					return
				}
				srcDir := path.Join(cwd,blobFolder,dirName)
				// after finding the commithash folder i will revert the changes
				err := RevertChanges(srcDir,cwd,commitHash)

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

// here i have commit object where user wants to revert , current working directory
// first i will read the index.txt file in the "srcDir"
// after reading i will create the maping "MAP" of the "fileToHash"
// for each entry in the "MAP" that is path i will load the content
// if content doenot have prefic commit then its easy
// if the entry has a prefix "COMMIT-" then i will get the key and navigate to that commit
// and from there i can get the file having the same hash that i have store in the format ["COMMIT-\COMMITID-\FILEHASH"]


func RevertChanges(prevCommitDir, cwd,currCommitId string) error {
    //fmt.Println("Starting RevertChanges from", prevCommitDir, "to", cwd)

	// i will read the index.txt from the prevCommitDir and converet it into MAP of FileToHash
	prevCommitIndexMap := SFileToHash{Files:make(map[string]string)}

   // Read index file of the previous commit
   indexFile, err := os.ReadFile(prevCommitDir+"/index.txt")
   if err != nil {
	   return fmt.Errorf("error reading source directory: %v", err)
   }

   // converting the index File into struct
   json.Unmarshal(indexFile,&prevCommitIndexMap)
   // now i have the map of all the files and their location
   // now i will iterate thru all the entries in the map and for each entry i have location
   // according to that location i will update the content
   // map----> key(path) : value(commit-/commitid-/filehash)



	// Get map keys
   for k,v := range prevCommitIndexMap.Files{
		values := strings.Split(v, "---")
		if len(values)>1{
			BringAndUpdateFromReferencedCommit(k,currCommitId,values[1],values[2])
		} else {
			err = BringAndUpdateFromThisCommit(k,currCommitId,values[0])
			if err!=nil{
				fmt.Printf("error updating file from %v commit ---> %v\n",currCommitId,err)
			}
		}
   }

   // i will iterate the map and for each key which is path i will create the directory here and update its content


	// Read all files from the prevCommitDir
    // files, err := os.ReadDir(prevCommitDir)
    // if err != nil {
    //     return fmt.Errorf("error reading source directory: %v", err)
    // }

    // for _, file := range files {
    //     srcFilePath := filepath.Join(prevCommitDir, file.Name())
    //     destFilePath := filepath.Join(cwd, file.Name())

	// 	if file.IsDir(){
	// 		CopyDirAndDecompress(srcFilePath,destFilePath)
	// 		continue
	// 	}

    //     // Ensure we are reading compressed files correctly
    //     if !strings.HasSuffix(file.Name(), ".gz"){
    //         fmt.Println("Skipping non-gzipped file:", srcFilePath)
    //         continue
    //     }

    //     err := CopyFileAndDecompress(srcFilePath, destFilePath)
    //     if err != nil {
    //         fmt.Println("Error decompressing file:", err)
    //         return err
    //     }
    // }

    fmt.Println("✅ Successfully reverted changes")
    return nil
}
