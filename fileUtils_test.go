package main

import (
	"os"
	"testing"
	"time"
)

func TestDirEntriesSortedByModDateOldestFirst(t *testing.T) {

	originalFiles := []os.DirEntry{createFile("test1", false, 2), createFile("test2", false, 1), createFile("test3", false, 5), createFile("test41", false, 10)}

	newFiles := make([]os.DirEntry, len(originalFiles))
	copy(newFiles, originalFiles)

	sortFilesByModTime(newFiles)

	if originalFiles[1].Name() != newFiles[3].Name() {
		t.Fail()
	}
}


func TestIneligbleFilesFiledInExcluded(t *testing.T) {
	f1 := createFile("ThisIsAtestFile", false, 0)
	fileNamePatterns := excludes{"test", "notMatchingPattern"}
	excludedFiles = fileNamePatterns

	if isEligibleFile(f1.MockInfo) {
		t.Fail()
	}
}
func TestIneligbleFilesIsDir(t *testing.T) {
	f1 := createFile("ThisIsAtestFile", true, 0)
	fileNamePatterns := excludes{"NotExcluded", "notMatchingPattern"}
	excludedFiles = fileNamePatterns

	if isEligibleFile(f1.MockInfo) {
		t.Fail()
	}
}

func TestIneligbleFilesFiledNotInExcluded(t *testing.T) {
	f1 := createFile("ThisIsAtestFile", false, 0)
	f2 := createFile("Sure", false, 0)
	fileNamePatterns := excludes{"NotExcluded", "ForSure"}
	excludedFiles = fileNamePatterns

	if !isEligibleFile(f1.MockInfo) {
		t.Fail()
	}
	if !isEligibleFile(f2.MockInfo) {
		t.Fail()
	}
}
func TestRemovesIneligibleFiles(t *testing.T) {
	f1 := createFile("ThisIsAtestFile", false, 0)
	f2 := createFile("Sure", false, 0)
	fileNamePatterns := excludes{"test", "ForSure"}
	excludedFiles = fileNamePatterns

	fileInput := []os.DirEntry{f1, f2}

	fileOutput := removeIneligibleFiles(fileInput)

	for _, f := range fileOutput {
		if f.Name() == f1.MockInfo.FileName {
			t.Fail()
		}
	}
}

func createFile(name string, isDir bool, howLongAgo int) MockDirEntry {
	t := time.Now().Add(-time.Hour * time.Duration(howLongAgo))
	mfi := MockFileInfo{FileName: name, IsDirectory: isDir, LastModTime: t}
	return MockDirEntry{FileName: name, IsDirectory: isDir, MockInfo: mfi}

}
