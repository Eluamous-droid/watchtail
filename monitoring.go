package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
)

type monitoredFile struct {
	file os.DirEntry
	tailProcess *os.Process
}

func MonitorDir(path string, maxTails int) {
	files, err := os.ReadDir(path)
	if err != nil {
		println("Unable to read directory, exiting.")
		os.Exit(1)
	}

	monitoredFiles := make([]monitoredFile, 0, maxTails)
	filesForMonitoring := getFilesForMonitoring(files, maxTails)

	for _, f := range filesForMonitoring {
		if !f.IsDir() {
			finfo, err := f.Info()
			if err != nil {

				println("Unable to read file %s , skipping", f.Name())
				continue
			}

			monitoredFiles = append(monitoredFiles, tailFile(filepath.Join(path,finfo.Name()), f))
		}
	}

	defer killAllTails(monitoredFiles)
	startWatching(path, monitoredFiles, maxTails)

}

func startWatching(path string, tails []monitoredFile, maxTails int) {

	w := watcher.New()
	w.FilterOps(watcher.Create)
	defer w.Close()
	err := w.Add(path)
	if err != nil {
		println("Unable to add watcher")
		os.Exit(1)
	}
	go func() {
		for {
			select {
			case err, ok := <-w.Error:
				if !ok { // Channel was closed (i.e. Watcher.Close() was called).
					return
				}
				println("ERROR: %s", err)

			case event, ok := <-w.Event:
				if !ok {
					println("event was not ok")
					return
				}
				tails = newFileCreated(event.Path, maxTails, tails)
			}
		}
	}()

	// We should never leave this function unless the program ends
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

}

func newFileCreated(path string, maxTails int, tails []monitoredFile) []monitoredFile {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		println("New file cannot be read: ", path)
	return tails
	}
	finfo, _ := f.Stat()

	if !isEligibleFile(finfo) {
	return tails
	}
	
	if len(tails) == maxTails {
		tails = sortMonitoredFilesByModTime(tails)
		mf := tails[len(tails) - 1]
		mf.tailProcess.Kill()
		mf.tailProcess.Wait()
		tails[len(tails)-1] = tailFile(path, fs.FileInfoToDirEntry(finfo))
	}else{	
		tails = append(tails, tailFile(path, fs.FileInfoToDirEntry(finfo)))
	}
	return tails
}

func tailFile(pathToFile string, file os.DirEntry) monitoredFile {

	app := "tail"
	args := "-f"

	cmd := exec.Command(app, args, pathToFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	mf := monitoredFile{file: file, tailProcess: cmd.Process}

	return mf
}

func killAllTails(moniteredFiles []monitoredFile) {

	for _,mf := range moniteredFiles {
		mf.tailProcess.Kill()
		mf.tailProcess.Wait()
	}
}
