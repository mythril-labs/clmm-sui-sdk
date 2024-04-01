package constants

import (
	"math/big"
)

// The default factory enabled fee amounts, denominated in hundredths of bips.

const (
	FeeLowest uint64 = 100
	FeeLow    uint64 = 500
	FeeMedium uint64 = 2500
	FeeHigh   uint64 = 10000

	FeeMax uint64 = 1000000
)

// The default factory tick spacings by fee amount.
var TickSpacings = map[uint64]int{
	FeeLowest: 1,
	FeeLow:    10,
	FeeMedium: 60,
	FeeHigh:   200,
}

var (
	NegativeOne = big.NewInt(-1)
	Zero        = big.NewInt(0)
	One         = big.NewInt(1)

	// used in liquidity amount math
	Q64 = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
)
