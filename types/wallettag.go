package types

import (
	"fmt"
)

// WalletTag in Sumus blockchain
type WalletTag uint16

const (
	// WalletTagNode is node wallet
	WalletTagNode WalletTag = iota
	// WalletTagGenesisNode is node wallet (TODO: clarify)
	WalletTagGenesisNode
	// WalletTagSupervisor is controller wallet who can tag other wallets
	WalletTagSupervisor
	// WalletTagOwner is a fee accumulator
	WalletTagOwner
	// WalletTagEmission emits token without a fee
	WalletTagEmission
	// WalletTagData can send UserData transactions without a fee
	WalletTagData
)

// WalletTagToString definition
var WalletTagToString = map[WalletTag]string{
	WalletTagNode:        "Node",
	WalletTagGenesisNode: "GenesisNode",
	WalletTagSupervisor:  "SupervisorWallet",
	WalletTagOwner:       "OwnerWallet",
	WalletTagEmission:    "EmissionWallet",
	WalletTagData:        "DataWallet",
}

// String representation
func (t WalletTag) String() string {
	ret, ok := WalletTagToString[t]
	if !ok {
		return ""
	}
	return ret
}

// ParseWalletTag from string
func ParseWalletTag(s string) (WalletTag, error) {
	for i, v := range WalletTagToString {
		if s == v {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Unknown wallet tag name `%v`", s)
}

// ValidWalletTag as uint16
func ValidWalletTag(u uint16) bool {
	_, ok := WalletTagToString[WalletTag(u)]
	return ok
}
