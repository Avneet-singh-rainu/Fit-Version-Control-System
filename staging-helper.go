package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type SFileToHash struct{
	ParentCommitId string
	Files map[string]string
}




// CopyDirAndCompress recursively copies and compresses a directory.
func CopyDirAndCompress(src, dest string , fileToHash , latestCommitIndex *map[string]string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, 0750); err != nil {
		return err
	}

	for _, entry := range entries {


		// i will store the dest file by his calculated hash...
		hashDestFilePath,err := CalculateHash(entry.Name())
		if err!=nil{
			fmt.Println("error in calculating hash",err)
		}

		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively compress subdirectories
			if err := CopyDirAndCompress(srcPath, destPath,fileToHash,latestCommitIndex); err != nil {
				return err
			}
		} else {
			// Compress individual files

            fmt.Println("Staging file:", srcPath, "->", hashDestFilePath)
			if err := CopyFileAndCompress(srcPath, destPath,fileToHash,latestCommitIndex); err != nil {
				return err
			} else {
                // Append file to index
                entry := fmt.Sprintf("%s %s\n", entry.Name(), hashDestFilePath)
                if err := EntryHashToFile(indexFile, entry); err != nil {
                    fmt.Println("Error writing to index file:", err)
                }
            }
		}
	}

	return nil
}



// CopyFileAndCompress compresses a file and saves it.
func CopyFileAndCompress(srcFilePath, destFilePath string, newCommitIndexInfo,latestCommitIndex *map[string]string) error {
	// calculating the hash of the current file so that i can compare it with the entries in the latestComitIndexFile
	srcFileHash,err := CalculateHash(srcFilePath)
	if err != nil {
		return err
	}

	cwd,err := os.Getwd()
	if err!=nil{
		color.Red("error getting cwd...")
	}

	blastCommitId,err := os.ReadFile(filepath.Join(cwd,".fit","HEAD","index.txt"))
	if err!=nil{
		color.Red("error reading head index txt...")
	}

	// var hashList []string
	// for _, v := range *latestCommitIndex {
    //     vparts := strings.Split(v, "---")
	// 	if len(vparts)>1{
	// 		hashList = append(hashList, vparts[2])
	// 	}else{
	// 		hashList = append(hashList, v)
	// 	}
    // }

    // Use slices.Contains
	for _, hash := range *latestCommitIndex {

		hashParts := strings.Split(hash, "---")
		fmt.Println("stsaging hashparts of latest commit",hashParts)
		if len(hashParts)>1 && hashParts[2]==srcFileHash{
			(*newCommitIndexInfo)[srcFilePath] = "commit---"+string(blastCommitId)+"---"+hashParts[2]
			return nil
		}else if hash == srcFileHash {
			(*newCommitIndexInfo)[srcFilePath] = "commit---"+string(blastCommitId)+"---"+hash
			return nil
		} else {
			fmt.Println("File hash does not exist in the map.")
		}
	}

	// adding the entry in the map

	(*newCommitIndexInfo)[srcFilePath] = srcFileHash

	// after making the filepath --> file hash entry
	// i need to add the file in the stagin area

	stageFilePath := filepath.Join(cwd, ".fit", "stage", srcFileHash+".gz")
	destFile, err := os.Create(stageFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	zw := gzip.NewWriter(destFile)
	defer zw.Close()



	// reading the src file to compress it
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Use io.Copy instead of reading the entire file into memory
	_, err = io.Copy(zw, srcFile)
	if err != nil {
		return err
	}

	fmt.Println("Compression successful:", destFilePath)
	return nil
}




func CalculateHash(srcFilePath string) (string,error) {

	f, err := os.Open(srcFilePath)
	if err != nil {
		return "",err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "",err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash,nil

}





func EntryHashToFile(filePath,hashedFileName string) error {
	f, err := os.OpenFile(stageIndexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(filePath+" "+hashedFileName + "\n")
	return err
}




func CalculatePreviousHashes(lastCommitIndexInfo *SFileToHash) error {

	cwd,err := os.Getwd()
	if err!=nil{
		return err
	}

	bfile , err := os.ReadFile(".fit/HEAD/index.txt")
	if err!=nil{
		return err
	}

	// get the last commit index file to prepare a map path ----> hash
	blastIndexFile,err := os.ReadFile(filepath.Join(cwd,".fit","object",string(bfile),"index.txt"))
	if err!=nil{
		return err
	}

	json.Unmarshal(blastIndexFile,&lastCommitIndexInfo)

	// dirs , err := os.ReadDir(path.Join(cwd,".fit","object",previousCommitId))
	// if err!=nil{
	// 	return err
	// }

	// for _,dir := range dirs{
	// 	if dir.IsDir(){continue}
	// 	fmt.Println(dir.Name())

	// 	if dir.Name()=="index.txt"{
	// 		file,err := os.ReadFile(".fit/object/"+previousCommitId+"/"+dir.Name())
	// 		if err!=nil{
	// 			return err
	// 		}
	// 		json.Unmarshal(file,lastCommitIndexInfo)
	// 	}

	// }

	return nil
}
