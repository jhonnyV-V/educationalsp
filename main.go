package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"educationalsp/analysis"
	"educationalsp/lsp"
	"educationalsp/rpc"
)

func main() {
	logfile, err := os.OpenFile("/home/jv/code/go/educationalsp/log.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger := getLogger(logfile)
	logger.Println("starting")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s\n", err)
			continue
		}
		writer := os.Stdout
		handleMessage(logger, writer, &state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *analysis.State, method string, content []byte) {
	logger.Printf("receive message with method %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Failed to parse: %s\n", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		msg := lsp.NewInitializedResponse(request.ID)
		writeResponse(writer, msg, logger)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s\n", err)
		}
		logger.Printf("textDocument/didOpen: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		msg := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				Method: "textDocument/publishDiagnostics",
				RPC:    "2.0",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		}
		writeResponse(writer, msg, logger)
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didChange: %s\n", err)
		}
		logger.Printf("Change: %s", request.Params.TextDocument.TextDocumentIdentifier)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			msg := lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					Method: "textDocument/publishDiagnostics",
					RPC:    "2.0",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			}
			writeResponse(writer, msg, logger)
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/hover: %s\n", err)
		}
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response, logger)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/definition: %s\n", err)
		}
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response, logger)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s\n", err)
		}
		logger.Printf("textDocument/codeAction: %s", request.Params.TextDocument.URI)

		response := state.CodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response, logger)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/completion: %s\n", err)
		}
		logger.Printf("textDocument/completion: %s", request.Params.TextDocument.URI)

		response := state.Completion(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response, logger)
	}
}

func getLogger(logfile io.Writer) *log.Logger {
	return log.New(logfile, "[educationalsp]", log.Ldate|log.Lshortfile|log.Ltime)
}

func writeResponse(writer io.Writer, msg any, logger *log.Logger) {
	reply := rpc.EncodeMessage(msg)
	if _, err := writer.Write([]byte(reply)); err != nil {
		log.Printf("Failed to write: %s\n", err)
	}

	logger.Printf("send response %s", reply)
}
