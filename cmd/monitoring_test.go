package cmd

import (
	"os"
	"testing"
	"time"
)

const dirName ="testDir"

func TestDirEntriesSortedByModDate(t *testing.T){

	originalFiles := []os.FileInfo{createFile("test1", false,0), createFile("test5", false,1), createFile("test2", false,0), createFile("test11", false,0)}


	newfiles = sortFilesByModTime(originalFiles)

	if originalFiles[1].Name() != newFiles[1].Name(){
		t.Fail()
	}
}


func createDir() {
	err := os.Mkdir(dirName,0644)
	if err != nil {
		panic(err)
	}
}

func createFile(name string, isDir bool, howLongAgo int) MockFileInfo{
	t := time.Now().Add(-time.Hour * time.Duration(howLongAgo))
	return MockFileInfo{FileName: name, IsDirectory: isDir, LastModTime: t}

}

