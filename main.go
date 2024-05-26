/*
Copyright © 2024 Mikkel Hansen
*/
package main

import (
	"flag"
	"fmt"
	"os"
)

type excludes []string

var excludedFiles excludes

func main() {

	maxTails := flag.Int("m", 10, "Max tails running at once")
	flag.Var(&excludedFiles, "exclude", "File to be excluded")
	flag.Parse()
	if len(flag.Args()) < 1 {
		println("Must provide path to monitor")
		os.Exit(1)
	}
	MonitorDir(flag.Arg(0), *maxTails)
}

func (n *excludes) Set(value string) error {
	*n = append(*n, value)
	return nil
}

func (n *excludes) String() string {
	return fmt.Sprintf("%s", *n)
}
