package main

import "flag"

var serverPort = flag.Uint("-p", 8080, "Server port")
var ch = make(chan OutlierDetectOutput)

func main() {
	go OutliersReporter(ch)
	go DataSetsChecker(ch)
	StartServer(*serverPort)
}
