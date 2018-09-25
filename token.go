package sumus

import (
	"fmt"
	"strings"
)

// Token (asset) in Sumus blockchain
type Token uint16

const (
	// TokenMNT is MNT
	TokenMNT Token = iota
	// TokenGOLD is GOLD
	TokenGOLD
)

// String representation
func (t Token) String() string {
	switch t {
	case TokenMNT:
		return "MNT"
	case TokenGOLD:
		return "GOLD"
	}
	return ""
}

// ParseToken from string
func ParseToken(s string) (Token, error) {
	s = strings.ToLower(s)
	switch s {
	case "0", "utility", "mnt", "mint":
		return TokenMNT, nil
	case "1", "commodity", "gold":
		return TokenGOLD, nil
	}
	return 0, fmt.Errorf("Unknown token `%v`", s)
}

// ToToken from uint16
func ToToken(u uint16) (Token, error) {
	switch u {
	case 0:
		return TokenMNT, nil
	case 1:
		return TokenGOLD, nil
	}
	return 0, fmt.Errorf("Unknown token `%v`", u)
}
