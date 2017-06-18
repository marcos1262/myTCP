package myTCP

import (
	"fmt"
	"os"
)

// Enables debugging messages
const DEBUG = false

// When an error is thrown
const FATAL_ERROR = "Fatal error: %s"


func debug(s string) {
	if DEBUG {
		fmt.Println("DEBUG: " + s)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, FATAL_ERROR, err.Error())
		os.Exit(1)
	}
}