package transaction

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/amount"
	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/signer"
)

func TestRegisterNode(t *testing.T) {

	spvt, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	sig, _ := signer.NewSignerFromPK(spvt)

	_, tx, err := RegisterNode(sig, 0, "chupachups")
	if err != nil {
		t.Fatal(err)
	}

	if tx != "0100000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201d6368757061636875707300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001f2def544ab5fec51d764951d0000932ff08bd67b3a7c63d78dc9202f9d8ffe4284cb8669b49f7b1376ecd3415d40acfe35844ba22c865ea05807df07aa8d1e01" {
		t.Fatal(tx)
	}

	// ---

	_, tx, err = RegisterNode(sig, 0, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!")
	if err == nil {
		t.Fatal("Should fail due to node name length")
	}
}

func TestUnregisterNode(t *testing.T) {

	spvt, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	sig, _ := signer.NewSignerFromPK(spvt)

	_, tx, err := UnregisterNode(sig, 1)
	if err != nil {
		t.Fatal(err)
	}

	if tx != "0200000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201d01f614b9c093360dedb0d3b8cc5fc5e3f30f2b3020ba5ab16ff8af58bb53aad4c900ee1a7d79d06b573258231e9f6bfde833bc8f80d40c1a79778f0afa2b258904" {
		t.Fatal(tx)
	}
}

func TestTransferAsset(t *testing.T) {

	srcpk, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	src, _ := signer.NewSignerFromPK(srcpk)

	dstpk, _ := sumus.Unpack58("FhM2u3UMtexZ3TU57G6d9iDpcmynBSpzmTZq6YaMPeA6DHFdEht3jcZUDpXyVbXGoXoWiYB9z8QVKjGhZuKCqMGYZE2P6")
	dst, _ := signer.NewSignerFromPK(dstpk)

	_, tx, err := TransferAsset(src, 2, dst.PublicKey(), 0, amount.NewFloatString("1000"))
	if err != nil {
		t.Fatal(err)
	}

	if tx != "03000000000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201df42378223753e3f5410b427d4c49df8dee069d798eb5cfb0a4e3bd197b0797b7000000000000000000000010000000010e4b042527eafe9f5c8d90da41d4e062fd044a84e3c1dbcda9342b4921798d9ee56310dda763c137e0ec4e521d2738249120edc7149018eb15240ba373e6090a" {
		t.Fatal(tx)
	}
}

func TestTransferAssetValidation(t *testing.T) {

	nonce := uint64(2)
	token := sumus.TokenMNT
	tokenAmount := amount.NewFloatString("123.666")

	// ---

	srcpk, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	src, _ := signer.NewSignerFromPK(srcpk)

	dstpk, _ := sumus.Unpack58("FhM2u3UMtexZ3TU57G6d9iDpcmynBSpzmTZq6YaMPeA6DHFdEht3jcZUDpXyVbXGoXoWiYB9z8QVKjGhZuKCqMGYZE2P6")
	dst, _ := signer.NewSignerFromPK(dstpk)

	_, txHex, err := TransferAsset(src, nonce, dst.PublicKey(), token, tokenAmount)
	if err != nil {
		t.Fatal(err)
	}

	txBytes, _ := hex.DecodeString(txHex)

	// ---

	// tx:
	// 0300000000000000 - nonce, 8B --------------------------|
	// 0000 - token, 2B                                       |
	// eea072...01d - pub key from, 32B                       |--- payload, 89 bytes
	// f42378...7b7 - pub key to, 32B                         |
	// 000000000000000000000010000000 - amount, 15B ----------|
	// 01 - signed byte, 1B
	// 0e4b04...90a - signature, 64B

	// get payload and signature
	des := serializer.NewDeserializer(txBytes)
	txPayload := des.GetBytes(89)
	tSigned := des.GetByte()
	tSignature := des.GetBytes(64)
	err = des.Error()
	if err != nil {
		t.Fatal(err, "Failed to get payload and signature")
	}

	// is signed?
	if tSigned != 1 {
		t.Fatal("Is not signed")
	}

	// check if signed
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

	if tNonce != nonce+1 {
		t.Fatal("Invalid nonce")
	}
	if tToken != uint16(token) {
		t.Fatal("Invalid token")
	}
	if !bytes.Equal(tSource, src.PublicKey()) {
		t.Fatal("Invalid source address")
	}
	if !bytes.Equal(tDestination, dst.PublicKey()) {
		t.Fatal("Invalid destination address")
	}
	if tTokenAmount == nil || tTokenAmount.Value == nil || tTokenAmount.Value.Cmp(tokenAmount.Value) != 0 {
		t.Fatal("Invalid token amount")
	}
}

func TestUnpackHash(t *testing.T) {
	tests := []struct {
		name      string
		hash      string
		wantAddr  string
		wantNonce uint64
		wantErr   bool
	}{
		{"ok", "cqG4tLhKKNd4ZirnFv7HqaYKDdD6c8GuUXdoWwgE6TmBZ6eu885fgkT2BEoJ", "qY4dBwxN7LfAjNeVhoJfKsAk8DjtCY9WGBMTeqvRvBJqcThNp", 1, false},
		{"fail", "2XfAbdqgBp69XHZfFPJH54XY4Rh6qPpKXG8e8YK6BgG6yQgBjmdvYJGGZDsrg1BRmjPHq3M7D2H6QsZ3YH2i", "qY4dBwxN7LfAjNeVhoJfKsAk8DjtCY9WGBMTeqvRvBJqcThNp", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotNonce, gotErr := UnpackHash(tt.hash)

			if (gotErr != nil) != tt.wantErr {
				t.Errorf("UnpackHash() got err %v, want %v", (gotErr != nil), tt.wantErr)
			}
			if gotErr == nil {
				if sumus.Pack58(gotAddr) != tt.wantAddr {
					t.Errorf("UnpackHash() got addr %v, want %v", sumus.Pack58(gotAddr), tt.wantAddr)
				} else if gotNonce != tt.wantNonce {
					t.Errorf("UnpackHash() got nonce %v, want %v", gotNonce, tt.wantNonce)
				}
			}
		})
	}
}
