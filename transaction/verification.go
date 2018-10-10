package transaction

import (
	"fmt"

	"github.com/void616/gm-sumus-lib/signer"
	"golang.org/x/crypto/sha3"
)

// Verify transaction payload
func Verify(address, payload, signature []byte) error {
	if address == nil || payload == nil || signature == nil {
		return fmt.Errorf("Null argument specified")
	}

	// make payload digest
	hasher := sha3.New256()
	_, err := hasher.Write(payload)
	if err != nil {
		return err
	}
	digest := hasher.Sum(nil)

	// verify
	return signer.Verify(address, digest, signature)
}
