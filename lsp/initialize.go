package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	//there is a lot more params
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
