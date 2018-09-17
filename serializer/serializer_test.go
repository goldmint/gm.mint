package serializer

import (
	"testing"

	"github.com/void616/gm-sumus-lib/amount"
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
		{"1", NewSerializer(), amount.NewFloatString("1234.000000000000001234"), "003412000000000000003412000000"},
		{"2", NewSerializer(), amount.NewFloatString("-0.123400000000004321"), "012143000000000034120000000000"},
		{"3", NewSerializer(), amount.NewFloatString("1.123456789123456789"), "008967452391785634120100000000"},
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
