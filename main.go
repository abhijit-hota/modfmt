package main

import (
	"fmt"
	"log"
	"os"
	"slices"
)

const gomodName = "go.mod"

func main() {
	updatedContents, err := MergeRequires(gomodName)
	if err != nil {
		log.Fatal(err)
	}

	// get arguments
	if slices.Contains(os.Args, "--replace") {
		// write updated contents to go.mod
		if err = os.WriteFile(gomodName, updatedContents, 0o644); err != nil {
			log.Fatal(fmt.Errorf("failed to write updated go.mod: %v", err))
		}

		return
	}

	fmt.Println(string(updatedContents))
}
