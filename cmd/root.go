package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"sort"

	"github.com/fsnotify/fsnotify"
)

func MonitorDir(path string, maxTails int){

		queue, counter := tailExistingFiles(path, maxTails)
		defer killAllTails(queue)
		monitorDirectory(path, queue, counter, maxTails)

}

func monitorDirectory(path string, queue chan *os.Process, counter int, maxTails int){

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	err = watcher.Add(path)
	if err != nil {
		panic(err)
	}
	// We should never leave this function unless the program ends
	for {
		select{
		case err, ok := <-watcher.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			println("ERROR: %s", err)

		case event,ok := <-watcher.Events:
			if !ok{
				println("event was not ok") 
				return
			}

			if event.Has(fsnotify.Create){
				if counter == maxTails{
					process := <- queue
					process.Kill()
					counter--
				}
				queue <- tailFile(event.Name)
				counter++
			}
		}	
	}
}

func killAllTails(queue chan *os.Process){

	for p := range queue{
		p.Kill()
		
	}

}
func printFiles(files []fs.DirEntry) {
	for _, file := range files {
		fileInfo, _:= file.Info()
		fmt.Println(file.Name(), fileInfo.ModTime())
	}
}
func tailExistingFiles(path string, maxLength int) (tails chan *os.Process, tailCount int){
	files, err := os.ReadDir(path)
	queue := make(chan *os.Process, maxLength)
	counter := 0
	if err != nil {
		panic(err)
	}
	files = removeDirs(files)
	files = sortFilesByModTime(files)

	sliceSize := getSmallestInt(len(files), maxLength)
	filesSlice := files[:sliceSize]

	for _, file := range filesSlice{
		if !file.IsDir() {
			queue <- tailFile(path + "/" + file.Name())
			counter++
		}

	}
	
	return queue, counter
}
func getSmallestInt(a int, b int) int {
	if a > b{ return b}
	return a

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

func sortFilesByModTime(files []os.DirEntry) []os.DirEntry{
	sort.Slice(files, func(i,j int) bool{
		fileI, err := files[i].Info()
		if err != nil {
			panic(err)
		}
		fileJ, err := files[j].Info()
		if err != nil {
			panic(err)
		}
		return fileI.ModTime().After(fileJ.ModTime())

	})
	return files
}
func removeDirs(s []os.DirEntry) []os.DirEntry{
	var dirless []os.DirEntry
	for _,file := range s{
		if !file.IsDir() {
			dirless = append(dirless, file)
		}
	}

	return dirless
}
