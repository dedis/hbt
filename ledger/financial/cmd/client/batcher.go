package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"ledger/financial/handler"
	"ledger/financial/payload"
	"ledger/financial/state"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"google.golang.org/protobuf/proto"

	tx_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
)

// batcher defines the primitive to send transactions
type batcher interface {
	sendTransaction(payload payload.FinancialPayload) (string, error)
}

// newBatcher returns a new initialized batcher. url must be the REST endpoint,
// and key must be the filepath to a private key that is allowed to send
// batches.
func newBatcher(url url.URL, keyfilePath string) (batcher, error) {
	privateKey, err := loadKey(keyfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key: %v", err)
	}

	pk := signing.NewSecp256k1Context().GetPublicKey(privateKey)
	fmt.Println("pk:", pk.AsHex())

	cryptoFac := signing.NewCryptoFactory(signing.NewSecp256k1Context())
	signer := *cryptoFac.NewSigner(privateKey)

	return basicBatcher{
		url:    url,
		signer: signer,
	}, nil
}

// batcher holds the functionality of creating and sending batches
type basicBatcher struct {
	url    url.URL
	signer signing.Signer
}

// sendTransaction implements batcher
func (b basicBatcher) sendTransaction(payload payload.FinancialPayload) (string, error) {

	payloadBuff, err := payload.ToBytes()
	if err != nil {
		return "", fmt.Errorf("failed to transform payload into bytes: %v", err)
	}

	hash := sha512.New()
	hash.Write(payloadBuff)

	payloadSha512 := strings.ToLower(hex.EncodeToString(hash.Sum(nil)))

	address := state.KeyAddress

	nonce := make([]byte, 16)

	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to get random nonce: %v", err)
	}

	rawTxHeader := tx_pb2.TransactionHeader{
		SignerPublicKey:  b.signer.GetPublicKey().AsHex(),
		FamilyName:       handler.FamilyName,
		FamilyVersion:    handler.FamilyVersion,
		Dependencies:     []string{},
		Nonce:            base64.RawStdEncoding.EncodeToString(nonce),
		BatcherPublicKey: b.signer.GetPublicKey().AsHex(),
		Inputs:           []string{address},
		Outputs:          []string{address},
		PayloadSha512:    payloadSha512,
	}

	txHeader, err := proto.Marshal(&rawTxHeader)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction header: %v", err)
	}

	txHeaderSig := hex.EncodeToString(b.signer.Sign(txHeader))

	transaction := tx_pb2.Transaction{
		Header:          txHeader,
		HeaderSignature: txHeaderSig,
		Payload:         payloadBuff,
	}

	rawBatchList, err := b.createBatch([]*tx_pb2.Transaction{&transaction})
	if err != nil {
		return "", fmt.Errorf("failed to create batch: %v", err)
	}

	// batchID := rawBatchList.Batches[0].HeaderSignature

	batchList, err := proto.Marshal(&rawBatchList)
	if err != nil {
		return "", fmt.Errorf("failed to proto marshal batch list: %v", err)
	}

	response, err := b.sendRequest(batchList, BatchSubmitAPI, ContentTypeOctetStream)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	return response, nil
}

func (b basicBatcher) sendRequest(data []byte, apiSuffix, contentType string) (string, error) {

	url := b.url.String() + "/" + apiSuffix

	logger.Infof("sending request to %q", url)

	var resp *http.Response
	var err error

	if len(data) > 0 {
		resp, err = http.Post(url, contentType, bytes.NewBuffer(data))
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return "", fmt.Errorf("failed to send http query: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		buff, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %q: %s", resp.Status, buff)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %v", err)
	}

	return string(body), nil
}

func (b basicBatcher) createBatch(txs []*tx_pb2.Transaction) (batch_pb2.BatchList, error) {

	txSignatures := make([]string, len(txs))
	for i, tx := range txs {
		txSignatures[i] = tx.GetHeaderSignature()
	}

	rawBatchHeader := batch_pb2.BatchHeader{
		SignerPublicKey: b.signer.GetPublicKey().AsHex(),
		TransactionIds:  txSignatures,
	}

	batchHeader, err := proto.Marshal(&rawBatchHeader)
	if err != nil {
		return batch_pb2.BatchList{}, fmt.Errorf("failed to proto marshal batch header: %v", err)
	}

	batchHeaderSignature := hex.EncodeToString(b.signer.Sign(batchHeader))

	batch := batch_pb2.Batch{
		Header:          batchHeader,
		HeaderSignature: batchHeaderSignature,
		Transactions:    txs,
	}

	return batch_pb2.BatchList{
		Batches: []*batch_pb2.Batch{&batch},
	}, nil
}

// loadKey loads an existing private or creates one
func loadKey(keyfilePath string) (signing.PrivateKey, error) {
	var privateKey signing.PrivateKey

	_, err := os.Stat(keyfilePath)
	if err != nil && os.IsExist(err) {
		return nil, fmt.Errorf("failed to check the private key: %v", err)
	}

	if os.IsNotExist(err) {
		privateKey = signing.NewSecp256k1Context().NewRandomPrivateKey()

		err = os.MkdirAll(filepath.Dir(keyfilePath), 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to create folder: %v", err)
		}

		err = os.WriteFile(keyfilePath, []byte(privateKey.AsHex()), 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to write key file: %v", err)
		}
	} else {
		privateKeyHex, err := os.ReadFile(keyfilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %v", err)
		}

		privateKeyBuff, err := hex.DecodeString(string(bytes.TrimSpace(privateKeyHex)))
		if err != nil {
			return nil, fmt.Errorf("failed to decode key hex: %v", err)
		}

		privateKey = signing.NewSecp256k1PrivateKey(privateKeyBuff)
	}

	return privateKey, nil
}
