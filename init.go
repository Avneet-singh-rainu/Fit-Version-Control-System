package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Init(){

	messaageForFitign := ".fitign\n.mod\n.exe\n.git"
	_,err := os.Getwd()
	directories := []string{"object","HEAD","commit"}

	if(err!=nil){
		fmt.Print(err)
		return
	}

	//dirName := path.Join(currDirectory,".fit")

	e := os.Mkdir(".fit",0750)
	// if fit is already initiated then return
	if(e!=nil){
		color.Red("already initiated provide other commands")
		return
	}

	for _,dir := range directories {
		e:=os.Mkdir(".fit/"+dir , 0750)
		if(e!=nil){
			fmt.Print(e)
		}
	}


	// i need to create a .fitign file where i can add filenames,dirs and file extensions to avoid staging and all
	err = os.WriteFile(".fitign",[]byte(messaageForFitign),0750)

	if(err != nil){
		fmt.Print("Error creating the fitign file",err)
		return
	}

	fmt.Println("fit initiated the control system has been established...")
	color.Yellow("--------Make sure you add large files in the fitign file.--------")

}
