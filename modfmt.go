package main

import (
	"os"

	"golang.org/x/mod/modfile"
)

func MergeRequires(goModFilename string) ([]byte, error) {
	contents, err := os.ReadFile(goModFilename)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.ParseLax(goModFilename, contents, nil)
	if err != nil {
		return nil, err
	}

	if err := mergeRequires(mod); err != nil {
		return nil, err
	}

	updatedContents, err := mod.Format()
	if err != nil {
		return nil, err
	}

	return updatedContents, nil
}

func mergeRequires(mod *modfile.File) (err error) {
	defer func() {
		if r := recover(); r != nil {
			possibleErr, ok := r.(error)
			if ok {
				err = possibleErr
			}
		}
	}()

	allRequires := make([]modfile.Require, len(mod.Require))
	for i, reqs := range mod.Require {
		// Save all the requires to a new slice
		allRequires[i] = *reqs

		// while removing them from the original slice
		mod.DropRequire(reqs.Mod.Path)
	}

	// Cleanup the modfile
	// This removes the empty require blocks
	mod.Cleanup()

	// Add the requires back to the modfile
	for _, reqs := range allRequires {
		mod.AddNewRequire(reqs.Mod.Path, reqs.Mod.Version, reqs.Indirect)
	}
	mod.Cleanup()

	// Sort the require blocks
	mod.SortBlocks()
	mod.Cleanup()

	// Set the require blocks to separate indirects
	mod.SetRequireSeparateIndirect(mod.Require)
	mod.Cleanup()

	return
}
