package transaction

import (
	"fmt"
	"io"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/signer"
)

var _ = Transactioner(&SetWalletTag{})

// SetWalletTag transaction data
type SetWalletTag struct {
	Address sumuslib.PublicKey
	Tag     sumuslib.WalletTag
}

// Sign impl
func (t *SetWalletTag) Sign(signer *signer.Signer, nonce uint64) (*SignedTransaction, error) {
	ctor := newConstructor(nonce)
	ctor.PutPublicKey(signer.PublicKey()) // signer public key
	ctor.PutPublicKey(t.Address)          // address / public key
	ctor.PutByte(uint8(t.Tag))            // tag
	return ctor.Sign(signer)
}

// Parse impl
func (t *SetWalletTag) Parse(r io.Reader) (*ParsedTransaction, error) {
	pars, err := newParser(r)
	if err != nil {
		return nil, err
	}
	from := pars.GetPublicKey()     // signer public key
	t.Address = pars.GetPublicKey() // address / public key
	tagCode := pars.GetByte()       // tag
	// ensure tag is valid
	if !sumuslib.ValidWalletTag(tagCode) {
		return nil, fmt.Errorf("unknown wallet tag with code `%v`", tagCode)
	}
	t.Tag = sumuslib.WalletTag(tagCode)
	return pars.Complete(from)
}

// Code impl
func (t *SetWalletTag) Code() Code {
	return SetWalletTagTx
}
