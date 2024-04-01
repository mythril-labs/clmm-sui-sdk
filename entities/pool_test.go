package entities

import (
	"math/big"
	"testing"

	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mythril-labs/clmm-sui-sdk/constants"
	"github.com/mythril-labs/clmm-sui-sdk/utils"
	"github.com/stretchr/testify/assert"
)

var (
	USDC     = entities.NewToken(1, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 6, "USDC", "USD Coin")
	DAI      = entities.NewToken(1, common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")
	OneEther = big.NewInt(1e18)
)

func TestNewPool(t *testing.T) {
	_, err := NewPool(USDC, entities.WETH9[3], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, entities.ErrDifferentChain, "cannot be used for tokens on different chains")

	_, err = NewPool(USDC, entities.WETH9[1], 1e6, 0, utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, ErrFeeTooHigh, "fee cannot be more than 1e6'")

	_, err = NewPool(USDC, USDC, constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, entities.ErrSameAddress, "cannot be used for the same token")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX64, "price must be within tick price bounds")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], new(big.Int).Add(utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(1)), big.NewInt(0), -1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX64, "price must be within tick price bounds")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, constants.TickSpacings[constants.FeeMedium], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool medium fee")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool low fee")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeHigh, constants.TickSpacings[constants.FeeHigh], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool high fee")
}
func TestToken0(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token0, DAI, "always is the token that sorts before")

	pool, _ = NewPool(DAI, USDC, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token0, DAI, "always is the token that sorts before")
}

func TestToken1(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token1, USDC, "always is the token that sorts after")

	pool, _ = NewPool(DAI, USDC, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token1, USDC, "always is the token that sorts after")
}

func TestInvolvesToken(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, constants.TickSpacings[constants.FeeLow], utils.EncodeSqrtRatioX64(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.True(t, pool.InvolvesToken(USDC), "involves USDC")
	assert.True(t, pool.InvolvesToken(DAI), "involves DAI")
	assert.False(t, pool.InvolvesToken(entities.WETH9[1]), "does not involve WETH9")
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
	inputAmount := entities.FromRawAmount(USDC, big.NewInt(100))
	outputAmount, _, err := pool.GetOutputAmount(inputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, outputAmount.Currency.Equal(DAI))
	assert.Equal(t, outputAmount.Quotient(), big.NewInt(98))

	// DAI -> USDC
	inputAmount = entities.FromRawAmount(DAI, big.NewInt(100))
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
	outputAmount := entities.FromRawAmount(DAI, big.NewInt(98))
	inputAmount, _, err := pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, inputAmount.Currency.Equal(USDC))
	assert.Equal(t, inputAmount.Quotient(), big.NewInt(100))

	// DAI -> USDC
	outputAmount = entities.FromRawAmount(USDC, big.NewInt(98))
	inputAmount, _, err = pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, inputAmount.Currency.Equal(DAI))
	assert.Equal(t, inputAmount.Quotient(), big.NewInt(100))
}
