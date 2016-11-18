package main

import "os"

const name string = "nrec-moody"

var (
	version  = "0.1.0"
	revision string
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
