package utils

import (
	"math/big"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
)

var MaxFee = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

func ComputeSwapStep(sqrtRatioCurrentX64, sqrtRatioTargetX64, liquidity, amountRemaining *big.Int, feePips uint64) (sqrtRatioNextX64, amountIn, amountOut, feeAmount *big.Int, err error) {
	zeroForOne := sqrtRatioCurrentX64.Cmp(sqrtRatioTargetX64) >= 0
	exactIn := amountRemaining.Cmp(constants.Zero) >= 0

	if exactIn {
		amountRemainingLessFee := new(big.Int).Div(new(big.Int).Mul(amountRemaining, new(big.Int).Sub(MaxFee, big.NewInt(int64(feePips)))), MaxFee)
		if zeroForOne {
			amountIn = GetAmount0Delta(sqrtRatioTargetX64, sqrtRatioCurrentX64, liquidity, true)
		} else {
			amountIn = GetAmount1Delta(sqrtRatioCurrentX64, sqrtRatioTargetX64, liquidity, true)
		}
		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX64 = sqrtRatioTargetX64
		} else {
			sqrtRatioNextX64, err = GetNextSqrtPriceFromInput(sqrtRatioCurrentX64, liquidity, amountRemainingLessFee, zeroForOne)
			if err != nil {
				return
			}
		}
	} else {
		if zeroForOne {
			amountOut = GetAmount1Delta(sqrtRatioTargetX64, sqrtRatioCurrentX64, liquidity, false)
		} else {
			amountOut = GetAmount0Delta(sqrtRatioCurrentX64, sqrtRatioTargetX64, liquidity, false)
		}
		if new(big.Int).Mul(amountRemaining, constants.NegativeOne).Cmp(amountOut) >= 0 {
			sqrtRatioNextX64 = sqrtRatioTargetX64
		} else {
			sqrtRatioNextX64, err = GetNextSqrtPriceFromOutput(sqrtRatioCurrentX64, liquidity, new(big.Int).Mul(amountRemaining, constants.NegativeOne), zeroForOne)
			if err != nil {
				return
			}
		}
	}

	max := sqrtRatioTargetX64.Cmp(sqrtRatioNextX64) == 0

	if zeroForOne {
		if !(max && exactIn) {
			amountIn = GetAmount0Delta(sqrtRatioNextX64, sqrtRatioCurrentX64, liquidity, true)
		}
		if !(max && !exactIn) {
			amountOut = GetAmount1Delta(sqrtRatioNextX64, sqrtRatioCurrentX64, liquidity, false)
		}
	} else {
		if !(max && exactIn) {
			amountIn = GetAmount1Delta(sqrtRatioCurrentX64, sqrtRatioNextX64, liquidity, true)
		}
		if !(max && !exactIn) {
			amountOut = GetAmount0Delta(sqrtRatioCurrentX64, sqrtRatioNextX64, liquidity, false)
		}
	}

	if !exactIn && amountOut.Cmp(new(big.Int).Mul(amountRemaining, constants.NegativeOne)) > 0 {
		amountOut = new(big.Int).Mul(amountRemaining, constants.NegativeOne)
	}

	if exactIn && sqrtRatioNextX64.Cmp(sqrtRatioTargetX64) != 0 {
		// we didn't reach the target, so take the remainder of the maximum input as fee
		feeAmount = new(big.Int).Sub(amountRemaining, amountIn)
	} else {
		feeAmount = MulDivRoundingUp(amountIn, big.NewInt(int64(feePips)), new(big.Int).Sub(MaxFee, big.NewInt(int64(feePips))))
	}

	return
}
