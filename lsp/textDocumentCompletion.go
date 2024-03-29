package lsp

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	Position Position `json:"position"`
}

type CompletionResponse struct {
	Response
	Result CompletionResult `json:"result"`
}

type CompletionResult struct {
	Items []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}
