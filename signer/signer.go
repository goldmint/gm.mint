package signer

import (
	"fmt"

	"github.com/void616/gm-sumus-lib/signer/ed25519"
)

// NewSigner made from random keypair
func NewSigner() (*Signer, error) {

	// generate pk
	_, pk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	// pk contains seed+public - prehash it
	pkPrehashed := pk.Prehash()

	return &Signer{
		privateKey:  pkPrehashed,
		publicKey:   ed25519.PublicKeyFromPrehashedPK(pkPrehashed),
		initialized: true,
	}, nil
}

// NewSignerFromPK makes keypair from prehashed seed/PK
func NewSignerFromPK(b []byte) (*Signer, error) {

	// check pk size
	if len(b) != 64 {
		return nil, fmt.Errorf("Private key has invalid size %v, expected %v", len(b), 64)
	}

	// copy pk bytes
	pvt := make([]byte, 64)
	copy(pvt, b)

	ret := &Signer{
		privateKey:  pvt,
		publicKey:   ed25519.PublicKeyFromPrehashedPK(pvt),
		initialized: true,
	}

	return ret, nil
}

// ---

// Signer data
type Signer struct {
	initialized bool
	privateKey  []byte
	publicKey   []byte
}

func (s *Signer) assert() {
	if !s.initialized {
		panic("Signer is not initialized")
	}
}

// Sign message with a key
func (s *Signer) Sign(message []byte) []byte {
	s.assert()
	return ed25519.SignWithPrehashed(s.privateKey, s.publicKey, message)
}

// PrivateKey of the signer
func (s *Signer) PrivateKey() []byte {
	s.assert()
	var ret [64]byte
	copy(ret[:], s.privateKey)
	return ret[:]
}

// PublicKey of the signer
func (s *Signer) PublicKey() []byte {
	s.assert()
	var ret [32]byte
	copy(ret[:], s.publicKey)
	return ret[:]
}
