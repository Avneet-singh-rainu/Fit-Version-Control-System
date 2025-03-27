package main

import (
	"fmt"
	"os"
)

func Init(){
	messaageForFitign := ".fitign\n.mod\n.exe\n.git"
	_,err := os.Getwd()
	directories := []string{"object","HEAD","commit","stage"}

	if(err!=nil){
		fmt.Print(err)
		return
	}

	//dirName := path.Join(currDirectory,".fit")

	e := os.Mkdir(".fit",0750)

	if(e!=nil){
		fmt.Print(e)
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

	fmt.Print("fit initiated the control system has been established...")

}
