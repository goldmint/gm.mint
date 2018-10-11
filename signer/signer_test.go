package signer

import (
	"encoding/hex"
	"testing"

	"github.com/void616/gm-sumus-lib"
)

func TestSignerGeneration(t *testing.T) {

	gen1, _ := New()
	gen2, _ := New()

	if hex.EncodeToString(gen1.PrivateKey()) == hex.EncodeToString(gen2.PrivateKey()) {
		t.Fatal(hex.EncodeToString(gen1.PrivateKey()), "==", hex.EncodeToString(gen2.PrivateKey()))
	}
}

func TestSignerFromSumusPrivateKey(t *testing.T) {

	spvt, _ := sumus.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	spub, _ := sumus.Unpack58("2p6QCcwAMLSSXfFFVQT4vYCe8VPwm3rvK4zdNGAM7zeLBqrVLW")

	sig, _ := FromPK(spvt)

	x := sig.PublicKey()
	if hex.EncodeToString(x) != hex.EncodeToString(spub) {
		t.Fatal(hex.EncodeToString(x), "!=", hex.EncodeToString(spub))
	}
}

func TestSignerFromSumusPrivateKey2(t *testing.T) {

	spvt, _ := sumus.Unpack58("4CdzVBba43H7B12zNoSCE8dz8RM9ggUSagfxPdZ1kQ7hbrXLqNNUwGQiiV1VxU3xuEcj4ybxTZPnjq8BAhBUuJxzU8XxQ")
	spub, _ := sumus.Unpack58("2PztA94iHZdeX8d5hPJbQfUGcN6WWUhfmU6G5ySJQ9cnUueiuk")

	sig, _ := FromPK(spvt)

	x := sig.PublicKey()
	if hex.EncodeToString(x) != hex.EncodeToString(spub) {
		t.Fatal(hex.EncodeToString(x), "!=", hex.EncodeToString(spub))
	}
}
