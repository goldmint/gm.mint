package fee

import (
	"math/big"

	"github.com/void616/gm.mint/amount"
)

var (
	goldMinFixed = amount.MustFromString("0.00002")
	goldMaxFixed = amount.MustFromString("0.002")
	mntFixed     = amount.MustFromString("0.02")
	mnt10        = amount.MustFromString("10")
	mnt1_000     = amount.MustFromString("1000")
	mnt10_000    = amount.MustFromString("10000")
	mntPerByte   = amount.MustFromString("0.004")
)

var (
	zero = big.NewInt(0)
	ten  = big.NewInt(10)
	five = big.NewInt(5)
)

// GoldFee estimates fee for a transaction sending `principalGold` GOLD from a sender with `balanceMNT` MNT balance
func GoldFee(principalGOLD *amount.Amount, balanceMNT *amount.Amount) (feeGOLD *amount.Amount) {
	ret := new(big.Int).Set(principalGOLD.Value)

	switch {
	// at least 10 000 MNT -> 0.003%, max fee 0.002 GOLD
	case balanceMNT.Value.Cmp(mnt10_000.Value) >= 0:
		ret.Mul(ret, big.NewInt(3))
		div(ret, 100_000)
		if ret.Cmp(goldMaxFixed.Value) > 0 {
			ret.Set(goldMaxFixed.Value)
		}
	// at least 1 000 MNT -> 0.003%
	case balanceMNT.Value.Cmp(mnt1_000.Value) >= 0:
		ret.Mul(ret, big.NewInt(3))
		div(ret, 100_000)
	// at least 10 MNT -> 0.03%
	case balanceMNT.Value.Cmp(mnt10.Value) >= 0:
		ret.Mul(ret, big.NewInt(3))
		div(ret, 10_000)
	// less than 10 MNT -> 0.1%
	default:
		div(ret, 1_000)
	}

	// min fee 0.00002 GOLD
	if ret.Cmp(goldMinFixed.Value) < 0 {
		ret.Set(goldMinFixed.Value)
	}

	return amount.FromBig(ret)
}

// MntFee estimates fee for a transaction sending `principalMNT` MNT
func MntFee(principalMNT *amount.Amount) (feeMNT *amount.Amount) {
	return amount.FromAmount(mntFixed)
}

// UserDataFee estimates fee (in MNT) for a user-data transaction with payload message length of `messageSize` bytes
func UserDataFee(messageSize uint32) (feeMNT *amount.Amount) {
	ret := big.NewInt(int64(messageSize))
	ret.Mul(ret, mntPerByte.Value)
	return amount.FromBig(ret)
}

// PurgeGold estimates address clearing transaction (both principal and fee, in GOLD) from an sender with `balanceMNT` MNT balance.
// Returned `ok` is false if the transaction is impossible
func PurgeGold(balanceGOLD *amount.Amount, balanceMNT *amount.Amount) (principalGOLD, feeGOLD *amount.Amount, ok bool) {
	g := new(big.Int).Set(balanceGOLD.Value)

	// min fee 0.00002 GOLD
	if g.Cmp(goldMinFixed.Value) <= 0 {
		return
	}

	f := new(big.Int).Set(g)

	switch {
	// at least 10 000 MNT -> 0.003%, max fee 0.002 GOLD
	case balanceMNT.Value.Cmp(mnt10_000.Value) >= 0:
		div(f.Mul(f, big.NewInt(1000)), 100003)
		div(f.Mul(f, big.NewInt(3)), 1000)
		if f.Cmp(goldMaxFixed.Value) > 0 {
			f.Set(goldMaxFixed.Value)
		}
	// at least 1 000 MNT -> 0.003%
	case balanceMNT.Value.Cmp(mnt1_000.Value) >= 0:
		div(f.Mul(f, big.NewInt(1000)), 100003)
		div(f.Mul(f, big.NewInt(3)), 1000)
	// at least 10 MNT -> 0.03%
	case balanceMNT.Value.Cmp(mnt10.Value) >= 0:
		div(f.Mul(f, big.NewInt(100)), 10003)
		div(f.Mul(f, big.NewInt(3)), 100)
	// less than 10 MNT -> 0.1%
	default:
		div(f, 1001)
	}

	// min fee 0.00002 GOLD
	if f.Cmp(goldMinFixed.Value) < 0 {
		f.Set(goldMinFixed.Value)
	}
	p := new(big.Int).Sub(g, f)

	principalGOLD = amount.FromBig(p)
	feeGOLD = amount.FromBig(f)
	ok = principalGOLD.Value.Cmp(zero) > 0
	return
}

// PurgeMnt estimates address clearing transaction (both principal and fee, in MNT).
// Returned `ok` is false if the transaction is impossible
func PurgeMnt(balanceMNT *amount.Amount) (principalMNT, feeMNT *amount.Amount, ok bool) {
	m := new(big.Int).Set(balanceMNT.Value)

	// min fee 0.02 MNT
	if m.Cmp(mntFixed.Value) <= 0 {
		return
	}

	principalMNT = amount.FromBig(new(big.Int).Sub(m, mntFixed.Value))
	feeMNT = amount.FromBig(new(big.Int).Set(mntFixed.Value))
	ok = true
	return
}

// ---

func div(x *big.Int, y int64) {
	x.Mul(x, ten)
	x.Div(x, big.NewInt(y))
	m := new(big.Int).Mod(x, ten)
	x.Div(x, ten)
	if m.Cmp(five) >= 0 {
		x.Add(x, big.NewInt(1))
		return
	}
}
