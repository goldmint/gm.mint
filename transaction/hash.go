package transaction

import (
	"fmt"

	sumus "github.com/void616/gm-sumus-lib"
	"github.com/void616/gm-sumus-lib/serializer"
)

// PackHash for specific addr/nonce
func PackHash(addr []byte, nonce uint64) (string, []byte, error) {
	if addr == nil || len(addr) != 32 {
		return "", nil, fmt.Errorf("Address should be 32 bytes length")
	}
	ser := serializer.NewSerializer()
	ser.PutBytes(addr)
	ser.PutUint64(nonce)
	b, err := ser.Data()
	if err != nil {
		return "", nil, err
	}
	return sumus.Pack58(b), b, nil
}

// UnpackHash and get addr/nonce
func UnpackHash(hash string) (addr []byte, nonce uint64, err error) {
	b, err := sumus.Unpack58(hash)
	if err != nil {
		return nil, 0, err
	}
	if len(b) != 40 {
		return nil, 0, fmt.Errorf("Invalid hash length")
	}

	des := serializer.NewDeserializer(b)
	baddr := des.GetBytes(32)
	bnonce := des.GetUint64()
	if des.Error() != nil {
		return nil, 0, des.Error()
	}

	return baddr, bnonce, nil
}
