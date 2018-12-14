package amount

import (
	"encoding/json"
	"math"
	"math/big"
	"testing"
)

func TestAmount_String(t *testing.T) {

	a := NewInteger(1)
	if a.String() != "1.000000000000000000" {
		t.Fatal(a.String())
	}

	a = NewInteger(123)
	if a.String() != "123.000000000000000000" {
		t.Fatal(a.String())
	}

	a = NewBig(big.NewInt(123456))
	if a.String() != "0.000000000000123456" {
		t.Fatal(a.String())
	}

	a = NewBig(big.NewInt(-666))
	if a.String() != "-0.000000000000000666" {
		t.Fatal(a.String())
	}

	a = NewBig(big.NewInt(123456))
	a.Value = a.Value.Add(a.Value, NewInteger(123456).Value)
	a.Value = a.Value.Neg(a.Value)
	if a.String() != "-123456.000000000000123456" {
		t.Fatal(a.String())
	}

	a = NewFloatString("0.1")
	if a.String() != "0.100000000000000000" {
		t.Fatal(a.String())
	}

	a = NewFloatString("-123456.000000000000123456444")
	if a.String() != "-123456.000000000000123456" {
		t.Fatal(a.String())
	}

	a = NewFloatString("-123456.000000000000123456999")
	if a.String() != "-123456.000000000000123457" {
		t.Fatal(a.String())
	}
}

func TestNewBigString(t *testing.T) {
	type args struct {
		s    string
		base int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"dec", args{"01000", 10}, "0.000000000000001000"},
		{"hex", args{"003e8", 16}, "0.000000000000001000"},
		{"oct", args{"01750", 8}, "0.000000000000001000"},
		{"dec-0", args{"-1000", 0}, "-0.000000000000001000"},
		{"hex-0", args{"-0x3e8", 0}, "-0.000000000000001000"},
		{"oct-0", args{"-01750", 0}, "-0.000000000000001000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBigString(tt.args.s, tt.args.base); got.String() != tt.want {
				t.Errorf("NewBigString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmount_Fraction(t *testing.T) {
	tests := []struct {
		name   string
		a      *Amount
		width1 uint
		width2 uint
		want1  string
		want2  string
	}{
		{"1", NewFloatString("0"), 10, Precision, "0000000000", "000000000000000000"},
		{"2", NewFloatString("-123.456"), 0, Precision, "123", "456000000000000000"},
		{"3", NewFloatString("0.000000000000000001"), 0, Precision, "0", "000000000000000001"},
		{"4", NewFloatString("666"), 0, Precision, "666", "000000000000000000"},
		{"5", NewFloatString("616.000000000000000000666"), 10, Precision, "0000000616", "000000000000000001"},
		{"6", NewFloatString("-999999999999999999.111222333444555666444"), 0, Precision, "999999999999999999", "111222333444555666"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Integer(tt.width1); got != tt.want1 {
				t.Errorf("Amount.Integer() = %v, want1 %v", got, tt.want1)
			}
			if got := tt.a.Fraction(tt.width2); got != tt.want2 {
				t.Errorf("Amount.Fraction() = %v, want2 %v", got, tt.want2)
			}
		})
	}
}

func TestAmount_ToFromJson(t *testing.T) {

	a := NewFloatString("-987654321987654321.123456789123456789")

	jbytes, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	b := NewInteger(0)
	err = json.Unmarshal(jbytes, b)
	if err != nil {
		t.Fatal(err)
	}

	if a.Value.Cmp(b.Value) != 0 {
		t.Fatal("a != b")
	}

	// ---

	type Tst struct {
		X *Amount `json:"x,omitempty"`
		Y *Amount `json:"y,omitempty"`
		Z *Amount `json:"z,omitempty"`
	}
	tst := Tst{
		X: NewFloatString("-987654321987654321.123456789123456789"),
		Y: nil,
		Z: NewFloatString("666"),
	}

	jbytes, err = json.Marshal(tst)
	if err != nil {
		t.Fatal(err)
	}

	tst2 := Tst{
		X: NewInteger(0),
		Y: NewInteger(0),
		Z: NewInteger(0),
	}

	err = json.Unmarshal(jbytes, &tst2)
	if err != nil {
		t.Fatal(err)
	}

	if tst.X.Value.Cmp(tst2.X.Value) != 0 {
		t.Fatal("x1 != x2")
	}

	if tst2.Y.Value.Cmp(big.NewInt(0)) != 0 {
		t.Fatal("y2 != 0")
	}

	if tst.Z.Value.Cmp(tst2.Z.Value) != 0 {
		t.Fatal("z1 != z2")
	}
}

func TestAmount_Float64(t *testing.T) {
	tests := []struct {
		a    string
		want float64
	}{
		{"1.0000011", 1.000001},
		{"1234.0000019", 1234.000001},
		{"-1234.0000019", -1234.000001},
		{"123123123123.123123", 123123123123.123123},
		{"123123123123.1231231", 123123123123.12312},
		{"-123123123123.1231239", -123123123123.12312},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := NewFloatString(tt.a).Float64(); got != tt.want {
				t.Errorf("Amount.Float64() = %v, want %v (diff %v)", got, tt.want, math.Abs(got-tt.want))
			}
		})
	}
}
