package entities

import (
	"math/big"
	"testing"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
	"github.com/mythril-labs/clmm-sui-sdk/utils"
	"github.com/stretchr/testify/assert"
)

var (
	USDC     = NewToken(1, "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", 6, "USDC", "USD Coin")
	DAI      = NewToken(1, "0x6B175474E89094C44Da98b954EedeAC495271d0F", 18, "DAI", "Dai Stablecoin")
	OneEther = big.NewInt(1e18)

	WETH9 = map[uint]*Token{
		1:  NewToken(1, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", 18, "WETH", "Wrapped Ether"),
		3:  NewToken(3, "0xc778417E063141139Fce010982780140Aa0cD5Ab", 18, "WETH", "Wrapped Ether"),
		4:  NewToken(4, "0xc778417E063141139Fce010982780140Aa0cD5Ab", 18, "WETH", "Wrapped Ether"),
		5:  NewToken(5, "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6", 18, "WETH", "Wrapped Ether"),
		42: NewToken(42, "0xd0A1E359811322d97991E03f863a0C30C2cF029C", 18, "WETH", "Wrapped Ether"),

		10: NewToken(10, "0x4200000000000000000000000000000000000006", 18, "WETH", "Wrapped Ether"),
		69: NewToken(69, "0x4200000000000000000000000000000000000006", 18, "WETH", "Wrapped Ether"),

		42161:  NewToken(42161, "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1", 18, "WETH", "Wrapped Ether"),
		421611: NewToken(421611, "0xB47e6A5f8b33b3F17603C83a0535A9dcD7E32681", 18, "WETH", "Wrapped Ether"),
	}
)

func TestNewPool(t *testing.T) {
	_, err := NewPool(USDC, WETH9[1], 1e6, 0, utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, ErrFeeTooHigh, "fee cannot be more than 1e6'")

	_, err = NewPool(USDC, WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX64, "price must be within tick price bounds")

	_, err = NewPool(USDC, WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], new(big.Int).Add(utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(1)), big.NewInt(0), -1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX64, "price must be within tick price bounds")

	_, err = NewPool(USDC, WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool medium fee")

	_, err = NewPool(USDC, WETH9[1], constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool low fee")

	_, err = NewPool(USDC, WETH9[1], constants.FeeHigh, constants.TickSpacings[constants.FeeHigh], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool high fee")
}

func TestInvolvesToken(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.True(t, pool.InvolvesToken(USDC), "involves USDC")
	assert.True(t, pool.InvolvesToken(DAI), "involves DAI")
	assert.False(t, pool.InvolvesToken(WETH9[1]), "does not involve WETH9")
}

func newTestPool() *Pool {
	ticks := []Tick{
		{
			Index:          NearestUsableTick(utils.MinTick, constants.TickSpacings[constants.FeeLow]),
			LiquidityNet:   OneEther,
			LiquidityGross: OneEther,
		},
		{
			Index:          NearestUsableTick(utils.MaxTick, constants.TickSpacings[constants.FeeLow]),
			LiquidityNet:   new(big.Int).Mul(OneEther, constants.NegativeOne),
			LiquidityGross: OneEther,
		},
	}

	p, err := NewTickListDataProvider(ticks, constants.TickSpacings[constants.FeeLow])
	if err != nil {
		panic(err)
	}

	pool, err := NewPool(USDC, DAI, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), OneEther, 0, p)
	if err != nil {
		panic(err)
	}
	return pool
}
func TestGetOutputAmount(t *testing.T) {
	pool := newTestPool()

	// USDC -> DAI
	inputAmount := FromRawAmount(USDC, big.NewInt(100))
	outputAmount, _, err := pool.GetOutputAmount(inputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, outputAmount.Currency.Equal(DAI))
	assert.Equal(t, outputAmount.Quotient(), big.NewInt(98))

	// DAI -> USDC
	inputAmount = FromRawAmount(DAI, big.NewInt(100))
	outputAmount, _, err = pool.GetOutputAmount(inputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, outputAmount.Currency.Equal(USDC))
	assert.Equal(t, outputAmount.Quotient(), big.NewInt(98))
}

func TestGetInputAmount(t *testing.T) {
	pool := newTestPool()

	// USDC -> DAI
	outputAmount := FromRawAmount(DAI, big.NewInt(98))
	inputAmount, _, err := pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, inputAmount.Currency.Equal(USDC))
	assert.Equal(t, inputAmount.Quotient(), big.NewInt(100))

	// DAI -> USDC
	outputAmount = FromRawAmount(USDC, big.NewInt(98))
	inputAmount, _, err = pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, inputAmount.Currency.Equal(DAI))
	assert.Equal(t, inputAmount.Quotient(), big.NewInt(100))
}
