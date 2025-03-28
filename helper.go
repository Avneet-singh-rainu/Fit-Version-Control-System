package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MoveDir moves an entire directory and deletes the original after copying
func MoveDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively move subdirectories
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Move individual files
			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	// Remove the source directory after moving all contents
	return os.RemoveAll(src)
}









// CopyDirAndDecompress recursively copies and decompresses a directory.
func CopyDirAndDecompress(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, 0750); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, strings.TrimSuffix(entry.Name(), ".gz")) // Remove `.gz` before decompressing

		if entry.IsDir() {
			// Recursively decompress subdirectories
			if err := CopyDirAndDecompress(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Decompress individual files
			if err := CopyFileAndDecompress(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}













// GetFitignFiles reads the `.fitign` file and retrieves ignored files and directories.
func GetFitignFiles() (ignoreFiles []string, ignoreDirs []string, err error) {
	file, err := os.Open(".fitign")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	rawEntries := strings.Split(string(content), "\n")

	for _, entry := range rawEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Check if entry is a directory by detecting `/`
		if strings.Contains(entry, "/") {
			ignoreDirs = append(ignoreDirs, entry)
		} else {
			ignoreFiles = append(ignoreFiles, entry)
		}
	}

	fmt.Println("Ignored Directories:", ignoreDirs)
	fmt.Println("Ignored Files:", ignoreFiles)

	return ignoreFiles, ignoreDirs, nil
}






















// CopyFileAndDecompress decompresses a gzip file and saves it.
func CopyFileAndDecompress(srcFilePath, destFilePath string) error {
    // Check if the file actually has a .gz extension
    if !strings.HasSuffix(srcFilePath, ".gz") {
        return fmt.Errorf("file %s is not a gzip compressed file", srcFilePath)
    }

    // Open the compressed source file
    srcCompressedFile, err := os.Open(srcFilePath)
    if err != nil {
        return fmt.Errorf("error opening compressed file: %v", err)
    }
    defer srcCompressedFile.Close()

    // Create the gzip reader
    zr, err := gzip.NewReader(srcCompressedFile)
    if err != nil {
        return fmt.Errorf("error creating gzip reader: %v", err)
    }
    defer zr.Close()

    // Remove .gz extension for the destination file
    destFilePath = strings.TrimSuffix(destFilePath, ".gz")

    // Create the destination file
    destFile, err := os.Create(destFilePath)
    if err != nil {
        return fmt.Errorf("error creating destination file: %v", err)
    }
    defer destFile.Close()

    // Copy decompressed content into the destination file
    if _, err := io.Copy(destFile, zr); err != nil {
        return fmt.Errorf("error copying decompressed data: %v", err)
    }

    fmt.Println("✅ Decompression successful:", destFilePath)
    return nil
}

























// CopyFile copies a file from src to dest.
func CopyFile(srcFilePath, destFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}

// CopyDir copies a directory and its contents.
func CopyDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, 0750); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy individual files
			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}
