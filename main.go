package main

import (
	"fmt"
	"log"
)

const gomodName = "go.mod"

func main() {

	updatedContents, err := MergeRequires(gomodName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(updatedContents))
}
