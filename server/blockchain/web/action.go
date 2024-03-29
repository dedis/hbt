package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/mino/proxy"
	"go.dedis.ch/kyber/v3/suites"

	"golang.org/x/xerrors"
)

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

const separator = ":"
const malformedEncoded = "malformed encoded: %s"

// RegisterAction is an action to register the HTTP handlers
//
// - implements node.ActionTemplate
type RegisterAction struct{}

// Execute implements node.ActionTemplate. It registers the handlers using the
// default proxy from the the injector.
func (a *RegisterAction) Execute(ctx node.Context) error {
	var p proxy.Proxy
	err := ctx.Injector.Resolve(&p)
	if err != nil {
		return xerrors.Errorf("failed to resolve proxy: %v", err)
	}

	router := mux.NewRouter()

	s := &secretHandler{&ctx}
	router.HandleFunc("/secret", s.addSecret).Methods("POST")
	router.HandleFunc("/secret/list", s.listSecrets).Methods("GET")
	router.HandleFunc("/secret", s.getSecret).Methods("GET")
	router.HandleFunc("/secret/reveal", s.revealSecret).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowedHandler)

	p.RegisterHandler("/secret/", router.ServeHTTP)

	dela.Logger.Info().Msg("proxy handlers registered")

	return nil
}

type DocID []byte

// a secretData is a struct to hold the secret data: the document ID and the
// encrypted key to access the document
type secretData struct {
	docID        DocID  `json:"docid"`
	encryptedKey string `json:"secret"`
}

type secretHandler struct {
	ctx *node.Context
}

// addSecret adds a new secret in the blockchain
func (s *secretHandler) addSecret(w http.ResponseWriter, r *http.Request) {
	ctx := *(s.ctx)

	// Decode the request
	var sec secretData
	err := json.NewDecoder(r.Body).Decode(&sec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add the secret to the blockchain
	// the secret is added to the blockchain with the document ID as the key
	// and the encrypted key as the value
	// TODO add it to the blockchain
	dela.Logger.Info().Msgf("secret added to the blockchain: ID=%v secret=%v", sec.docID,
		sec.encryptedKey)
}

// listSecrets lists all secrets in the blockchain
func (s *secretHandler) listSecrets(w http.ResponseWriter, r *http.Request) {
	ctx := *(s.ctx)

	// list all secrets from the blockchain

}

// getSecret gets a secret from the blockchain
func (s *secretHandler) getSecret(w http.ResponseWriter, r *http.Request) {
	ctx := *(s.ctx)

	// Decode the request
	var id DocID
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the secret from the blockchain
	// TODO

}

// -----------------------------------------------------------------------------
// Helper functions

// HTTPError defines the standard error format
type HTTPError struct {
	Title   string
	Code    uint
	Message string
	Args    map[string]interface{}
}

// notFoundHandler defines a generic handler for 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := HTTPError{
		Title:   "Not found",
		Code:    http.StatusNotFound,
		Message: "The requested endpoint was not found",
		Args: map[string]interface{}{
			"url":    r.URL.String(),
			"method": r.Method,
		},
	}

	buf, _ := json.MarshalIndent(&err, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, string(buf))
}

// notAllowedHandler degines a generic handler for 405
func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := HTTPError{
		Title:   "Not allowed",
		Code:    http.StatusMethodNotAllowed,
		Message: "The requested endpoint was not allowed",
		Args: map[string]interface{}{
			"url":    r.URL.String(),
			"method": r.Method,
		},
	}

	buf, _ := json.MarshalIndent(&err, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintln(w, string(buf))
}
