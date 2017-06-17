package myTCP

import "strings"

func Listen(ipPort string) {
	ipPortSplit := strings.Split(ipPort, ":")
	ip, port := ipPortSplit[0], ipPortSplit[1]

	Printa("IP: " + ip + ", PORT: " + port)
}
