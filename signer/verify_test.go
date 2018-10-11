package signer

import "testing"

func TestVerify(t *testing.T) {
	msg := []byte{0x0, 0x1, 0x2, 0x3}
	sig, _ := New()
	s := sig.Sign(msg)
	if Verify(sig.PublicKey(), msg, s) != nil {
		t.Fatal()
	}
}
