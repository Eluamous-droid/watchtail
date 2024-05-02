package cmd

import (
	"os"
	"os/exec"

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
		defer killAllTails(queue, counter)

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

func killAllTails(queue chan *os.Process, counter int){

	for i:= 0; i < counter; i++{
		p := <- queue
		p.Kill()
		
	}

}

func tailExistingFiles(path string, maxLength int) (tails chan *os.Process, tailCount int){
	files, err := os.ReadDir(path)
	queue := make(chan *os.Process, maxLength)
	counter := 0
	if err != nil {
		panic(err)
	}

	for _, file := range files{
		if !file.IsDir() {
			queue <- tailFile(path + "/" + file.Name())
			counter++
		}

		if counter >= maxLength{
			oldest := <- queue
			oldest.Kill()
			counter--
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


