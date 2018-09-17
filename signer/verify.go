package signer

import (
	"fmt"

	"github.com/void616/gm-sumus-lib/signer/ed25519"
)

// Verify a message with a public key of a signer
func Verify(pub, message, sig []byte) error {

	// check public key size
	if len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("Public key has invalid size %v, expected %v", len(pub), ed25519.PublicKeySize)
	}

	// check message size
	if len(message) == 0 {
		return fmt.Errorf("Message has invalid size %v, expected at least 1", len(message))
	}

	// check signature size
	if len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("Signature has invalid size %v, expected %v", len(sig), ed25519.SignatureSize)
	}

	// verify
	if !ed25519.Verify(pub, message, sig) {
		return fmt.Errorf("Invalid signature for this message")
	}

	return nil
}
