package key

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

// NewAsymmetric generates a new kyber V3 asymmetric key pair
func NewAsymmetric() (kyber.Point, kyber.Scalar) {
	suite := suites.MustFind("Ed25519")

	// Create a public/private keypair
	sk := suite.Scalar().Pick(suite.RandomStream())
	pk := suite.Point().Mul(sk, nil)

	return pk, sk
}
