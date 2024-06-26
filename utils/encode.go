package utils

import "math/big"

/**
 * Returns the sqrt ratio as a Q64.64 corresponding to a given ratio of amount1 and amount0
 * @param amount1 The numerator amount i.e., the amount of token1
 * @param amount0 The denominator amount i.e., the amount of token0
 * @returns The sqrt ratio
 */
func EncodeSqrtRatioX64(amount1 *big.Int, amount0 *big.Int) *big.Int {
	numerator := new(big.Int).Lsh(amount1, 128)
	denominator := amount0
	ratioX128 := new(big.Int).Div(numerator, denominator)
	return new(big.Int).Sqrt(ratioX128)
}
