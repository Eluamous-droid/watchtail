/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"retailer/cmd"
)

type excludes []string

func main() {

	var excludedFiles excludes

	maxTails := flag.Int("m", 10, "Max tails running at once")
	flag.Var(&excludedFiles, "exclude","File to be excluded")
	flag.Parse()
	if len(flag.Args()) < 1 {
		println("Must provide path to monitor")
		os.Exit(1)
	}
fmt.Println("Excludes: ", excludedFiles)
	cmd.MonitorDir(flag.Arg(0), *maxTails)
}

func (n *excludes) Set(value string) error {
    *n = append(*n, value)
    return nil
}

func (n *excludes) String() string {
    return fmt.Sprintf("%s", *n)
}
