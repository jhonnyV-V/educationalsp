package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"educationalsp/rpc"
)

func main() {
	logfile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger := getLogger(logfile)
	logger.Println("starting")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}

func getLogger(logfile io.Writer) *log.Logger {
	return log.New(logfile, "[educationalsp]", log.Ldate|log.Lshortfile|log.Ltime)
}
