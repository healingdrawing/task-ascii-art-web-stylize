package main

import (
	"flag"
	"log"
)

var (
	debug        bool
	defaultflags = 0
	logflags     = defaultflags
)

// executed before "main()" to provide flag "-debug" functionality for server log.
func init() {
	flag.BoolVar(&debug, "debug", false, "set to true for debug logging")
}

func main() {
	flag.Parse()
	if debug {
		logflags = log.LstdFlags | log.Lshortfile
	}
	log.SetFlags(logflags)

	startServer()
}
