package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	gomodName := "go.mod"
	if len(os.Args) > 1 {
		gomodName = os.Args[1]
	}

	updatedContents, err := MergeRequires(gomodName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(updatedContents))
}
