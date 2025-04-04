package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"golang.org/x/mod/modfile"
)

const (
	testFileName        = "testdata/go.mod"
	updatedTestFileName = "testdata/updated_go.mod"
)

func TestRun(t *testing.T) {
	t.Run("Should output error with invalid argument", func(t *testing.T) {
		args := []string{"--invalid"}
		output := bytes.NewBuffer(nil)

		err := run(args, output)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrFailedToParseFlags) {
			t.Fatalf("Expected error to be %v, got %v", ErrFailedToParseFlags, err)
		}

		if output.Len() == 0 {
			t.Fatal("Expected output, got nothing")
		}

		if !strings.Contains(output.String(), "flag provided but not defined: -invalid") {
			t.Fatalf("Expected output to contain 'flag provided but not defined: -invalid', got '%s'", output.String())
		}
	})
}

func TestMergeRequires(t *testing.T) {
	// fmt the go.mod file
	updatedContents, err := MergeRequires(testFileName)
	if err != nil {
		t.Fatal(err)
	}

	// create a new file with the updated contents
	newFile, err := os.Create(updatedTestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer newFile.Close()

	if _, err = newFile.Write(updatedContents); err != nil {
		t.Fatal(err)
	}

	// parse the old go.mod file
	oldContents, err := os.ReadFile(testFileName)
	if err != nil {
		t.Fatal(err)
	}
	oldmod, err := modfile.ParseLax(testFileName, oldContents, nil)
	if err != nil {
		t.Fatal(err)
	}

	// parse the updated go.mod file
	newContents, err := os.ReadFile(updatedTestFileName)
	if err != nil {
		t.Fatal(err)
	}
	newmod, err := modfile.ParseLax(updatedTestFileName, newContents, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check if both mod reqs have the same length
	if len(oldmod.Require) != len(newmod.Require) {
		t.Errorf("Require length mismatch: %d != %d", len(oldmod.Require), len(newmod.Require))
	}

	// Check if both mod reqs have the same content
	for _, req := range oldmod.Require {
		found := false
		for _, newReq := range newmod.Require {
			if req.Mod.Path == newReq.Mod.Path && req.Mod.Version == newReq.Mod.Version {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Require mismatch: %s@%s not found in updated go.mod", req.Mod.Path, req.Mod.Version)
		}
	}
}
