package main

import (
	"errors"
	"os"
)

func FitExists() error {
	_, err := os.Open(".fit")
	if err != nil {
		return err
	}
	return nil
}

func StageExixts() error {
	// dirs, err := os.ReadDir(".fit")
	// if err != nil {
	// 	return err
	// }

	// for _,dir := range dirs{
	// 	if dir.Name()=="stage"{
	// 		return nil
	// 	}
	// }
	// return errors.New("---no files found inside the staging area, there is nothing to commit---")
	_,err := os.ReadDir(".fit/stage")
	if err!=nil{
		return errors.New("---no files found inside the staging area, there is nothing to commit---")
	}
	return nil
}

func CommitsExists() error {

	dirs , err := os.ReadDir(".fit/Commit")
	if err!=nil{
		return err
	}

	if len(dirs)==0{
		return errors.New("no previous versions found, you are currently on the latest version")
	}

	return nil
}
