package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
)

type monitoredFile struct {
	file os.FileInfo
	tailProcess *os.Process
}

func MonitorDir(path string, maxTails int) {
	files, err := os.ReadDir(path)
	if err != nil {
		println("Unable to read directory, exiting.")
		os.Exit(1)
	}

	monitoredFiles := make([]monitoredFile, maxTails)
	filesForMonitoring := getFilesForMonitoring(files, maxTails)

	for i, f := range filesForMonitoring {
		if !f.IsDir() {
			finfo, err := f.Info()
			if err != nil {

				println("Unable to read file %s , skipping", f.Name())
				continue
			}

			monitoredFiles[i] = tailFile(filepath.Join(path,finfo.Name()), finfo)
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
	tails = sortMonitoredFilesByModTime(tails)
	if len(tails) == maxTails {
		mf := tails[len(tails) - 1]
		mf.tailProcess.Kill()
		mf.tailProcess.Wait()
	}
	tails[len(tails)-1] = tailFile(path,finfo)
	return tails
}

func tailFile(pathToFile string, file os.FileInfo) monitoredFile {

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
