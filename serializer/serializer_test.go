package serializer

import (
	"testing"

	"github.com/void616/gm.mint/amount"
)

func TestSerializer_PutString64(t *testing.T) {
	s := NewSerializer()
	s.PutString64("asdasdфывфыв")
	_, err := s.Hex()
	if err != nil {
		t.Fatal(err)
	}

	s = NewSerializer()
	s.PutString64("言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言言語でゼロ埋め言")
	_, err = s.Hex()
	if err == nil {
		t.Fatal("Should fail on long string")
	}
}

func TestSerializer_PutAmount(t *testing.T) {
	tests := []struct {
		name string
		s    *Serializer
		args *amount.Amount
		want string
	}{
		{"1", NewSerializer(), amount.MustFromString("1234.000000000000001234"), "00341200000000000000341200000000"},
		{"2", NewSerializer(), amount.MustFromString("-0.123400000000004321"), "01214300000000003412000000000000"},
		{"3", NewSerializer(), amount.MustFromString("1.123456789123456789"), "00896745239178563412010000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.PutAmount(tt.args)
			h, err := tt.s.Hex()
			if err != nil {
				t.Error(err)
			} else if h != tt.want {
				t.Errorf("Serializer.PutAmount() = %v, want %v", h, tt.want)
			}
		})
	}
}
