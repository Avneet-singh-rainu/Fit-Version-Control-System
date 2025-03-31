package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

type PreviousCommitFileHashes struct{
	Files map[string]string
}


func Add(filename string) {

	// before staging the content make sure the user have initiated the fit
	// we can ensure it by checking if the ".fit" folder exists or not
	err := FitExists()
	if err!=nil{
		color.Red("Please initiate the fit first")
		color.Cyan("You can do so by using the command --> fit init")
		return
	}

	// Ensure the staging folder exists
	if err = os.MkdirAll(stageFolder, 0755); err != nil {
		color.Red("Error creating staging directory")
		fmt.Println(err)
		return
	}

	if filename =="."{
		stageAllFiles()
		color.Green("staged all files")
		return
	}

	// Generate a unique staged filename
	//extension := path.Ext(filename)
	//uniqueName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Intn(100000), extension)
	stagedFilePath := stageFolder + filename

	// Open the source file
	sourceFile, err := os.Open(filename)
	if err != nil {
		color.Green("Error opening source file")
		fmt.Println(err)
		return
	}
	defer sourceFile.Close()

	// Read the content of the source file
	content, err := io.ReadAll(sourceFile)
	if err != nil {
		color.Green("Error reading source file")
		fmt.Println( err)
		return
	}

	// Write content to the new staged file
	if err := os.WriteFile(stagedFilePath, content, 0644); err != nil {
		color.Green("Error writing to staged file:")
		fmt.Println(err)
		return
	}

	// Store filename mapping in the index file
	entry := fmt.Sprintf("%s %s\n", filename, filename)
	if err := EntryHashToFile(indexFile, entry); err != nil {
		color.Red("Error writing to index file")
		fmt.Println(err)
		return
	}

	color.Green("File staged successfully:", stagedFilePath)
}



func stageAllFiles() {
	// before staging the files i need to get the parent commit hash
	// using this hash i will locate the latest commit
	// after locating the latest commit i will have to read its index file to make a map [ path --> hash of the file]
	// i store it in the map
	// now for each file in the working dir i will calculate its hash
	// if the current file hash is in the latestCommitInfo map then i will return and
	// add the latest commit reference in the FIleToHash stuct
	// [ filepath -----Map to---> previous commit id ]
	// and i can locate the content of the file using the commit id in the object and after that i can can its index for the filepath that
	// i want and from there i will get the hash of the file and load its content...
	ignoreFiles := []string{}
	ignoreDirs := []string{}
	rootDir, err := os.Getwd()
    if err != nil {
		color.Set(color.FgRed)
        fmt.Println("Error getting current directory:", err)
		color.Unset()
        return
    }

	latestCommitHash,err := os.ReadFile(filepath.Join(rootDir,".fit","HEAD","index.txt"))
	if err != nil {
		// for the first staging there will be no latest commit so seeding the head.txt
        os.WriteFile(filepath.Join(rootDir,".fit","HEAD","index.txt"),[]byte{00},0666)
    }


	lastCommitIndexInfo := SFileToHash{ParentCommitId : "",Files: make(map[string]string)}
	err = CalculatePreviousHashes(&lastCommitIndexInfo)
	if err!=nil{
		color.Set(color.FgRed)
		fmt.Println("error calculating previous hashes...",err)
		color.Unset()
	}

	// initiating a fileToHash struct so that i can store the current files info
	newCommitIndexInfo := SFileToHash{ParentCommitId : "",Files: make(map[string]string)}

	// calculating the stage path for further use
    stagePath := filepath.Join(rootDir, stageFolder)
    if err := os.MkdirAll(stagePath, 0755); err != nil {
		color.Set(color.FgRed)
        fmt.Println("Error creating stage folder:", err)
		color.Unset()
        return
    }


	// read all the dirs in the cwd so that i can stage them
    entries, err := os.ReadDir(rootDir)
    if err != nil {
		color.Set(color.FgRed)
        fmt.Println("Error reading directory:", err)
		color.Unset()
        return
    }
	// iteration over the files and dirs in the cwd
    for _, entry := range entries {

		// if the entry extension is in the fitign the dont add it to the staging area...
		entryExtension := path.Ext(entry.Name())
		ignoreFiles,ignoreDirs,err = GetFitignFiles()
		if(err != nil){
			color.Set(color.FgRed)
			fmt.Println("error in GetFitignFiles",err)
			color.Unset()
		}
		if(entry.Name()==".fit" || entry.Name()==".git" || entry.Name()=="fit.exe" ){continue}

        srcPath := filepath.Join(rootDir, entry.Name())
        destPath := filepath.Join(stagePath, entry.Name())

        if entry.IsDir() {
			if Contains(ignoreDirs,"/"+entry.Name()){
				continue
			}
            //fmt.Println("Staging directory:", srcPath, "->", destPath)
            if err := CopyDirAndCompress(srcPath, destPath,&newCommitIndexInfo.Files,&lastCommitIndexInfo.Files); err != nil {
				color.Set(color.FgRed)
                fmt.Println("Error staging directory:", err)
				color.Unset()
            }
        } else {
			if Contains(ignoreFiles,entryExtension){
				continue
			}
			// i will store the dest file by his calculated hash...
            //fmt.Println("Staging file:", srcPath, "->", destPath)
            if err := CopyFileAndCompress(srcPath, destPath,&newCommitIndexInfo.Files,&lastCommitIndexInfo.Files); err != nil {
				color.Set(color.FgRed)
                fmt.Println("Error staging file:", err)
				color.Unset()
				return
             } //else {
            //     // Append file to index
            //     entry := fmt.Sprintf("%s %s\n", entry.Name(), destPath)
            //     if err := EntryHashToFile(indexFile, entry); err != nil {
            //         fmt.Println("Error writing to index file:", err)
            //     }
            // }
        }
    }



	newCommitIndexInfo.ParentCommitId=string(latestCommitHash)
	bres,err := json.Marshal(newCommitIndexInfo)

	if err!=nil{
		color.Set(color.FgRed)
		fmt.Println("error in Marshaling to json",err)
		color.Unset()
	}

    //fmt.Println(string(bres))

	file,err := os.Create(indexFile)
	if(err !=nil){
		color.Set(color.FgRed)
		fmt.Println("errror in creating stage index file",err)
		color.Unset()
	}

	file.Write(bres)
	file.Close()

	color.Set(color.FgYellow)
	fmt.Println("🔍 Ignored Directories -->", ignoreDirs)
	fmt.Println("📂 Ignored Files -->", ignoreFiles)
	color.Unset()

}



// appendToFile adds a new entry to the index file
// func appendToFile(filepath, content string) error {
// 	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	_, err = f.WriteString(content)
// 	return err
// }



func Contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}
