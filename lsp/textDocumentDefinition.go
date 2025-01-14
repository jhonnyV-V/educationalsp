package lsp

type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
	Position Position `json:"position"`
}

type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
