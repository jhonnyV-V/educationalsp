package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"educationalsp/lsp"
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
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s\n", err)
			continue
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
	logger.Printf("receive message with method %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Failed to parse: %s\n", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		writer := os.Stdout
		msg := lsp.NewInitializedResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		if _, err := writer.Write([]byte(reply)); err != nil {
			logger.Printf("Failed to write: %s\n", err)
		}
		logger.Printf("send response %s", reply)
	}
}

func getLogger(logfile io.Writer) *log.Logger {
	return log.New(logfile, "[educationalsp]", log.Ldate|log.Lshortfile|log.Ltime)
}
