package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/core/execution"
	"go.dedis.ch/dela/core/txn"
	"go.dedis.ch/dela/core/txn/signed"
	"go.dedis.ch/dela/crypto/bls"
	"go.dedis.ch/dela/mino/proxy"
	"go.dedis.ch/hbt/server/blockchain/calypso"
	purbkv "go.dedis.ch/purb-db/store/kv"

	"golang.org/x/xerrors"
)

// RegisterAction is an action to register the HTTP handlers
//
// - implements node.ActionTemplate
type RegisterAction struct{}

// Execute implements node.ActionTemplate. It registers the handlers using the
// default proxy from the injector.
func (a *RegisterAction) Execute(ctx node.Context) error {
	var p proxy.Proxy
	err := ctx.Injector.Resolve(&p)
	if err != nil {
		return xerrors.Errorf("failed to resolve proxy: %v", err)
	}

	router := mux.NewRouter()

	s := &secretHandler{ctx}
	router.HandleFunc("/secret/smc", s.advertiseSmc).Methods("POST")

	router.HandleFunc("/secret", s.addSecret).Methods("POST")

	router.HandleFunc("/secret/admin/list", s.listSecrets).Methods("GET")
	router.HandleFunc("/secret/admin", s.getSecret).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowedHandler)

	p.RegisterHandler("/secret/", router.ServeHTTP)

	dela.Logger.Info().Msg("proxy handlers registered")

	return nil
}

type DocID []byte

type secretHandler struct {
	ctx node.Context
}

// advertiseSmc advertises the SMC public key and its roster to the blockchain
func (s *secretHandler) advertiseSmc(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal().Err(err)
	}

	smckey := r.FormValue("smckey")
	roster := r.FormValue("roster")
	dela.Logger.Info().Msgf("received SMC pubkey %v from SMC roster %v", smckey, roster)

	// get the calypso contract
	var c calypso.Contract
	err = s.ctx.Injector.Resolve(&c)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve calypso contract")
		http.Error(w, fmt.Sprintf("failed to resolve calypso contract: %v", err),
			http.StatusInternalServerError)
		return
	}

	var db purbkv.DB
	err = s.ctx.Injector.Resolve(&db)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve PURB database")
		http.Error(w, fmt.Sprintf("failed to resolve database: %v", err),
			http.StatusInternalServerError)
		return
	}

	err = db.Update(func(txn purbkv.WritableTx) error {
		b, err := txn.GetBucketOrCreate([]byte("bucket:secret"))
		if err != nil {
			return err
		}

		err = c.Execute(b, makeStep(calypso.CmdArg, string(calypso.CmdAdvertiseSmc),
			calypso.SmcPublicKeyArg, smckey, calypso.RosterArg, roster))

		return err
	})

	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to advertise SMC to the blockchain")
		http.Error(w, fmt.Sprintf("failed to advertise SMC to the blockchain: %v", err),
			http.StatusInternalServerError)
		return
	}

	dela.Logger.Info().Msg("SMC advertised to the blockchain")

	return
}

// addSecret adds a new secret in the blockchain
func (s *secretHandler) addSecret(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal().Err(err)
	}

	smckey := r.FormValue("smckey")
	secret := r.FormValue("secret")
	id := r.FormValue("id")
	dela.Logger.Info().Msgf("received doc ID=%v with secret=%v", id, secret)

	// get the calypso contract
	var c calypso.Contract
	err = s.ctx.Injector.Resolve(&c)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve calypso contract")
		http.Error(w, fmt.Sprintf("failed to resolve calypso contract: %v", err),
			http.StatusInternalServerError)
		return
	}

	var db purbkv.DB
	err = s.ctx.Injector.Resolve(&db)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve database")
		http.Error(w, fmt.Sprintf("failed to resolve PURB database: %v", err),
			http.StatusInternalServerError)
		return
	}

	// add the secret to the blockchain
	// the secret is added to the blockchain with the document ID as the key
	// and the encrypted key as the value

	err = db.Update(func(txn purbkv.WritableTx) error {
		b, err := txn.GetBucketOrCreate([]byte("bucket:secret"))
		if err != nil {
			return err
		}

		err = c.Execute(b, makeStep(calypso.CmdArg, string(calypso.CmdCreateSecret),
			calypso.SmcPublicKeyArg, smckey,
			calypso.SecretNameArg, id, calypso.SecretArg, secret))

		return err
	})

	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to add secret to the blockchain")
		http.Error(w, fmt.Sprintf("failed to add secret to the blockchain: %v", err),
			http.StatusInternalServerError)
		return
	}

	dela.Logger.Info().Msgf("secret added to the blockchain: ID=%v secret=%v", id, secret)
}

// listSecrets lists all secrets in the blockchain
func (s *secretHandler) listSecrets(w http.ResponseWriter, r *http.Request) {
	// list all secrets from the blockchain
	r.ParseForm()

	pubkey := r.Form.Get("pubkey")
	dela.Logger.Info().Msgf("received request from %v to list the secrets", pubkey)

	// get the calypso contract
	var c calypso.Contract
	err := s.ctx.Injector.Resolve(&c)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve calypso contract")
		http.Error(w, fmt.Sprintf("failed to resolve calypso contract: %v", err),
			http.StatusInternalServerError)
		return
	}

	var db purbkv.DB
	err = s.ctx.Injector.Resolve(&db)
	if err != nil {
		dela.Logger.Error().Err(err).Msg("failed to resolve database")
		http.Error(w, fmt.Sprintf("failed to resolve database: %v", err),
			http.StatusInternalServerError)
		return
	}

	err = db.View(func(txn purbkv.ReadableTx) error {
		b := txn.GetBucket([]byte("bucket:secret"))

		err = c.Execute(b, makeStep(calypso.CmdArg, string(calypso.CmdListSecrets)))

		return err
	})

}

// getSecret gets a secret from the blockchain
func (s *secretHandler) getSecret(w http.ResponseWriter, r *http.Request) {
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

// -----------------------------------------------------------------------------
// Utility functions
var nounce uint64

func makeStep(args ...string) execution.Step {
	return execution.Step{Current: makeTx(args...)}
}

func makeTx(args ...string) txn.Transaction {
	signer := bls.NewSigner()

	options := []signed.TransactionOption{}
	for i := 0; i < len(args)-1; i += 2 {
		options = append(options, signed.WithArg(args[i], []byte(args[i+1])))
	}

	tx, err := signed.NewTransaction(nounce, signer.GetPublicKey(), options...)
	if err != nil {
		nounce++
	}

	return tx
}
