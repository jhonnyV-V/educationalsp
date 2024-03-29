package lsp

type PublishDiagnosticsNotification struct {
	Notification
	Params PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range Range `json:"range"`
	//1 = error 2 = warning 3 = information 4 = hint
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}
