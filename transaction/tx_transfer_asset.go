package transaction

import (
	"fmt"
	"io"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/amount"
	"github.com/void616/gm-sumuslib/signer"
)

var _ = Transactioner(&TransferAsset{})

// TransferAsset transaction data
type TransferAsset struct {
	Address sumuslib.PublicKey
	Token   sumuslib.Token
	Amount  *amount.Amount
}

// Sign impl
func (t *TransferAsset) Sign(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {
	ctor := newConstructor(nonce)
	ctor.PutUint16(uint16(t.Token))       // token
	ctor.PutPublicKey(signer.PublicKey()) // signer public key
	ctor.PutPublicKey(t.Address)          // address / public key
	ctor.PutAmount(t.Amount)              // amount
	return ctor.Sign(signer)
}

// Parse impl
func (t *TransferAsset) Parse(r io.Reader) (*ParsedTransaction, error) {
	pars, err := newParser(r)
	if err != nil {
		return nil, err
	}
	tokenCode := pars.GetUint16()   // token
	from := pars.GetPublicKey()     // signer public key
	t.Address = pars.GetPublicKey() // address / public key
	t.Amount = pars.GetAmount()     // amount
	// ensure token is valid
	if !sumuslib.ValidToken(tokenCode) {
		return nil, fmt.Errorf("unknown token with code `%v`", tokenCode)
	}
	t.Token = sumuslib.Token(tokenCode)
	return pars.Complete(from)
}

// Code impl
func (t *TransferAsset) Code() Code {
	return TransferAssetTx
}
