package transaction

import (
	"testing"

	"github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/signer"
	"github.com/void616/gm-sumus-lib/types"
	"github.com/void616/gm-sumus-lib/types/amount"
)

func TestRegisterNode(t *testing.T) {

	spvt, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	sig, _ := signer.NewSignerFromPK(spvt)

	tx, err := RegisterNode(sig, 1, "chupachups")
	if err != nil {
		t.Fatal(err)
	}

	if tx.Data != "0100000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201d6368757061636875707300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001f2def544ab5fec51d764951d0000932ff08bd67b3a7c63d78dc9202f9d8ffe4284cb8669b49f7b1376ecd3415d40acfe35844ba22c865ea05807df07aa8d1e01" {
		t.Fatal(tx)
	}

	// ---

	_, err = RegisterNode(sig, 0, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!")
	if err == nil {
		t.Fatal("Should fail due to node name length")
	}
}

func TestUnregisterNode(t *testing.T) {

	spvt, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	sig, _ := signer.NewSignerFromPK(spvt)

	tx, err := UnregisterNode(sig, 2)
	if err != nil {
		t.Fatal(err)
	}

	if tx.Data != "0200000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201d01f614b9c093360dedb0d3b8cc5fc5e3f30f2b3020ba5ab16ff8af58bb53aad4c900ee1a7d79d06b573258231e9f6bfde833bc8f80d40c1a79778f0afa2b258904" {
		t.Fatal(tx.Data)
	}
}

func TestTransferAsset(t *testing.T) {

	srcpk, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	src, _ := signer.NewSignerFromPK(srcpk)

	dstpk, _ := sumus.Unpack58("FhM2u3UMtexZ3TU57G6d9iDpcmynBSpzmTZq6YaMPeA6DHFdEht3jcZUDpXyVbXGoXoWiYB9z8QVKjGhZuKCqMGYZE2P6")
	dst, _ := signer.NewSignerFromPK(dstpk)

	tx, err := TransferAsset(src, 3, dst.PublicKey(), types.TokenMNT, amount.NewFloatString("1000"))
	if err != nil {
		t.Fatal(err)
	}

	if tx.Data != "03000000000000000000eea0728dfee30d6a65ff2e5c07ddbc4c304cc9005abe2640822adc1ec944201df42378223753e3f5410b427d4c49df8dee069d798eb5cfb0a4e3bd197b0797b7000000000000000000000010000000010e4b042527eafe9f5c8d90da41d4e062fd044a84e3c1dbcda9342b4921798d9ee56310dda763c137e0ec4e521d2738249120edc7149018eb15240ba373e6090a" {
		t.Fatal(tx.Data)
	}
}
