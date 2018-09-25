package amount

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Precision of amount
const Precision = 18

// New amount
func New() *Amount {
	return &Amount{
		Value: big.NewInt(0),
	}
}

// NewInteger amount: 100 => 100.000000000000000000
func NewInteger(i int64) *Amount {
	v := big.NewInt(0).
		Mul(
			big.NewInt(i),
			big.NewInt(0).Exp(big.NewInt(10), big.NewInt(Precision), nil),
		)
	return &Amount{
		Value: v,
	}
}

// NewBig amount: 100 => 0.000000000000000100
func NewBig(i *big.Int) *Amount {
	return &Amount{
		Value: big.NewInt(0).Set(i),
	}
}

// NewFloatString amount:
// "1.000000000000000000123" => 1.000000000000000000
// "1.000000000000000000999" => 1.000000000000000001
func NewFloatString(s string) *Amount {
	f, ok := big.NewRat(1, 1).SetString(s)
	if !ok {
		return nil
	}
	t := strings.Replace(
		f.FloatString(Precision),
		".", "", -1,
	)
	return NewBigString(t, 10)
}

// NewBigString amount:
// "0100" (base 10) => 0.000000000000000100, "100" (base 00) => 0.000000000000000100
// "000A" (base 16) => 0.000000000000000010, "0xA" (base 00) => 0.000000000000000010
// "0144" (base 08) => 0.000000000000000100, "012" (base 00) => 0.000000000000000010
// etc. (see big.SetString())
func NewBigString(s string, base int) *Amount {
	v, ok := big.NewInt(0).SetString(s, base)
	if !ok {
		return nil
	}
	return &Amount{
		Value: v,
	}
}

// ---

// Amount in Sumus
type Amount struct {
	Value *big.Int
}

// String representation: -100.000000000000000123
func (a *Amount) String() string {
	sign := ""
	if a.Value.Cmp(big.NewInt(0)) < 0 {
		sign = "-"
	}
	abs := big.NewInt(0).Abs(a.Value)
	ret := fmt.Sprintf(fmt.Sprintf("%%0%ds", Precision+1), abs.Text(10))
	return fmt.Sprintf("%s%s.%s", sign, ret[:len(ret)-Precision], ret[len(ret)-Precision:])
}

// IsNeg value
func (a *Amount) IsNeg() bool {
	return a.Value.Cmp(big.NewInt(0)) < 0
}

// Integer part as string:
// -123.000000000000456000 => "123" (width 0)
// -123.000000000000456000 => "00123" (width 5)
func (a *Amount) Integer(width uint) string {
	ret := big.NewInt(0).Abs(a.Value)
	ret.Div(ret, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(Precision), nil))
	return fmt.Sprintf(fmt.Sprintf("%%0%ds", width), ret.Text(10))
}

// Fraction part as string:
// -123.000000000000456000 => "456000" (width 0)
// -123.000000000000456000 => "000000000000456000" (width 18)
func (a *Amount) Fraction(width uint) string {
	ret := big.NewInt(0).Abs(a.Value)
	ret.Mod(ret, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(Precision), nil))
	return fmt.Sprintf(fmt.Sprintf("%%0%ds", width), ret.Text(10))
}

// ---

// MarshalJSON ...
func (a *Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON ...
func (a *Amount) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp := NewFloatString(s)
	if tmp == nil {
		return errors.New("Failed to parse amount from `" + s + "`")
	}

	*a = *tmp

	return nil
}
