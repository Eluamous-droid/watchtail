package main

import (
	"os"
	"sort"
	"strings"
)

func getFilesForMonitoring(files []os.DirEntry, maxLength int) []os.DirEntry {
	files = removeIneligibleFiles(files)
	files = sortDirEntryByModTime(files)

	sliceSize := getSmallestInt(len(files), maxLength)
	filesSlice := files[len(files)-sliceSize:]

	return filesSlice
}

func sortDirEntryByModTime(files []os.DirEntry) []os.DirEntry {
	sort.Slice(files, func(i, j int) bool {
		fileI, err := files[i].Info()
		if err != nil {
			println("Unable to read file %s , while sorting", fileI.Name())
			return true
		}
		fileJ, err := files[j].Info()
		if err != nil {
			println("Unable to read file %s , while sorting", fileJ.Name())
			return true
		}
		return fileI.ModTime().Before(fileJ.ModTime())

	})
	return files
}

func removeDeletedFileFromMonitoredFileSlice(files []monitoredFile, deletedFile string) []monitoredFile {
	newFiles := make([]monitoredFile, 0, 0)
	for _, mf := range files {
		if fi, err := mf.file.Info(); err != nil {
			untailFile(mf)
		} else {
			if fi.Name() == deletedFile {
				untailFile(mf)
			} else {
				newFiles = append(newFiles, mf)
			}
		}
	}
	return newFiles
}

func sortMonitoredFilesByModTime(files []monitoredFile) []monitoredFile {
	sort.Slice(files, func(i, j int) bool {
		fileI, err := files[i].file.Info()
		if err != nil {
			println("Unable to read file %s , while sorting", fileI.Name())
			return true
		}
		fileJ, err := files[j].file.Info()
		if err != nil {
			println("Unable to read file %s , while sorting", fileJ.Name())
			return true
		}
		return fileI.ModTime().After(fileJ.ModTime())

	})
	return files
}

func removeIneligibleFiles(s []os.DirEntry) []os.DirEntry {
	var eligibleFiles []os.DirEntry
	for _, file := range s {
		fi, err := file.Info()
		if err != nil {
			println("File is unreadable: ", file.Name())
		}
		if isEligibleFile(fi) {
			eligibleFiles = append(eligibleFiles, file)
		}
	}
	return eligibleFiles
}

func isEligibleFile(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	fn := fi.Name()
	for _, s := range excludedFiles {
		if strings.Contains(fn, s) {
			println("Name is in excludes")
			return false
		}
	}

	return true
}

func untailFile(mf monitoredFile) {
	mf.tailProcess.Kill()
	mf.tailProcess.Wait()
}

func getSmallestInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
