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
	os.RemoveAll(testFilesDir)
}

func TestNewFileCreatedFullSliceDoesntRemoveFirstAddedFile(t *testing.T) {
	os.Mkdir(testFilesDir, 0755)
	defer os.RemoveAll(testFilesDir)

	maxTails := 2
	file1 := createFile(filepath.Join(testFilesDir, "test1"))
	file2 := createFile(filepath.Join(testFilesDir, "test2"))
	file3 := createFile(filepath.Join(testFilesDir, "test3"))

	mfs := make([]monitoredFile, 0, 0)
	mfs = newFileCreated(file1, maxTails, mfs)
	mfs = newFileCreated(file2, maxTails, mfs)

	f, err := os.OpenFile(file1, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file1ExtraInput := "file1 extra input"
	if _, err = f.WriteString(file1ExtraInput); err != nil {
		panic(err)
	}

	mfs = newFileCreated(file3, maxTails, mfs)

	if len(mfs) != 2 {
		println("len is not 2, it is: " + strconv.FormatInt(int64(len(mfs)), 10))
		t.Fail()
	}

	for _, mf := range mfs {
		if mf.file.Name() == "test2" {
			println(file2 + " should have been removed")
			t.Fail()
		}

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
