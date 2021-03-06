package main

import (
	"fmt"
	"os"
	"bufio"
)

// Enables debugging messages
const DEBUG = true

// When an error is thrown
const FATAL_ERROR = "Fatal error: %s"

func debug(s string) {
	if DEBUG {
		fmt.Println("DEBUG: " + s)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(FATAL_ERROR, err.Error()))
		os.Exit(1)
	}
}

func CLIclearLine() {
	x := ""
	for ; len(x) < 10000; x += "\b" {
	}
	fmt.Print(x)
}

// Receive and process data from in
func receiveData(in *bufio.Reader, qtd int) ([]byte, error) {
	data := make([]byte, qtd) // in buffer
	n, err := in.Read(data)   // Reading from network
	return data[:n], err
}

//// Prepare and send data to out
//func sendData(out *bufio.Writer, data []byte) (error) {
//	_, err := out.Write(data)
//	out.Flush()
//	return err
//}
