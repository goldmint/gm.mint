package mnemonic

import (
	"bytes"
	"testing"
)

func TestMnemonic(t *testing.T) {

	// new
	phrase, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("Phrase: %v", phrase)

	// get private key from phrase
	pk, err := Recover(phrase, "password")
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("Private Key: %v", pk.String())

	// phrase is valid
	if !Valid(phrase) {
		t.Fatal("Should be valid phrase")
	}

	// phrase is invalid
	if Valid(phrase + " hello") {
		t.Fatal("Should be invalid phrase")
	}

	// recover with invalid phrase
	if _, err := Recover(phrase+" hello", "password"); err == nil {
		t.Fatal("Should fail on invalid phrase")
	}

	// recover with valid phrase but wrong password
	wrongPK, err := Recover(phrase, "wrong password")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(wrongPK.Bytes(), pk.Bytes()) {
		t.Fatal("Should be different private key")
	}

	// recover with valid phrase but empty password
	wrongPK, err = Recover(phrase, "")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(wrongPK.Bytes(), pk.Bytes()) {
		t.Fatal("Should be different private key")
	}

	// recover again with
	samePK, err := Recover(phrase, "password")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(samePK.Bytes(), pk.Bytes()) {
		t.Fatal("Should be same private key")
	}
}
