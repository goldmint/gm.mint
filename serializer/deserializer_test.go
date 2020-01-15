package serializer

import (
	"encoding/hex"
	"testing"

	"github.com/void616/gm.mint/amount"
)

func TestDerializer(t *testing.T) {

	var b = byte(142)
	var u16 = uint16(0xDEAD)
	var u32 = uint32(0xDEADBEEF)
	var u64 = uint64(0xDEADBEEF1337C0DE)
	var str64 = "961D2014E3E93AC701A6A5F25824DB66"
	var str64Full = "1EF8C0F73B2370D14330C487A70618E0333EAEBA8313EC87131B8F67D964D097"
	var amo1 = amount.MustFromString("1234567890.123456789123456789")
	var amo2 = amount.MustFromString("-987654321.102030405060708090")
	var amo3 = amount.MustFromString("1000")
	var amo4 = amount.MustFromString("1")
	var amo5 = amount.MustFromString("0")

	ser := NewSerializer()
	ser.PutByte(b)
	ser.PutUint16(u16)
	ser.PutUint32(u32)
	ser.PutUint64(u64)
	ser.PutString64(str64)
	ser.PutString64(str64Full)
	ser.PutAmount(amo1)
	ser.PutAmount(amo2)
	ser.PutAmount(amo3)
	ser.PutAmount(amo4)
	ser.PutAmount(amo5)
	datHex, err := ser.Hex()
	if err != nil {
		t.Fatal(err)
	}

	// ---

	datBytes, err := hex.DecodeString(datHex)
	if err != nil {
		t.Fatal(err)
	}

	des := NewDeserializer(datBytes)
	if des.GetByte() != b {
		t.Fatal(des.Error())
	}
	if des.GetUint16() != u16 {
		t.Fatal(des.Error())
	}
	if des.GetUint32() != u32 {
		t.Fatal(des.Error())
	}
	if des.GetUint64() != u64 {
		t.Fatal(des.Error())
	}
	if des.GetString64() != str64 {
		t.Fatal(des.Error())
	}
	if des.GetString64() != str64Full {
		t.Fatal(des.Error())
	}
	damo1 := des.GetAmount()
	if damo1 == nil {
		t.Fatal(des.Error())
	}
	if damo1.Value.Cmp(amo1.Value) != 0 {
		t.Fatal(damo1.String(), "!=", amo1.String())
	}
	damo2 := des.GetAmount()
	if damo2 == nil {
		t.Fatal(des.Error())
	}
	if damo2.Value.Cmp(amo2.Value) != 0 {
		t.Fatal(damo2.String(), "!=", amo2.String())
	}
	damo3 := des.GetAmount()
	if damo3 == nil {
		t.Fatal(des.Error())
	}
	if damo3.Value.Cmp(amo3.Value) != 0 {
		t.Fatal(damo3.String(), "!=", amo3.String())
	}
	damo4 := des.GetAmount()
	if damo4 == nil {
		t.Fatal(des.Error())
	}
	if damo4.Value.Cmp(amo4.Value) != 0 {
		t.Fatal(damo4.String(), "!=", amo4.String())
	}
	damo5 := des.GetAmount()
	if damo5 == nil {
		t.Fatal(des.Error())
	}
	if damo5.Value.Cmp(amo5.Value) != 0 {
		t.Fatal(damo5.String(), "!=", amo5.String())
	}

	if des.Error() != nil {
		t.Fatal(des.Error())
	}
}
