package sproxy

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/dkg"
	"go.dedis.ch/dela/mino/proxy"
	eproxy "go.dedis.ch/hbt/server/smc/proxy"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"

	"golang.org/x/xerrors"
)

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

const separator = ":"
const malformedEncoded = "malformed encoded: %s"

/*
const (
	smcPath            = "/smc"
	smcPubKeyPath      = "/smc/pubkey"
	smcReencryptPath   = "/smc/reencrypt"
	resolveActorFailed = "failed to resolve actor, did you call listen?: %v"
)
*/

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

	pk := &pubKeyHandler{&ctx}
	router.HandleFunc("/smc/pubkey", pk.ServeHTTP).Methods("GET")

	re := &reencryptHandler{&ctx}
	router.HandleFunc("/smc/reencrypt", re.ServeHTTP).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(eproxy.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(eproxy.NotAllowedHandler)

	p.RegisterHandler("/smc/", router.ServeHTTP)

	dela.Logger.Info().Msg("proxy handlers registered")

	return nil
}

type pubKeyHandler struct {
	ctx *node.Context
}

func (h *pubKeyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {

	c := *(h.ctx)

	var actor dkg.Actor
	err := c.Injector.Resolve(&actor)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resolve DKG actor: %v", err),
			http.StatusInternalServerError)
		return
	}

	pk, err := actor.GetPublicKey()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed retrieving public key: %v", err),
			http.StatusInternalServerError)
		return
	}

	b, err := pk.MarshalBinary()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal public key: %v", err),
			http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(b)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to respond: %v", err),
			http.StatusInternalServerError)
		return
	}
}

type reencryptHandler struct {
	ctx *node.Context
}

func (h *reencryptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := *(h.ctx)

	var actor dkg.Actor
	err := c.Injector.Resolve(&actor)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resolve DKG actor: %v", err),
			http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Err(err).Msg("failed to parse form")
		http.Error(w, "failed to parse form", http.StatusInternalServerError)
		return
	}

	// XHATENC=$(smccli --config /tmp/smc1 dkg reencrypt --encrypted ${CIPHER} --pubk ${PUBK})

	// retrieve the public key
	pubkString := r.FormValue("pubk")
	pubk, err := decodePublicKey(pubkString)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode public key str: %v", err),
			http.StatusInternalServerError)
		return
	}

	// retriev the encrypted cypher
	encrypted := r.FormValue("encrypted")
	k, _, err := decodeEncrypted(encrypted)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode encrypted str: %v", err),
			http.StatusInternalServerError)
		return
	}

	// re-encrypt the message
	hatenc, err := actor.Reencrypt(k, pubk)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to re-encrypt: %v", err),
			http.StatusInternalServerError)
		return
	}

	// write back the re-encrypted message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(hatenc)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err),
			http.StatusInternalServerError)
		return
	}

	dela.Logger.Debug().Msgf("Re-encrypted message: %v", hatenc)
}

// -----------------------------------------------------------------------------
// Helper functions
func decodePublicKey(str string) (kyber.Point, error) {
	pkbuff, err := hex.DecodeString(str)
	if err != nil {
		return nil, xerrors.Errorf(malformedEncoded, str)
	}

	pk := suite.Point()
	err = pk.UnmarshalBinary(pkbuff)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal pk: %v", err)
	}

	return pk, nil
}

func decodeEncrypted(str string) (kyber.Point, []kyber.Point, error) {
	parts := strings.Split(str, separator)
	if len(parts) < 2 {
		return nil, nil, xerrors.Errorf(malformedEncoded, str)
	}

	// Decode K
	kbuff, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to decode k point: %v", err)
	}

	k := suite.Point()
	err = k.UnmarshalBinary(kbuff)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to unmarshal k point: %v", err)
	}

	// Decode Cs
	cs := make([]kyber.Point, 0, len(parts)-1)
	for _, p := range parts[1:] {
		cbuff, err := hex.DecodeString(p)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to decode c point: %v", err)
		}

		c := suite.Point()
		err = c.UnmarshalBinary(cbuff)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to unmarshal c point: %v", err)
		}

		cs = append(cs, c)
	}

	dela.Logger.Debug().Msgf("Decoded K: %v and Cs: %v", k, cs)

	return k, cs, nil
}
