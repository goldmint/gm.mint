package transaction

import (
	"errors"

	"github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/serializer"
	"github.com/void616/gm-sumus-lib/signer"
	"github.com/void616/gm-sumus-lib/types"
	"github.com/void616/gm-sumus-lib/types/amount"
	"golang.org/x/crypto/sha3"
)

// Transaction data
type Transaction struct {
	// Name
	Name string
	// Nonce
	Nonce uint64
	// Signer public key, packed
	Signer string
	// Hash, packed
	Hash string
	// Digest, packed
	Digest string
	// Data hex
	Data string
}

// ---

type payloadWriter func(s *serializer.Serializer) types.Transaction

func construct(signer *signer.Signer, nonce uint64, write payloadWriter) (*Transaction, error) {

	ser := serializer.NewSerializer()

	// write nonce
	ser.PutUint64(nonce)

	// write payload
	txtype := write(ser)

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
	digest := hasher.Sum(nil)

	// sign digest
	signature := signer.Sign(digest)

	// signature
	ser.
		PutByte(1).         // append a byte - "signed bit"
		PutBytes(signature) // signature

	// hex of txdata
	txdata, err := ser.Hex()
	if err != nil {
		return nil, err
	}

	// transaction hash
	txhash, err := PackHash(signer.PublicKey(), nonce)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Name:   txtype.String(),
		Nonce:  nonce,
		Hash:   txhash,
		Data:   txdata,
		Signer: sumus.Pack58(signer.PublicKey()),
		Digest: sumus.Pack58(digest),
	}, nil
}

// ---

// RegisterNode transaction
func RegisterNode(signer *signer.Signer, nonce uint64, address string) (*Transaction, error) {

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutString64(address)         // node address
		return types.TransactionRegisterNode
	})
}

// UnregisterNode transaction
func UnregisterNode(signer *signer.Signer, nonce uint64) (*Transaction, error) {

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutBytes(signer.PublicKey()) // signer public key
		return types.TransactionUnregisterNode
	})
}

// TransferAsset transaction
func TransferAsset(signer *signer.Signer, nonce uint64, address []byte, token types.Token, am *amount.Amount) (*Transaction, error) {

	if address == nil || len(address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutUint16(uint16(token))     // token
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(address)            // address / public key
		ser.PutAmount(am)                // amount
		return types.TransactionTransferAssets
	})
}

// UserData transaction
func UserData(signer *signer.Signer, nonce uint64, data []byte) (*Transaction, error) {

	if data == nil {
		return nil, errors.New("Data is empty")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutUint32(uint32(len(data))) // data size
		ser.PutBytes(data)               // data bytes
		return types.TransactionUserData
	})
}

// RegisterSysWallet transaction
func RegisterSysWallet(signer *signer.Signer, nonce uint64, address []byte, tag types.WalletTag) (*Transaction, error) {

	if address == nil || len(address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(address)            // address / public key
		ser.PutUint16(uint16(tag))       // tag
		return types.TransactionRegisterSystemWallet
	})
}

// UnregisterSysWallet transaction
func UnregisterSysWallet(signer *signer.Signer, nonce uint64, address []byte, tag types.WalletTag) (*Transaction, error) {

	if address == nil || len(address) != 32 {
		return nil, errors.New("Destination address is invalid")
	}

	return construct(signer, nonce, func(ser *serializer.Serializer) types.Transaction {
		ser.PutBytes(signer.PublicKey()) // signer public key
		ser.PutBytes(address)            // address / public key
		ser.PutUint16(uint16(tag))       // tag
		return types.TransactionUnregisterSystemWallet
	})
}
