package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
)

type SFileToHash struct{
	Files map[string]string
}




// CopyDirAndCompress recursively copies and compresses a directory.
func CopyDirAndCompress(src, dest string , fileToHash , previousFileToHash *SFileToHash) error {
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
			if err := CopyDirAndCompress(srcPath, destPath,fileToHash,previousFileToHash); err != nil {
				return err
			}
		} else {
			// Compress individual files

            fmt.Println("Staging file:", srcPath, "->", hashDestFilePath)
			if err := CopyFileAndCompress(srcPath, destPath,fileToHash,previousFileToHash); err != nil {
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
func CopyFileAndCompress(srcFilePath, destFilePath string, fileToHash,previousFileToHash *SFileToHash) error {

	// calculating the hash of the current file so that i can compare it with the entries in the preciousFileToHAsh
	destFileHash,err := CalculateHash(srcFilePath)
	if err != nil {
		return err
	}

	var hashList []string
    for _, v := range previousFileToHash.Files {
        hashList = append(hashList, v)
    }

    // Use slices.Contains
    if slices.Contains(hashList, destFileHash) {
		for _, hash := range previousFileToHash.Files {
			if hash == destFileHash {
				fileToHash.Files[srcFilePath] = "commit-"+previousCommitId
				return nil
			}
		}

    } else {
        fmt.Println("File hash does not exist in the map.")
    }

	// adding the entry in the map

	fileToHash.Files[srcFilePath] = destFileHash



	// after making the filepath --> file hash entry
	// i need to add the file in the stagin area

	cwd,err:=os.Getwd()
	if err != nil {
		return err
	}

	stageFilePath := filepath.Join(cwd, ".fit", "stage", destFileHash+".gz")
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




func CalculatePreviousHashes(previousFileToHash *SFileToHash) error {

	bfile , err := os.ReadFile(".fit/HEAD/index.txt")
	if err!=nil{
		return err
	}
	previousCommitId = string(bfile)

	cwd,err := os.Getwd()
	if err!=nil{
		return err
	}

	dirs , err := os.ReadDir(path.Join(cwd,".fit","object",previousCommitId))
	if err!=nil{
		return err
	}

	for _,dir := range dirs{
		if dir.IsDir(){continue}
		fmt.Println(dir.Name())

		if dir.Name()=="index.txt"{
			file,err := os.ReadFile(".fit/object/"+previousCommitId+"/"+dir.Name())
			if err!=nil{
				return err
			}
			json.Unmarshal(file,previousFileToHash)
		}

	}

	return nil
}
