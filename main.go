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

	maxTails := flag.Int("m", 10, "Max tails running at once")
	flag.Parse()
	if len(flag.Args()) < 1 {
		println("Must provide path to monitor")
		os.Exit(1)
	}
	cmd.MonitorDir(flag.Arg(0), *maxTails)
}
