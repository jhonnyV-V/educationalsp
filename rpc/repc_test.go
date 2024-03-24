package rpc_test

import (
	"educationalsp/rpc"
	"testing"
)

type EcodingExample struct {
	Method string `json:"method"`
}

func TestEncoding(t *testing.T) {
	expected := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	actual := rpc.EncodeMessage(EcodingExample{"hi"})
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestDecoding(t *testing.T) {
	expectedByteLength := 15
	method, content, err := rpc.DecodeMessage([]byte("Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"))
	contentLength := len(content)
	if err != nil {
		t.Fatal(err)
	}
	if contentLength != expectedByteLength {
		t.Fatalf("expected contentLength %d, got %d", expectedByteLength, contentLength)
	}

	if method != "hi" {
		t.Fatalf("expected method hi, got %s", method)
	}
}
