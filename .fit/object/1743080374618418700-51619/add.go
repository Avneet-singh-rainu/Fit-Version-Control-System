package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)



func Add(filename string) {
	// Ensure the staging folder exists
	if err := os.MkdirAll(stageFolder, 0755); err != nil {
		fmt.Println("Error creating staging directory:", err)
		return
	}

	if filename =="."{
		stageAllFiles()
		fmt.Println("staged all files...")
		return
	}

	// Generate a unique staged filename
	//extension := path.Ext(filename)
	//uniqueName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Intn(100000), extension)
	stagedFilePath := stageFolder + filename

	// Open the source file
	sourceFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	// Read the content of the source file
	content, err := io.ReadAll(sourceFile)
	if err != nil {
		fmt.Println("Error reading source file:", err)
		return
	}

	// Write content to the new staged file
	if err := os.WriteFile(stagedFilePath, content, 0644); err != nil {
		fmt.Println("Error writing to staged file:", err)
		return
	}

	// Store filename mapping in the index file
	entry := fmt.Sprintf("%s %s\n", filename, filename)
	if err := appendToFile(indexFile, entry); err != nil {
		fmt.Println("Error writing to index file:", err)
		return
	}

	fmt.Println("File staged successfully:", stagedFilePath)
}



func stageAllFiles() {
    rootDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current directory:", err)
        return
    }

    stagePath := filepath.Join(rootDir, stageFolder)

    if err := os.MkdirAll(stagePath, 0755); err != nil {
        fmt.Println("Error creating stage folder:", err)
        return
    }

    entries, err := os.ReadDir(rootDir)
    if err != nil {
        fmt.Println("Error reading directory:", err)
        return
    }

    for _, entry := range entries {
		// if the entry extension is in the fitign the dont add it to the staging area...
		entryExtension := path.Ext(entry.Name())

		ignoreFiles,ignoreDirs,err := GetFitignFiles()
		if(err != nil){
			fmt.Println("error in GetFitignFiles",err)
		}


		if(entry.Name()==".fit" || entry.Name()==".git" || entry.Name()=="fit.exe" ){continue}

        srcPath := filepath.Join(rootDir, entry.Name())
        destPath := filepath.Join(stagePath, entry.Name())

        if entry.IsDir() {
			if Contains(ignoreDirs,"/"+entry.Name()){
				continue
			}
            fmt.Println("Staging directory:", srcPath, "->", destPath)
            if err := CopyDir(srcPath, destPath); err != nil {
                fmt.Println("Error staging directory:", err)
            }
        } else {
			if Contains(ignoreFiles,entryExtension){
				continue
			}
            fmt.Println("Staging file:", srcPath, "->", destPath)
            if err := CopyFileAndCompress(srcPath, destPath); err != nil {
                fmt.Println("Error staging file:", err)
            } else {
                // Append file to index
                entry := fmt.Sprintf("%s %s\n", entry.Name(), destPath)
                if err := appendToFile(indexFile, entry); err != nil {
                    fmt.Println("Error writing to index file:", err)
                }
            }
        }
    }
    fmt.Println("All files staged successfully.")
}



// appendToFile adds a new entry to the index file
func appendToFile(filepath, content string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}



func Contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}
