package main

import (
	"fmt"
	"os"
)

func Init(){
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

	fmt.Print("fit initiated the control system has been established...")

}
