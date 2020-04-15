package transaction

import (
	"bytes"
	"testing"

	mint "github.com/void616/gm.mint"
	"github.com/void616/gm.mint/amount"
	"github.com/void616/gm.mint/serializer"
	"github.com/void616/gm.mint/signer"
)

func TestTransferAssetTransactionVerification(t *testing.T) {

	nonce := uint64(2)
	token := mint.TokenMNT
	tokenAmount := amount.MustFromString("123.666")

	// ---

	srcpk, _ := mint.ParsePrivateKey("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	src := signer.FromPrivateKey(srcpk)

	dstpk, _ := mint.ParsePrivateKey("FhM2u3UMtexZ3TU57G6d9iDpcmynBSpzmTZq6YaMPeA6DHFdEht3jcZUDpXyVbXGoXoWiYB9z8QVKjGhZuKCqMGYZE2P6")
	dst := signer.FromPrivateKey(dstpk)

	tx, err := (&TransferAsset{
		Address: dst.PublicKey(),
		Token:   token,
		Amount:  tokenAmount,
	}).Sign(src, nonce)
	if err != nil {
		t.Fatal(err)
	}

	txBytes := tx.Data

	// ---

	// tx:
	// 0300000000000000 - nonce, 8B --------------------------|
	// 0000 - token, 2B                                       |
	// eea072...01d - pub key from, 32B                       |--- payload, 90 bytes
	// f42378...7b7 - pub key to, 32B                         |
	// 000000000000000000000010000000 - amount, 16B ----------|
	// 01 - signed byte, 1B
	// 0e4b04...90a - signature, 64B

	// get payload and signature
	des := serializer.NewDeserializer(txBytes)
	txPayload := des.GetBytes(90)
	tSigned := des.GetByte()
	tSignatureBytes := des.GetBytes(64)
	err = des.Error()
	if err != nil {
		t.Fatal(err, "Failed to get payload and signature")
	}

	// is signed?
	if tSigned != 1 {
		t.Fatal("Is not signed")
	}

	var tSignature mint.Signature
	copy(tSignature[:], tSignatureBytes)

	// verify if signed
	err = Verify(src.PublicKey(), txPayload, tSignature)
	if err != nil {
		t.Fatal(err, "Failed to verify signature")
	}

	// and this should fail
	err = Verify(src.PublicKey(), txPayload[:len(txPayload)-1], tSignature)
	if err == nil {
		t.Fatal("Invalid signature is valid")
	}

	// read payload
	desPayload := serializer.NewDeserializer(txPayload)
	tNonce := desPayload.GetUint64()
	tToken := desPayload.GetUint16()
	tSource := desPayload.GetBytes(32)
	tDestination := desPayload.GetBytes(32)
	tTokenAmount := desPayload.GetAmount()
	if desPayload.Error() != nil {
		t.Fatal("Failed to read payload")
	}

	srcPub := src.PublicKey()
	dstPub := dst.PublicKey()

	if tNonce != nonce {
		t.Fatal("Invalid nonce")
	}
	if tToken != uint16(token) {
		t.Fatal("Invalid token")
	}
	if !bytes.Equal(tSource, srcPub[:]) {
		t.Fatal("Invalid source address")
	}
	if !bytes.Equal(tDestination, dstPub[:]) {
		t.Fatal("Invalid destination address")
	}
	if tTokenAmount == nil || tTokenAmount.Value == nil || tTokenAmount.Value.Cmp(tokenAmount.Value) != 0 {
		t.Fatal("Invalid token amount")
	}
}
