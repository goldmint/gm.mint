package transaction

import (
	"bytes"
	"fmt"

	sumus "github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/amount"
	"github.com/void616/gm-sumus-lib/serializer"
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

// VerifyAssetTransaction with payload check
func VerifyAssetTransaction(tx, sourceAddr, expDestAddr []byte, expNonce *uint64, expToken *sumus.Token, expTokenAmount *amount.Amount) error {

	if tx == nil || len(tx) == 0 {
		return fmt.Errorf("Transaction bytes array is null or empty")
	}

	// get payload and signature
	des := serializer.NewDeserializer(tx)
	txPayload := des.GetBytes(89)
	tSigned := des.GetByte()
	tSignature := des.GetBytes(64)
	err := des.Error()
	if err != nil {
		return fmt.Errorf("Failed to get payload and signature")
	}

	// is signed?
	if tSigned != 1 {
		return fmt.Errorf("Is not signed")
	}

	// verify signature
	err = Verify(sourceAddr, txPayload, tSignature)
	if err != nil {
		return fmt.Errorf("Failed to verify signature")
	}

	// read payload
	desPayload := serializer.NewDeserializer(txPayload)
	tNonce := desPayload.GetUint64()
	tToken := desPayload.GetUint16()
	tSource := desPayload.GetBytes(32)
	tDestination := desPayload.GetBytes(32)
	tTokenAmount := desPayload.GetAmount()
	err = desPayload.Error()
	if err != nil {
		return fmt.Errorf("Failed to read payload")
	}

	if expNonce != nil && tNonce != *expNonce {
		return fmt.Errorf("Invalid nonce")
	}
	if expToken != nil && tToken != uint16(*expToken) {
		return fmt.Errorf("Invalid token")
	}
	if !bytes.Equal(tSource, sourceAddr) {
		return fmt.Errorf("Invalid source address")
	}
	if expDestAddr != nil && !bytes.Equal(tDestination, expDestAddr) {
		return fmt.Errorf("Invalid destination address")
	}
	if tTokenAmount == nil || tTokenAmount.Value == nil || (expTokenAmount != nil && tTokenAmount.Value.Cmp(expTokenAmount.Value) != 0) {
		return fmt.Errorf("Invalid token amount")
	}

	return nil
}
