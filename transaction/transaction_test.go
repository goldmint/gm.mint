package transaction

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"

	"github.com/void616/gm-sumus-lib/signer"
	"github.com/void616/gm-sumus-lib/types"
	"github.com/void616/gm-sumus-lib/types/amount"
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
				NodeAddress: "127.0.0.1",
			},
		},
		{
			"UnregisterNode",
			&UnregisterNode{},
		},
		{
			"TransferAsset",
			&TransferAsset{
				Address: signer.PublicKey(),
				Token:   types.TokenGOLD,
				Amount:  amount.NewFloatString("1.666"),
			},
		},
		{
			"UserData",
			&UserData{
				Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
			},
		},
		{
			"RegisterSysWallet",
			&RegisterSysWallet{
				Address: signer.PublicKey(),
				Tag:     types.WalletTagSupervisor,
			},
		},
		{
			"UnregisterSysWallet",
			&UnregisterSysWallet{
				Address: signer.PublicKey(),
				Tag:     types.WalletTagEmission,
			},
		},
		{
			"DistributionFee",
			&DistributionFee{
				OwnerAddress: signer.PublicKey(),
				AmountMNT:    amount.NewFloatString("1.666"),
				AmountGOLD:   amount.NewFloatString("666.1"),
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
