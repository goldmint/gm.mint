package transaction

import (
	"fmt"

	"github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/amount"
	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/signer"
	"golang.org/x/crypto/sha3"
)

// New transaction
func New(signer *signer.Signer, nonce uint64) *Transaction {
	tx := &Transaction{
		nonce:  nonce,
		signer: signer,
		ser:    serializer.NewSerializer(),
	}

	// write nonce
	tx.ser.PutUint64(nonce)

	return tx
}

// ---

// Transaction data
type Transaction struct {
	nonce  uint64
	signer *signer.Signer
	ser    *serializer.Serializer
}

// Construct transaction
func (t *Transaction) Construct() (txhash string, txdata string, err error) {

	txhash = ""
	txdata = ""
	err = nil

	payload, err := t.ser.Data()
	if err != nil {
		return
	}

	// make payload digest
	hasher := sha3.New256()
	_, err = hasher.Write(payload)
	if err != nil {
		return
	}
	digest := hasher.Sum(nil)

	// sign digest
	signature := t.signer.Sign(digest)

	// signature
	t.ser.
		PutByte(1).         // append a byte - "signed bit"
		PutBytes(signature) // signature

	// hex of txdata
	txdata, err = t.ser.Hex()
	if err != nil {
		return
	}

	// transaction hash
	txhash, err = PackHash(t.signer.PublicKey(), t.nonce)
	if err != nil {
		return
	}

	return
}

// ---

// RegisterNode transaction
func RegisterNode(signer *signer.Signer, nonce uint64, nodename string) (string, string, error) {

	tx := New(signer, nonce)

	// payload
	tx.ser.PutBytes(signer.PublicKey()) // public key
	tx.ser.PutString64(nodename)        // node name as 256 bit

	return tx.Construct()
}

// UnregisterNode transaction
func UnregisterNode(signer *signer.Signer, nonce uint64) (string, string, error) {

	tx := New(signer, nonce)

	// payload
	tx.ser.PutBytes(signer.PublicKey()) // public key

	return tx.Construct()
}

// TransferAsset transaction
func TransferAsset(signer *signer.Signer, nonce uint64, address []byte, token sumus.Token, am *amount.Amount) (string, string, error) {

	if address == nil || len(address) != 32 {
		return "", "", fmt.Errorf("Destination address is invalid")
	}

	tx := New(signer, nonce)

	// payload
	tx.ser.PutUint16(uint16(token))     // token
	tx.ser.PutBytes(signer.PublicKey()) // public key
	tx.ser.PutBytes(address)            // address / public key
	tx.ser.PutAmount(am)                // amount

	return tx.Construct()
}

// UserData transaction
func UserData(signer *signer.Signer, nonce uint64, data []byte) (string, string, error) {

	if data == nil {
		return "", "", fmt.Errorf("Data is empty")
	}

	tx := New(signer, nonce)

	// payload
	tx.ser.PutBytes(signer.PublicKey()) // public key
	tx.ser.PutUint32(uint32(len(data))) // data size
	tx.ser.PutBytes(data)               // data bytes

	return tx.Construct()
}
