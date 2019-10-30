package transaction

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/amount"
	"github.com/void616/gm-sumuslib/signer"
)

func TestConstructParse(t *testing.T) {

	signer, _ := signer.New()

	tests := []struct {
		name string
		tx   interface{}
	}{
		{
			"RegisterNode",
			&RegisterNode{
				NodeAddress: signer.PublicKey(),
				NodeIP:      "127.0.0.1",
			},
		},
		{
			"UnregisterNode",
			&UnregisterNode{
				NodeAddress: signer.PublicKey(),
			},
		},
		{
			"TransferAsset",
			&TransferAsset{
				Address: signer.PublicKey(),
				Token:   sumuslib.TokenGOLD,
				Amount:  amount.MustFromString("1.666"),
			},
		},
		{
			"UserData",
			&UserData{
				Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
			},
		},
		{
			"SetWalletTag",
			&SetWalletTag{
				Address: signer.PublicKey(),
				Tag:     sumuslib.WalletTagSupervisor,
			},
		},
		{
			"UnsetWalletTag",
			&UnsetWalletTag{
				Address: signer.PublicKey(),
				Tag:     sumuslib.WalletTagEmission,
			},
		},
		{
			"DistributionFee",
			&DistributionFee{
				OwnerAddress: signer.PublicKey(),
				AmountMNT:    amount.MustFromString("1.666"),
				AmountGOLD:   amount.MustFromString("666.1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// construct and sign
			signed, err := tt.tx.(ITransaction).Construct(signer, rand.Uint64())
			if err != nil {
				t.Errorf("Construct() error = %v", err)
				return
			}

			// parse back
			buf := bytes.NewBuffer(signed.Data)
			tx := reflect.New(reflect.TypeOf(tt.tx).Elem()).Interface()
			_, err = tx.(ITransaction).Parse(buf)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if !reflect.DeepEqual(tt.tx, tx) {
				t.Errorf("Constructed and parsed are not equal: %#v != %#v", tt.tx, tx)
				return
			}
		})
	}
}
