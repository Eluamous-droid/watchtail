package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var testFilesDir = "testDir"

func TestNewFileCreatedEmptySlice(t *testing.T) {

	os.Mkdir(testFilesDir, 0755)
	defer os.RemoveAll(testFilesDir)

	maxTails := 10
	file1 := createFile("test1")
	mfs := make([]monitoredFile, 0, maxTails)
	mfs = newFileCreated(file1, maxTails, mfs)

	if len(mfs) != 1 {
		println("len is not 1, it is: " + strconv.FormatInt(int64(len(mfs)), 10))
		t.Fail()
	}

	killAllTails(mfs)
}

func TestNewFileCreatedFullSlice(t *testing.T) {
	os.Mkdir(testFilesDir, 0755)
	defer os.RemoveAll(testFilesDir)

	maxTails := 2
	file1 := createFile(filepath.Join(testFilesDir, "test1"))
	file2 := createFile(filepath.Join(testFilesDir, "test2"))
	file3 := createFile(filepath.Join(testFilesDir, "test3"))
	mfs := make([]monitoredFile, 0, maxTails)
	mfs = newFileCreated(file1, maxTails, mfs)
	mfs = newFileCreated(file2, maxTails, mfs)
	mfs = newFileCreated(file3, maxTails, mfs)

	if len(mfs) != 2 {
		println("len is not 2, it is: " + strconv.FormatInt(int64(len(mfs)), 10))
		t.Fail()
	}

	killAllTails(mfs)
}

func createFile(name string) string {
	fileInput := []byte("File name is: " + name)
	err := os.WriteFile(name, fileInput, 0644)

	if err != nil {
		panic(err)
	}

	return name
}
