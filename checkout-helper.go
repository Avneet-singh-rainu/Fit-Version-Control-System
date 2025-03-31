package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
)

// i will have to recursively get the file from the foreign
func BringAndUpdateFromReferencedCommit(destAbsPath,checkoutCommitId,foreignCommitId, fileHash string) error {
	cwd,err:= os.Getwd()
	if err!=nil{
		color.Set(color.FgRed, color.Bold)
		fmt.Println("error in reading cwd")
		color.Unset()
		return err
	}
	//destFileName := filepath.Base(destAbsPath)
	//destTestingPath := filepath.Join(`E:\GoProjects\go-fit`, destFileName)

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
	foreignIndexInfoMap := SFileToHash{ParentCommitId:"", Files : map[string]string{}}
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
			err = CopyFileAndDecompress(srcCommitFile+".gz",destAbsPath)
			if err!=nil{
				color.Red("error while copying the referenced file",err)
			}
			//color.Yellow("Successfully reverted changes: ",srcCommitFile)
		}

	return nil
}


// accepting the dest file path from the index file and commit hash that user entered and the calculated file hash to locate the file
func BringAndUpdateFromThisCommit(destAbsPath,currCommitId,fileHash string) error {

	cwd,err:= os.Getwd()
	if err!=nil{
		color.Red("error in reading cwd")
		return err
	}

	srcPath := filepath.Join(cwd, ".fit", "object", currCommitId, fileHash+".gz")
	//destFileName := filepath.Base(destAbsPath)
	//testingDestPath := filepath.Join(`E:\GoProjects\go-fit`, destFileName)


	err = CopyFileAndDecompress(srcPath,destAbsPath)
	if err!=nil{
		color.Red("error in creating file")
		return err
	}
	return nil
}



func ChangeHead(userGivenCommit string) error {
	cwd,err:=os.Getwd()
	if err!=nil{
		color.Red("error getting cwd")
		return err
	}

	headFilePath := filepath.Join(cwd,".fit","HEAD","index.txt")
	os.WriteFile(headFilePath, []byte(userGivenCommit), 0666)
	return nil
}

func RemoveOrphanFilesAndDirs(targetCommitId string) error {
	cwd, err := os.Getwd()
	if err != nil {
		color.Red("error fetching working directory...")
		return err
	}

	// Fetch fitign files and dirs -----------------------------------------------------------
	fitignFiles, fitignDirs, err := GetFitignFiles()
	if err != nil {
		color.Red("error fetching fitign files...")
		return err
	}

	// Get target commit index paths -----------------------------------------------------------
	targetCommitIndexFilePath := filepath.Join(cwd, ".fit", "object", targetCommitId, "index.txt")
	bfile, err := os.ReadFile(targetCommitIndexFilePath)
	if err != nil {
		color.Red("error fetching target commit files...")
		return err
	}

	targetCommitIndexFile := SFileToHash{}
	json.Unmarshal(bfile, &targetCommitIndexFile)
	targetCommitIndexPaths := slices.Sorted(maps.Keys(targetCommitIndexFile.Files))

	// fmt.Println("fitignFiles...\n", fitignFiles)
	// fmt.Println("fitignDirs...\n", fitignDirs)
	// fmt.Println("targetCommitIndexPaths...\n", targetCommitIndexPaths)

	// Scan directory and collect files/dirs for removal -----------------------------------------------------------
	dirs, err := os.ReadDir(cwd)
	if err != nil {
		color.Red("error iterating cwd files...")
		return err
	}

	filesToRemove := []string{}
	dirsToRemove := []string{}

	// iterating the cwd files and dirs -----------------------------------------------------------

	for _, entry := range dirs {
		entryPath := filepath.Join(cwd, entry.Name())

		// skip dir if they are in `fitignDirs`
		if entry.IsDir() {

			if entry.Name()==".fit"{
				continue
			}

			if contains(fitignDirs, "/"+entry.Name()) {
				continue
			}

			// recursively check if the directory contains tracked files
			if !isDirectoryTracked(entryPath, targetCommitIndexPaths) {
				dirsToRemove = append(dirsToRemove, entryPath)
			}
		} else {
			if entry.Name()==".fitign"{
				continue
			}
			//fmt.Println("file entry path,,,,,,,,",entryPath)
			// make sure that i compare ext with fitign files and path with indexpath
			if contains(fitignFiles, path.Ext(entry.Name())) || contains(targetCommitIndexPaths, entryPath) {
				continue
			}else{
				filesToRemove = append(filesToRemove, entryPath)
			}
		}
	}

	// remove files first
	for _, file := range filesToRemove {
		color.Set(color.FgGreen)
		fmt.Println("Removing file:", file)
		os.Remove(file)
		color.Unset()
	}

	// remove directories only if empty
	for _, dir := range dirsToRemove {
		color.Set(color.FgGreen)
		fmt.Println("Removing directory:", dir)
		os.RemoveAll(dir)
		color.Unset()
	}

	return nil
}

// checks if a directory has any tracked files
func isDirectoryTracked(dirPath string, trackedFiles []string) bool {
	for _, filePath := range trackedFiles {
		if strings.HasPrefix(filePath, dirPath) {
			return true // the directory contains at least one tracked file
		}
	}
	return false
}



// helper function to check if a slice contains a value
func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
