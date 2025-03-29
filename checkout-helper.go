package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// i will have to recursively get the file from the foreign
//
func BringAndUpdateFromReferencedCommit(destAbsPath,checkoutCommitId,foreignCommitId, fileHash string) error {
	cwd,err:= os.Getwd()
	if err!=nil{
		fmt.Println("error in reading cwd")
		return err
	}
	destFileName := filepath.Base(destAbsPath)
	destTestingPath := filepath.Join(`E:\GoProjects\go-fit\testing`, destFileName)

	// i will have to go to the foreign commit and read its index.txt
	// read its index.txt into map
	// find the value of map[destAbsPath]
	// if the value is also referencing another commit the i will recursively find the commit where the value isnt referencing
	// after that i will take the value from the map [filehash]
	// after that i can update it

	// foreign commit dir -> cwd/.fit/object/commitId
	srcCommitDir := filepath.Join(cwd,".fit","object",foreignCommitId)
	// read the index.txt
	foreignIndexFilePath := filepath.Join(srcCommitDir,"index.txt")
	bfile , err := os.ReadFile(foreignIndexFilePath)
	if err!=nil{
		color.Red("error reading foreign index file")
		return err
	}
	// the bfile is json unmarshal it into the struct
	foreignIndexInfoMap := SFileToHash{map[string]string{}}
	err = json.Unmarshal(bfile,&foreignIndexInfoMap)
	if err!=nil{
		color.Red("error unmarshalling the foreign index file",err)
	}

	// now finding the value/filehash from the map
	requiredFileHash := foreignIndexInfoMap.Files[destAbsPath]
	requiredFileHashParts := strings.Split(requiredFileHash, "---")

	// if requiredFileHashParts has more than 1 parts then i will do recursion
	// else using the required file hash i will CopyAndDecompress the file
	if len(requiredFileHashParts) > 1 {
		err = BringAndUpdateFromReferencedCommit(destAbsPath,checkoutCommitId,requiredFileHashParts[1],fileHash)
		if err!=nil{
			color.Red("error recursively referencing commitid",err)
		}
		} else {
			srcCommitFile := filepath.Join(srcCommitDir,requiredFileHash)
		err = CopyFileAndDecompress(srcCommitFile+".gz",destTestingPath)
		if err!=nil{
			color.Red("error while copying the referenced file",err)
		}
	}

	return nil
}


// accepting the dest file path from the index file and commit hash that user entered and the calculated file hash to locate the file
func BringAndUpdateFromThisCommit(destAbsPath,currCommitId,fileHash string) error {

	cwd,err:= os.Getwd()
	if err!=nil{
		fmt.Println("error in reading cwd")
		return err
	}

	srcPath := filepath.Join(cwd, ".fit", "object", currCommitId, fileHash+".gz")
	destFileName := filepath.Base(destAbsPath)
	testingDestPath := filepath.Join(`E:\GoProjects\go-fit\testing`, destFileName)


	err = CopyFileAndDecompress(srcPath,testingDestPath)
	if err!=nil{
		fmt.Println("error in creataing  file")
		return err
	}
	return nil
}
