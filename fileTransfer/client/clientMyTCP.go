package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strconv"
	"myTCP/myTCP"
)

// Wait for server message and show to client
func ListenServer(in *bufio.Reader) {
	for {
		data, err := receiveData(in, 4096)
		checkError(err)

		fmt.Printf(string(data))
	}
	in = nil
}

func main() {
	host, port, filename := "127.0.0.1", "10101", "test.txt"
	//flag.StringVar(&host, "h", "", "Give a hostname with -h=#####")
	//flag.StringVar(&port, "p", "", "Give a port with -p=#####")
	//flag.StringVar(&filename, "f", "", "Give a filename with -f=#####")
	//flag.Parse()
	//
	//if port == "" || host == "" || filename == "" {
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}

	serverAddr, err := myTCP.ResolveName(host + ":" + port)
	//serverAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	checkError(err)

	//// for TCP version
	//localAddr, err := net.ResolveTCPAddr("tcp", host+":0")
	//checkError(err)

	conn, err := myTCP.Connect(serverAddr)
	//conn, err := net.DialTCP("tcp", localAddr, serverAddr)
	checkError(err)
	defer conn.Close()

	//in := bufio.NewReader(conn)
	//out := bufio.NewWriter(conn)

	debug("Connected with server at " + host + " " + port)

	//debug("Waiting for some server message")
	//go ListenServer(in)

	debug("Opening file " + filename)
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()

	debug("Send file to server")
	n, err := io.Copy(conn, file)
	checkError(err)

	debug("Sent " + strconv.Itoa(int(n)) + " bytes to server")
}
