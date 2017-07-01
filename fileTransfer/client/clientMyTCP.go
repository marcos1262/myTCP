package main

import (
	"myTCP/myTCP"
)

func main() {
	serverAddr, err := myTCP.ResolveName("127.0.0.1:10001")
	checkError(err)

	conn, err := myTCP.Connect(serverAddr)
	checkError(err)
	defer conn.Close()

	//i := 0
	//for {
	//	msg := strconv.Itoa(i)
	//	buf := []byte(msg)
	//	i++
	//	_, err := conn.Write(buf)
	//	if err != nil {
	//		fmt.Println(msg, err)
	//	}
	//	time.Sleep(time.Second * 1)
	//}
}
