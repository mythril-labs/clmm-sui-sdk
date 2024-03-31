package utils

import (
	"math/big"
	"testing"

	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/stretchr/testify/assert"
)

func TestEncodeSqrtRatioX64(t *testing.T) {
	assert.Equal(t, EncodeSqrtRatioX64(big.NewInt(1), big.NewInt(1)), constants.Q64, "1/1")

	r0, _ := new(big.Int).SetString("184467440737095516160", 10)
	assert.Equal(t, EncodeSqrtRatioX64(big.NewInt(100), big.NewInt(1)), r0, 10, "100/1")

	r1, _ := new(big.Int).SetString("1844674407370955161", 10)
	assert.Equal(t, EncodeSqrtRatioX64(big.NewInt(1), big.NewInt(100)), r1, 10, "1/100")

	r2, _ := new(big.Int).SetString("10650232656628343401", 10)
	assert.Equal(t, EncodeSqrtRatioX64(big.NewInt(111), big.NewInt(333)), r2, 10, "111/333")

	r3, _ := new(big.Int).SetString("31950697969885030203", 10)
	assert.Equal(t, EncodeSqrtRatioX64(big.NewInt(333), big.NewInt(111)), r3, 10, "333/111")
}
