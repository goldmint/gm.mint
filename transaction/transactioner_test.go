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
			"register_node",
			&RegisterNode{
				NodeAddress: signer.PublicKey(),
				NodeIP:      "127.0.0.1",
			},
		},
		{
			"unregister_node",
			&UnregisterNode{
				NodeAddress: signer.PublicKey(),
			},
		},
		{
			"transfer_asset",
			&TransferAsset{
				Address: signer.PublicKey(),
				Token:   sumuslib.TokenGOLD,
				Amount:  amount.MustFromString("1.666"),
			},
		},
		{
			"user_data",
			&UserData{
				Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
			},
		},
		{
			"set_wallet_tag",
			&SetWalletTag{
				Address: signer.PublicKey(),
				Tag:     sumuslib.WalletTagSupervisor,
			},
		},
		{
			"unset_wallet_tag",
			&UnsetWalletTag{
				Address: signer.PublicKey(),
				Tag:     sumuslib.WalletTagEmission,
			},
		},
		{
			"distribution_fee",
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
			signed, err := tt.tx.(Transactioner).Sign(signer, rand.Uint64())
			if err != nil {
				t.Errorf("Sign() error = %v", err)
				return
			}

			// parse back
			buf := bytes.NewBuffer(signed.Data)
			tx := reflect.New(reflect.TypeOf(tt.tx).Elem()).Interface()
			_, err = tx.(Transactioner).Parse(buf)
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
