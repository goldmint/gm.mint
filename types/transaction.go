package types

import (
	"fmt"
)

// Transaction in Sumus blockchain
type Transaction uint16

const (
	// TransactionRegisterNode registers a new node
	TransactionRegisterNode Transaction = iota
	// TransactionUnregisterNode unregisters existing node
	TransactionUnregisterNode
	// TransactionTransferAssets sends token between wallets
	TransactionTransferAssets
	// TransactionRegisterSystemWallet registers system wallet
	TransactionRegisterSystemWallet
	// TransactionUnregisterSystemWallet unregisters system wallet
	TransactionUnregisterSystemWallet
	// TransactionUserData contains custom payload
	TransactionUserData
	// TransactionDistributionFee does the magic
	TransactionDistributionFee
)

// TransactionToString definition
var TransactionToString = map[Transaction]string{
	TransactionRegisterNode:           "RegisterNodeTransaction",
	TransactionUnregisterNode:         "UnregisterNodeTransaction",
	TransactionTransferAssets:         "TransferAssetsTransaction",
	TransactionRegisterSystemWallet:   "RegisterSystemWalletTransaction",
	TransactionUnregisterSystemWallet: "UnregisterSystemWalletTransaction",
	TransactionUserData:               "UserDataTransaction",
	TransactionDistributionFee:        "DistributionFeeTransaction",
}

// String representation
func (t Transaction) String() string {
	ret, ok := TransactionToString[t]
	if !ok {
		return ""
	}
	return ret
}

// ParseTransaction from string
func ParseTransaction(s string) (Transaction, error) {
	for i, v := range TransactionToString {
		if s == v {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Unknown transaction name `%v`", s)
}

// ValidTransaction as uint16
func ValidTransaction(u uint16) bool {
	_, ok := TransactionToString[Transaction(u)]
	return ok
}
