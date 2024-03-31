package utils

import (
	"math/big"
	"testing"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetSqrtRatioAtTick(t *testing.T) {
	_, err := GetSqrtRatioAtTick(MinTick - 1)
	assert.ErrorIs(t, err, ErrInvalidTick, "tick tool small")

	_, err = GetSqrtRatioAtTick(MaxTick + 1)
	assert.ErrorIs(t, err, ErrInvalidTick, "tick tool large")

	rmax, _ := GetSqrtRatioAtTick(MinTick)
	assert.Equal(t, rmax, MinSqrtRatio, "returns the correct value for min tick")

	r0, _ := GetSqrtRatioAtTick(0)
	assert.Equal(t, r0, new(big.Int).Lsh(constants.One, 64), "returns the correct value for tick 0")

	rmin, _ := GetSqrtRatioAtTick(MaxTick)
	assert.Equal(t, rmin, MaxSqrtRatio, "returns the correct value for max tick")

	// Custom test: tick#97674
	tick, _ := GetSqrtRatioAtTick(97674)
	sqrtRatio, _ := new(big.Int).SetString("2436562986624311242090", 10)
	assert.Equal(t, sqrtRatio, tick, "returns the correct value")

	// Custom test: tick#97675
	tick, _ = GetSqrtRatioAtTick(97675)
	sqrtRatio, _ = new(big.Int).SetString("2436684811728091000041", 10)
	assert.Equal(t, sqrtRatio, tick, "returns the correct value")

	// Custom test: tick#-97674
	tick, _ = GetSqrtRatioAtTick(-97674)
	sqrtRatio, _ = new(big.Int).SetString("139656708564048263", 10)
	assert.Equal(t, sqrtRatio, tick, "returns the correct value")

	// Custom test: tick#-97675
	tick, _ = GetSqrtRatioAtTick(-97675)
	sqrtRatio, _ = new(big.Int).SetString("139649726252289079", 10)
	assert.Equal(t, sqrtRatio, tick, "returns the correct value")
}

func TestGetTickAtSqrtRatio(t *testing.T) {
	tmin, err := GetTickAtSqrtRatio(MinSqrtRatio)
	assert.NoError(t, err)
	assert.Equal(t, tmin, MinTick, "returns the correct value for sqrt ratio at min tick")

	tmax, err := GetTickAtSqrtRatio(new(big.Int).Sub(MaxSqrtRatio, constants.One))
	assert.NoError(t, err)
	assert.Equal(t, tmax, MaxTick-1, "returns the correct value for sqrt ratio at max tick")
}
