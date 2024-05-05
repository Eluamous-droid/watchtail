/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"os"
	"retailer/cmd"
)

func main() {

	maxTails := flag.Int("Max tails", 10, "Max tails running at once")
	cmd.MonitorDir(os.Args[1],*maxTails)
}
