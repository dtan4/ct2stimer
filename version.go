package main

import (
	"fmt"
)

var (
	// Version represents version number
	Version string
	// Revision represents the git commit hash at build time
	Revision string
)

func printVersion() {
	fmt.Println("ct2stimer version " + Version + ", build " + Revision)
}
