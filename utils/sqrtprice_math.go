package utils

import (
	"errors"
	"math/big"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
)

var (
	ErrSqrtPriceLessThanZero = errors.New("sqrt price less than zero")
	ErrLiquidityLessThanZero = errors.New("liquidity less than zero")
	ErrInvariant             = errors.New("invariant violation")
)
var MaxUint128, _ = new(big.Int).SetString("ffffffffffffffffffffffffffffffff", 16)
var MaxUint256, _ = new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)

func multiplyIn256(x, y *big.Int) *big.Int {
	product := new(big.Int).Mul(x, y)
	return new(big.Int).And(product, MaxUint256)
}

func addIn256(x, y *big.Int) *big.Int {
	sum := new(big.Int).Add(x, y)
	return new(big.Int).And(sum, MaxUint256)
}

func GetAmount0Delta(sqrtRatioAX64, sqrtRatioBX64, liquidity *big.Int, roundUp bool) *big.Int {
	if sqrtRatioAX64.Cmp(sqrtRatioBX64) >= 0 {
		sqrtRatioAX64, sqrtRatioBX64 = sqrtRatioBX64, sqrtRatioAX64
	}

	numerator1 := new(big.Int).Lsh(liquidity, 64)
	numerator2 := new(big.Int).Sub(sqrtRatioBX64, sqrtRatioAX64)

	if roundUp {
		return MulDivRoundingUp(MulDivRoundingUp(numerator1, numerator2, sqrtRatioBX64), constants.One, sqrtRatioAX64)
	}
	return new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(numerator1, numerator2), sqrtRatioBX64), sqrtRatioAX64)
}

func GetAmount1Delta(sqrtRatioAX64, sqrtRatioBX64, liquidity *big.Int, roundUp bool) *big.Int {
	if sqrtRatioAX64.Cmp(sqrtRatioBX64) >= 0 {
		sqrtRatioAX64, sqrtRatioBX64 = sqrtRatioBX64, sqrtRatioAX64
	}

	if roundUp {
		return MulDivRoundingUp(liquidity, new(big.Int).Sub(sqrtRatioBX64, sqrtRatioAX64), constants.Q64)
	}
	return new(big.Int).Div(new(big.Int).Mul(liquidity, new(big.Int).Sub(sqrtRatioBX64, sqrtRatioAX64)), constants.Q64)
}

func GetNextSqrtPriceFromInput(sqrtPX64, liquidity, amountIn *big.Int, zeroForOne bool) (*big.Int, error) {
	if sqrtPX64.Cmp(constants.Zero) <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Cmp(constants.Zero) <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX64, liquidity, amountIn, true)
	}
	return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX64, liquidity, amountIn, true)
}

func GetNextSqrtPriceFromOutput(sqrtPX64, liquidity, amountOut *big.Int, zeroForOne bool) (*big.Int, error) {
	if sqrtPX64.Cmp(constants.Zero) <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Cmp(constants.Zero) <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX64, liquidity, amountOut, false)
	}
	return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX64, liquidity, amountOut, false)
}

func getNextSqrtPriceFromAmount0RoundingUp(sqrtPX64, liquidity, amount *big.Int, add bool) (*big.Int, error) {
	if amount.Cmp(constants.Zero) == 0 {
		return sqrtPX64, nil
	}

	numerator1 := new(big.Int).Lsh(liquidity, 64)
	if add {
		product := multiplyIn256(amount, sqrtPX64)
		if new(big.Int).Div(product, amount).Cmp(sqrtPX64) == 0 {
			denominator := addIn256(numerator1, product)
			if denominator.Cmp(numerator1) >= 0 {
				return MulDivRoundingUp(numerator1, sqrtPX64, denominator), nil
			}
		}
		return MulDivRoundingUp(numerator1, constants.One, new(big.Int).Add(new(big.Int).Div(numerator1, sqrtPX64), amount)), nil
	} else {
		product := multiplyIn256(amount, sqrtPX64)
		if new(big.Int).Div(product, amount).Cmp(sqrtPX64) != 0 {
			return nil, ErrInvariant
		}
		if numerator1.Cmp(product) <= 0 {
			return nil, ErrInvariant
		}
		denominator := new(big.Int).Sub(numerator1, product)
		return MulDivRoundingUp(numerator1, sqrtPX64, denominator), nil
	}
}

func getNextSqrtPriceFromAmount1RoundingDown(sqrtPX64, liquidity, amount *big.Int, add bool) (*big.Int, error) {
	if add {
		var quotient *big.Int
		if amount.Cmp(MaxUint128) <= 0 {
			quotient = new(big.Int).Div(new(big.Int).Lsh(amount, 64), liquidity)
		} else {
			quotient = new(big.Int).Div(new(big.Int).Mul(amount, constants.Q64), liquidity)
		}
		return new(big.Int).Add(sqrtPX64, quotient), nil
	}

	quotient := MulDivRoundingUp(amount, constants.Q64, liquidity)
	if sqrtPX64.Cmp(quotient) <= 0 {
		return nil, ErrInvariant
	}
	return new(big.Int).Sub(sqrtPX64, quotient), nil
}
