package transaction

import (
	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/serializer"
	"github.com/void616/gm-sumuslib/signer"
	"golang.org/x/crypto/sha3"
)

// construct constructs and signs transaction
type constructor struct {
	*serializer.Serializer
}

// SignedTransaction data
type SignedTransaction struct {
	Digest    sumuslib.Digest
	Data      []byte
	Signature sumuslib.Signature
}

func newConstructor(nonce uint64) *constructor {
	c := &constructor{
		Serializer: serializer.NewSerializer(),
	}

	// write nonce/ID
	c.PutUint64(nonce)

	return c
}

// Sign signs transaction data
func (c *constructor) Sign(signer *signer.Signer) (*SignedTransaction, error) {
	// get payload
	payload, err := c.Data()
	if err != nil {
		return nil, err
	}

	// make payload digest
	var txdigest sumuslib.Digest
	{
		hasher := sha3.New256()
		_, err = hasher.Write(payload)
		if err != nil {
			return nil, err
		}
		digest := hasher.Sum(nil)
		copy(txdigest[:], digest)
	}

	// sign digest
	txsignature := signer.Sign(txdigest[:])

	// signature
	c.PutByte(1)               // append a byte - "signed bit"
	c.PutBytes(txsignature[:]) // signature

	// data
	txdata, err := c.Data()
	if err != nil {
		return nil, err
	}

	return &SignedTransaction{
		Data:      txdata,
		Digest:    txdigest,
		Signature: txsignature,
	}, nil
}
