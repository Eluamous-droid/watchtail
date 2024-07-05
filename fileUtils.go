package main

import (
	"os"
	"sort"
	"strings"
)

func getFilesForMonitoring(files []os.DirEntry, maxLength int) []os.DirEntry {
	files = removeIneligibleFiles(files)
	files = sortFilesByModTime(files)

	sliceSize := getSmallestInt(len(files), maxLength)
	filesSlice := files[len(files)-sliceSize:]

	return filesSlice
}

func sortFilesByModTime(files []monitoredFile) []monitoredFile{
	sort.Slice(files, func(i, j int) bool {
		fileI:= files[i].file
		fileJ:= files[j].file
		return fileI.ModTime().Before(fileJ.ModTime())

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
			return false
		}
	}

	return true
}

func getSmallestInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
