package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestDirEntriesSortedByModDateOldestFirst(t *testing.T) {
	defer teardown()

	originalFiles := []os.DirEntry{createFileMock("test1", false, 2), createFileMock("test2", false, 1), createFileMock("test3", false, 5), createFileMock("test41", false, 10)}

	newFiles := make([]os.DirEntry, len(originalFiles))
	copy(newFiles, originalFiles)

	sortDirEntryByModTime(newFiles)

	if originalFiles[1].Name() != newFiles[3].Name() {
		t.Fail()
	}
}

func TestSortMonitoredFilesByModTimeNewestFirst(t *testing.T) {
	defer teardown()

	originalFiles := []monitoredFile{createMonitoredFileMock("test1", false, 3), createMonitoredFileMock("test2", false, 5), createMonitoredFileMock("test3", false, 1), createMonitoredFileMock("test41", false, 10)}

	newFiles := make([]monitoredFile, len(originalFiles))
	copy(newFiles, originalFiles)

	sortMonitoredFilesByModTime(newFiles)

	if originalFiles[0].file.Name() != newFiles[1].file.Name() {
		t.Fail()
	}
}

func TestIneligbleFilesFiledInExcluded(t *testing.T) {
	defer teardown()
	f1 := createFileMock("ThisIsAtestFile", false, 0)
	fileNamePatterns := excludes{"test", "notMatchingPattern"}
	excludedFiles = fileNamePatterns

	if isEligibleFile(f1.MockInfo) {
		t.Fail()
	}
}
func TestIneligbleFilesIsDir(t *testing.T) {
	defer teardown()
	f1 := createFileMock("ThisIsAtestFile", true, 0)
	fileNamePatterns := excludes{"NotExcluded", "notMatchingPattern"}
	excludedFiles = fileNamePatterns

	if isEligibleFile(f1.MockInfo) {
		t.Fail()
	}
}

func TestIneligbleFilesFiledNotInExcluded(t *testing.T) {
	defer teardown()
	f1 := createFileMock("ThisIsAtestFile", false, 0)
	f2 := createFileMock("Sure", false, 0)
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
	defer teardown()
	f1 := createFileMock("ThisIsAtestFile", false, 0)
	f2 := createFileMock("Sure", false, 0)
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

func TestFileDeleted(t *testing.T) {

	os.Mkdir(testFilesDir, 0755)
	defer os.RemoveAll(testFilesDir)

	maxTails := 3
	file1 := createFile(filepath.Join(testFilesDir, "test1"))
	file2 := createFile(filepath.Join(testFilesDir, "test2"))
	file3 := createFile(filepath.Join(testFilesDir, "test3"))

	mfs := make([]monitoredFile, 0, 0)
	mfs = newFileCreated(file1, maxTails, mfs)
	mfs = newFileCreated(file2, maxTails, mfs)

	err := os.Remove(file2)
	if err != nil {
		println("Unable to remove file " + file2)
		t.Fail()
	}

	mfs = removeDeletedFileFromMonitoredFileSlice(mfs, "test2")

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

func teardown() {
	excludedFiles = excludes{}
}

func createFileMock(name string, isDir bool, howLongAgo int) MockDirEntry {
	t := time.Now().Add(-time.Hour * time.Duration(howLongAgo))
	mfi := MockFileInfo{FileName: name, IsDirectory: isDir, LastModTime: t}
	return MockDirEntry{FileName: name, IsDirectory: isDir, MockInfo: mfi}
}

func createMonitoredFileMock(name string, isDir bool, howLongAgo int) monitoredFile{
	t := time.Now().Add(-time.Hour * time.Duration(howLongAgo))
	mfi := MockFileInfo{FileName: name, IsDirectory: isDir, LastModTime: t}
	mde :=  MockDirEntry{FileName: name, IsDirectory: isDir, MockInfo: mfi}
	return monitoredFile{file: mde, tailProcess: nil}
}
