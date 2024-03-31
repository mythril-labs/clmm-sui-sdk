package utils

import (
	"math/big"
	"testing"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
	"github.com/stretchr/testify/assert"
)

func TestComputeSwapStep(t *testing.T) {
	currentSqrtPrice := mustFromString("2402403269835123476612")
	targetSqrtPrice := mustFromString("2379498185825388834695")
	liquidity := mustFromString("644166710458")
	amount := mustFromString("500000")
	feeRate := constants.FeeAmount(10000)

	sqrtPriceX64, amountIn, amountOut, feeAmount, err := ComputeSwapStep(currentSqrtPrice, targetSqrtPrice, liquidity, amount, feeRate)
	assert.NoError(t, err)
	assert.Equal(t, mustFromString("2402162869228603008056"), sqrtPriceX64)
	assert.Equal(t, big.NewInt(495000), amountIn)
	assert.Equal(t, big.NewInt(8394872681), amountOut)
	assert.Equal(t, big.NewInt(5000), feeAmount)
}

func mustFromString(decimal string) *big.Int {
	result, _ := new(big.Int).SetString(decimal, 10)
	return result
}
