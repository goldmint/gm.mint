package transaction

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/signer"
	"github.com/void616/gm-sumus-lib/types"
	"github.com/void616/gm-sumus-lib/types/amount"
	"golang.org/x/crypto/sha3"
)

// SignedTransaction data
type SignedTransaction struct {
	// Hash, 40b
	Hash []byte
	// Digest, 32b
	Digest []byte
	// Data of the transaction
	Data []byte
	// Signature, 64b
	Signature []byte
}

// ParsedTransaction data
type ParsedTransaction struct {
	// From address, 32b
	From []byte
	// Nonce
	Nonce uint64
	// Hash, 40b
	Hash []byte
	// Digest, 32b
	Digest []byte
	// Signature, 64b
	Signature []byte
}

// ---

// Callback to get transaction payload
type payloadWriter func(s *serializer.Serializer)

// Write nonce and payload, calc a digest and sign it
func construct(signer *signer.Signer, nonce uint64, write payloadWriter) (*SignedTransaction, error) {

	ser := serializer.NewSerializer()

	// write nonce
	ser.PutUint64(nonce)

	// write payload
	write(ser)

	// get payload
	payload, err := ser.Data()
	if err != nil {
		return nil, err
	}

	// make payload digest
	hasher := sha3.New256()
	_, err = hasher.Write(payload)
	if err != nil {
		return nil, err
	}
	txdigest := hasher.Sum(nil)

	// sign digest
	txsignature := signer.Sign(txdigest)

	// signature
	ser.
		PutByte(1).           // append a byte - "signed bit"
		PutBytes(txsignature) // signature

	// data
	txdata, err := ser.Data()
	if err != nil {
		return nil, err
	}

	// transaction hash
	_, txhash, err := PackHash(signer.PublicKey(), nonce)
	if err != nil {
		return nil, err
	}

	return &SignedTransaction{
		Hash:      txhash,
		Data:      txdata,
		Digest:    txdigest,
		Signature: txsignature,
	}, nil
}

// Callback to parse transaction payload. Wants a signer public key or an error
type payloadReader func(d *serializer.Deserializer) ([]byte, error)

// Parse transaction data from bytes
func parse(r io.Reader, read payloadReader) (*ParsedTransaction, error) {

	digestWriter := bytes.NewBuffer(make([]byte, 256))
	des := serializer.NewStreamDeserializer(io.TeeReader(r, digestWriter))

	// read nonce
	txnonce := des.GetUint64()
	if err := des.Error(); err != nil {
		return nil, err
	}

	// read payload, get signer pub key
	txsigner, rerr := read(des)
	if err := des.Error(); err != nil {
		return nil, err
	}
	if rerr != nil {
		return nil, rerr
	}

	// calc the digest
	hasher := sha3.New256()
	_, err := hasher.Write(digestWriter.Bytes())
	if err != nil {
		return nil, err
	}
	txdigest := hasher.Sum(nil)

	// "signed" byte
	txsigned := des.GetByte()
	if err := des.Error(); err != nil {
		return nil, err
	}

	txsignature := make([]byte, 64)
	if txsigned != 0 {
		// signature
		txsignature = des.GetBytes(64)
		if err := des.Error(); err != nil {
			return nil, err
		}
	} else {
		// digest
		_ = des.GetBytes(32)
		if err := des.Error(); err != nil {
			return nil, err
		}
	}

	// TODO: verify optionally

	// make a hash
	_, txhash, err := PackHash(txsigner, txnonce)
	if err != nil {
		return nil, err
	}

	return &ParsedTransaction{
		From:      txsigner,
		Nonce:     txnonce,
		Hash:      txhash,
		Digest:    txdigest,
		Signature: txsignature,
	}, nil
}

// ITransaction is generic interface
type ITransaction interface {
	Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error)
	Parse(r io.Reader) (*ParsedTransaction, error)
}

// ---

// RegisterNode transaction
type RegisterNode struct {
	NodeAddress string
}

// Construct ...
func (t *RegisterNode) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutString64(t.NodeAddress)   // node address
	})
}

// Parse ...
func (t *RegisterNode) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32)           // signer public key
		t.NodeAddress = des.GetString64() // node address
		return ret, nil
	})
}

// ---

// UnregisterNode transaction
type UnregisterNode struct {
}

// Construct ...
func (t *UnregisterNode) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey()) // signer public key
	})
}

// Parse ...
func (t *UnregisterNode) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32) // signer public key
		return ret, nil
	})
}

// ---

// TransferAsset transaction
type TransferAsset struct {
	Address []byte
	Token   types.Token
	Amount  *amount.Amount
}

// Construct ...
func (t *TransferAsset) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	if t.Address == nil || len(t.Address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutUint16(uint16(t.Token))   // token
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(t.Address)          // address / public key
		ser.PutAmount(t.Amount)          // amount
	})
}

// Parse ...
func (t *TransferAsset) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		tokenCode := des.GetUint16() // token
		ret := des.GetBytes(32)      // signer public key
		t.Address = des.GetBytes(32) // address / public key
		t.Amount = des.GetAmount()   // amount

		// ensure token is valid
		if !types.ValidToken(tokenCode) {
			return nil, fmt.Errorf("Unknown token with code `%v`", tokenCode)
		}
		t.Token = types.Token(tokenCode)

		return ret, nil
	})
}

// ---

// UserData transaction
type UserData struct {
	Data []byte
}

// Construct ...
func (t *UserData) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	if t.Data == nil {
		return nil, errors.New("Data is empty")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey())   // signer public key
		ser.PutUint32(uint32(len(t.Data))) // data size
		ser.PutBytes(t.Data)               // data bytes
	})
}

// Parse ...
func (t *UserData) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32)     // signer public key
		size := des.GetUint32()     // data size
		t.Data = des.GetBytes(size) // data bytes
		return ret, nil
	})
}

// ---

// RegisterSysWallet transaction
type RegisterSysWallet struct {
	Address []byte
	Tag     types.WalletTag
}

// Construct ...
func (t *RegisterSysWallet) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	if t.Address == nil || len(t.Address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(t.Address)          // address / public key
		ser.PutByte(uint8(t.Tag))        // tag
	})
}

// Parse ...
func (t *RegisterSysWallet) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32)      // signer public key
		t.Address = des.GetBytes(32) // address / public key
		tagCode := des.GetByte()     // tag

		// ensure tag is valid
		if !types.ValidWalletTag(tagCode) {
			return nil, fmt.Errorf("Unknown wallet tag with code `%v`", tagCode)
		}
		t.Tag = types.WalletTag(tagCode)

		return ret, nil
	})
}

// ---

// UnregisterSysWallet transaction
type UnregisterSysWallet struct {
	Address []byte
	Tag     types.WalletTag
}

// Construct ...
func (t *UnregisterSysWallet) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	if t.Address == nil || len(t.Address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(t.Address)          // address / public key
		ser.PutByte(uint8(t.Tag))        // tag
	})
}

// Parse ...
func (t *UnregisterSysWallet) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32)      // signer public key
		t.Address = des.GetBytes(32) // address / public key
		tagCode := des.GetByte()     // tag

		// ensure tag is valid
		if !types.ValidWalletTag(tagCode) {
			return nil, fmt.Errorf("Unknown wallet tag with code `%v`", tagCode)
		}
		t.Tag = types.WalletTag(tagCode)

		return ret, nil
	})
}

// ---

// DistributionFee transaction
type DistributionFee struct {
	OwnerAddress []byte
	AmountMNT    *amount.Amount
	AmountGOLD   *amount.Amount
}

// Construct ...
func (t *DistributionFee) Construct(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {

	if t.OwnerAddress == nil || len(t.OwnerAddress) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(t.OwnerAddress)     // owner address / public key
		ser.PutAmount(t.AmountMNT)       // mnt amount
		ser.PutAmount(t.AmountGOLD)      // gold amount
	})
}

// Parse ...
func (t *DistributionFee) Parse(r io.Reader) (*ParsedTransaction, error) {

	return parse(r, func(des *serializer.Deserializer) ([]byte, error) {
		ret := des.GetBytes(32)           // signer public key
		t.OwnerAddress = des.GetBytes(32) // owner address / public key
		t.AmountMNT = des.GetAmount()     // mnt amount
		t.AmountGOLD = des.GetAmount()    // gold amount
		return ret, nil
	})
}
