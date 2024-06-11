package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/radovskyb/watcher"
)

func MonitorDir(path string, maxTails int) {
	files, err := os.ReadDir(path)
	if err != nil {
		println("Unable to read directory, exiting.")
		os.Exit(1)
	}

	monitoredFiles := make([]*os.Process, maxTails)
	counter := 0
	filesForMonitoring := getFilesForMonitoring(files, maxTails)

	for _, f := range filesForMonitoring {
		if !f.IsDir() {
			finfo, err := f.Info()
			if err != nil {

				println("Unable to read file %s , skipping", f.Name())
				continue
			}

			monitoredFiles = append(monitoredFiles, tailFile(path + finfo.Name()))
			counter++
		}
	}

	defer killAllTails(monitoredFiles)
	startWatching(path, monitoredFiles, counter, maxTails)

}

func startWatching(path string, tails []*os.Process, counter int, maxTails int) {

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
				counter = newFileCreated(event.Path, counter, maxTails, tails)
			}
		}
	}()

	// We should never leave this function unless the program ends
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

}

func newFileCreated(path string, counter int, maxTails int, tails []*os.Process) int {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		println("New file cannot be read: ", path)
		return counter
	}
	fi, _ := f.Stat()

	if !isEligibleFile(fi) {
		return counter
	}
	if counter == maxTails {
		process := tails[counter - 1]
		process.Kill()
		process.Wait()
		counter--
	}
	tails[counter - 1] = tailFile(path)
	counter++
	sortFilesByModTime(tails)
	return counter
}

func tailFile(filePath string) *os.Process {

	app := "tail"
	args := "-f"

	cmd := exec.Command(app, args, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd.Process
}

func killAllTails(moniteredFiles []*os.Process) {

	for _,p := range moniteredFiles {
		p.Kill()
		p.Wait()
	}
}
