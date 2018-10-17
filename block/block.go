package block

import (
	"fmt"
	"io"
	"math/big"

	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/types"
)

// Header data
type Header struct {

	// Version of the blockchain
	Version uint16
	// PrevBlockDigest, 32 bytes
	PrevBlockDigest []byte
	// MerkleRoot, 32 bytes
	MerkleRoot []byte
	// Timestamp of the block
	Timestamp uint64
	// TransactionsCount in the block
	TransactionsCount uint16
	// BlockNumber, 32 bytes
	BlockNumber *big.Int
	// SignersCount
	SignersCount uint16
	// Signers list
	Signers []Signer
}

// Signer data
type Signer struct {

	// PublicKey, 32 bytes
	PublicKey []byte
	// Signature, 64 bytes
	Signature []byte
}

// CbkHeader for parsed header
type CbkHeader func(*Header) error

// CbkTransaction for parsed transaction
type CbkTransaction func(types.Transaction, *serializer.Deserializer, *Header) error

// ---

// Parse block
func Parse(r io.Reader, cbkHeader CbkHeader, cbkTransaction CbkTransaction) error {

	d := serializer.NewStreamDeserializer(r)

	// read header
	header := &Header{}
	header.Version = d.GetUint16()           // version
	header.PrevBlockDigest = d.GetBytes(32)  // previous block digest
	header.MerkleRoot = d.GetBytes(32)       // merkle root
	header.Timestamp = d.GetUint64()         // time
	header.TransactionsCount = d.GetUint16() // transactions
	header.BlockNumber = d.GetUint256()      // block
	header.SignersCount = d.GetUint16()      // signers
	if err := d.Error(); err != nil {
		return err
	}

	// read signers list
	header.Signers = make([]Signer, header.SignersCount)
	for i := uint16(0); i < header.SignersCount; i++ {

		sig := Signer{}
		sig.PublicKey = d.GetBytes(32) // address
		sig.Signature = d.GetBytes(64) // signature

		if err := d.Error(); err != nil {
			return err
		}
		header.Signers[i] = sig
	}

	// callback
	if err := cbkHeader(header); err != nil {
		return err
	}

	// read transactions
	for i := uint16(0); i < header.TransactionsCount; i++ {

		txCode := d.GetUint16() // code
		if err := d.Error(); err != nil {
			return err
		}

		// check the code
		if !types.ValidTransaction(txCode) {
			return fmt.Errorf("Unknown transaction with code `%v` with index %v", txCode, i)
		}
		txType := types.Transaction(txCode)

		// parse transaction outside
		if err := cbkTransaction(txType, d, header); err != nil {
			return err
		}
		if err := d.Error(); err != nil {
			return err
		}
	}

	return nil
}
