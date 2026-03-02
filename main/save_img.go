package main

import (
	"github.com/sw965/bmgt"
)

func main() {
	return
	err := bmgt.FetchAndSaveImage("Dark Magician Girl", "img")
	if err != nil {
		panic(err)
	}
}
