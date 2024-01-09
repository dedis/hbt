package main

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func IgnoreTestCreateDocument(t *testing.T) {
	apitest.New().
		Handler(newApp().UserRouter).
		Post("/document").
		Header("Content-Type", "application/json").
		Body(`{"name":"John Doe","passport":"12XY3456789","role":0,"picture":"1234567890987654321","registered":false}`).
		Expect(t).
		Body(`{"document_id": 1}`).
		Status(http.StatusOK).
		End()
}
