package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retailer",
	Short: "Monitors a folder and will run tail -f on new files",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		queue, counter := tailExistingFiles(args[0], 10)
		defer killAllTails(queue)
		monitorDirectory(args[0], queue, counter, 10)

	},

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.retailer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	files = sortFilesByModTime(files)
	filesSlice := files[:maxLength]

	for _, file := range filesSlice{
		if !file.IsDir() {
			queue <- tailFile(path + "/" + file.Name())
			counter++
		}

	}
	
	return queue, counter
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

