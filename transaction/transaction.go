package transaction

import (
	"bytes"
	"fmt"

	"github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/amount"
	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/signer"
	"golang.org/x/crypto/sha3"
)

// New transaction
func New(signer *signer.Signer, nonce uint64) *Transaction {
	nonce++

	tx := &Transaction{
		nonce:  nonce,
		signer: signer,
		ser:    serializer.NewSerializer(),
	}

	// write nonce
	tx.ser.PutUint64(nonce)

	return tx
}

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

// ---

// PackHash for specific addr/nonce
func PackHash(addr []byte, nonce uint64) (string, error) {
	if addr == nil || len(addr) != 32 {
		return "", fmt.Errorf("Address should be 32 bytes length")
	}
	ser := serializer.NewSerializer()
	ser.PutBytes(addr)
	ser.PutUint64(nonce)
	b, err := ser.Data()
	if err != nil {
		return "", err
	}
	return sumus.Pack58(b), nil
}

// UnpackHash and get addr/nonce
func UnpackHash(hash string) (addr []byte, nonce uint64, err error) {
	b, err := sumus.Unpack58(hash)
	if err != nil {
		return nil, 0, err
	}
	if len(b) != 40 {
		return nil, 0, fmt.Errorf("Invalid hash length")
	}

	des := serializer.NewDeserializer(b)
	baddr := des.GetBytes(32)
	bnonce := des.GetUint64()
	if des.Error() != nil {
		return nil, 0, des.Error()
	}

	return baddr, bnonce, nil
}

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
