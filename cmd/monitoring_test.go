package cmd

import (
	"os"
	"testing"
	"time"
)

const dirName ="testDir"

func TestDirEntriesSortedByModDateNewestFirst(t *testing.T){

	originalFiles := []os.DirEntry{createFile("test1", false,2), createFile("test2", false,1), createFile("test3", false,5), createFile("test41", false,10)}

	newFiles := make ([]os.DirEntry, len(originalFiles))
	copy(newFiles,originalFiles)

	sortFilesByModTime(newFiles)

	if originalFiles[1].Name() != newFiles[0].Name(){
		println(originalFiles[1].Name() + " " + newFiles[0].Name())
		t.Fail()
	}
}


func createDir() {
	err := os.Mkdir(dirName,0644)
	if err != nil {
		panic(err)
	}
}

func createFile(name string, isDir bool, howLongAgo int) MockDirEntry{
	t := time.Now().Add(-time.Hour * time.Duration(howLongAgo))
	mfi := MockFileInfo{FileName: name, IsDirectory: isDir, LastModTime: t}
	return MockDirEntry{FileName: name, IsDirectory: isDir, MockInfo: mfi}

}

