package entities

import (
	"errors"
	"math/big"

	"github.com/mythril-labs/clmm-sui-sdk/constants"
	"github.com/mythril-labs/clmm-sui-sdk/utils"
)

var (
	ErrFeeTooHigh               = errors.New("fee too high")
	ErrInvalidSqrtRatioX64      = errors.New("invalid sqrtRatioX64")
	ErrTokenNotInvolved         = errors.New("token not involved in pool")
	ErrSqrtPriceLimitX64TooLow  = errors.New("SqrtPriceLimitX64 too low")
	ErrSqrtPriceLimitX64TooHigh = errors.New("SqrtPriceLimitX64 too high")
)

type StepComputations struct {
	sqrtPriceStartX64 *big.Int
	tickNext          int
	initialized       bool
	sqrtPriceNextX64  *big.Int
	amountIn          *big.Int
	amountOut         *big.Int
	feeAmount         *big.Int
}

// Represents a V3 pool
type Pool struct {
	Token0           *Token
	Token1           *Token
	Fee              uint64
	TickSpacing      int
	SqrtRatioX64     *big.Int
	Liquidity        *big.Int
	TickCurrent      int
	TickDataProvider TickDataProvider
}

/**
 * Construct a pool
 * @param tokenA One of the tokens in the pool
 * @param tokenB The other token in the pool
 * @param fee The fee in hundredths of a bips of the input amount of every swap that is collected by the pool
 * @param sqrtRatioX64 The sqrt of the current ratio of amounts of token1 to token0
 * @param liquidity The current value of in range liquidity
 * @param tickCurrent The current tick of the pool
 * @param ticks The current state of the pool ticks or a data provider that can return tick data
 */
func NewPool(tokenA, tokenB *Token, fee uint64, tickSpacing int, sqrtRatioX64 *big.Int, liquidity *big.Int, tickCurrent int, ticks TickDataProvider) (*Pool, error) {
	if fee >= constants.FeeMax {
		return nil, ErrFeeTooHigh
	}

	tickCurrentSqrtRatioX64, err := utils.GetSqrtRatioAtTick(tickCurrent)
	if err != nil {
		return nil, err
	}
	nextTickSqrtRatioX64, err := utils.GetSqrtRatioAtTick(tickCurrent + 1)
	if err != nil {
		return nil, err
	}

	if sqrtRatioX64.Cmp(tickCurrentSqrtRatioX64) < 0 || sqrtRatioX64.Cmp(nextTickSqrtRatioX64) > 0 {
		return nil, ErrInvalidSqrtRatioX64
	}

	return &Pool{
		Token0:           tokenA,
		Token1:           tokenB,
		Fee:              fee,
		TickSpacing:      tickSpacing,
		SqrtRatioX64:     sqrtRatioX64,
		Liquidity:        liquidity,
		TickCurrent:      tickCurrent,
		TickDataProvider: ticks, // TODO: new tick data provider
	}, nil
}

/**
 * Returns true if the token is either token0 or token1
 * @param token The token to check
 * @returns True if token is either token0 or token
 */
func (p *Pool) InvolvesToken(token *Token) bool {
	return p.Token0.Equal(token) || p.Token1.Equal(token)
}

/**
 * Given an input amount of a token, return the computed output amount, and a pool with state updated after the trade
 * @param inputAmount The input amount for which to quote the output amount
 * @param sqrtPriceLimitX64 The Q64.64 sqrt price limit
 * @returns The output amount and the pool with updated state
 */
func (p *Pool) GetOutputAmount(inputAmount *CurrencyAmount, sqrtPriceLimitX64 *big.Int) (*CurrencyAmount, *Pool, error) {
	if !(inputAmount.Currency.IsToken() && p.InvolvesToken(inputAmount.Currency.Wrapped())) {
		return nil, nil, ErrTokenNotInvolved
	}
	zeroForOne := inputAmount.Currency.Equal(p.Token0)
	outputAmount, sqrtRatioX64, liquidity, tickCurrent, err := p.swap(zeroForOne, inputAmount.Quotient(), sqrtPriceLimitX64)
	if err != nil {
		return nil, nil, err
	}
	var outputToken *Token
	if zeroForOne {
		outputToken = p.Token1
	} else {
		outputToken = p.Token0
	}
	pool, err := NewPool(p.Token0, p.Token1, p.Fee, p.TickSpacing, sqrtRatioX64, liquidity, tickCurrent, p.TickDataProvider)
	if err != nil {
		return nil, nil, err
	}
	return FromRawAmount(outputToken, new(big.Int).Mul(outputAmount, constants.NegativeOne)), pool, nil
}

/**
 * Given a desired output amount of a token, return the computed input amount and a pool with state updated after the trade
 * @param outputAmount the output amount for which to quote the input amount
 * @param sqrtPriceLimitX64 The Q64.64 sqrt price limit. If zero for one, the price cannot be less than this value after the swap. If one for zero, the price cannot be greater than this value after the swap
 * @returns The input amount and the pool with updated state
 */
func (p *Pool) GetInputAmount(outputAmount *CurrencyAmount, sqrtPriceLimitX64 *big.Int) (*CurrencyAmount, *Pool, error) {
	if !(outputAmount.Currency.IsToken() && p.InvolvesToken(outputAmount.Currency.Wrapped())) {
		return nil, nil, ErrTokenNotInvolved
	}
	zeroForOne := outputAmount.Currency.Equal(p.Token1)
	inputAmount, sqrtRatioX64, liquidity, tickCurrent, err := p.swap(zeroForOne, new(big.Int).Mul(outputAmount.Quotient(), constants.NegativeOne), sqrtPriceLimitX64)
	if err != nil {
		return nil, nil, err
	}
	var inputToken *Token
	if zeroForOne {
		inputToken = p.Token0
	} else {
		inputToken = p.Token1
	}
	pool, err := NewPool(p.Token0, p.Token1, p.Fee, p.TickSpacing, sqrtRatioX64, liquidity, tickCurrent, p.TickDataProvider)
	if err != nil {
		return nil, nil, err
	}
	return FromRawAmount(inputToken, inputAmount), pool, nil
}

/**
 * Executes a swap
 * @param zeroForOne Whether the amount in is token0 or token1
 * @param amountSpecified The amount of the swap, which implicitly configures the swap as exact input (positive), or exact output (negative)
 * @param sqrtPriceLimitX64 The Q64.64 sqrt price limit. If zero for one, the price cannot be less than this value after the swap. If one for zero, the price cannot be greater than this value after the swap
 * @returns amountCalculated
 * @returns sqrtRatioX64
 * @returns liquidity
 * @returns tickCurrent
 */
func (p *Pool) swap(zeroForOne bool, amountSpecified, sqrtPriceLimitX64 *big.Int) (amountCalCulated *big.Int, sqrtRatioX64 *big.Int, liquidity *big.Int, tickCurrent int, err error) {
	if sqrtPriceLimitX64 == nil {
		if zeroForOne {
			sqrtPriceLimitX64 = new(big.Int).Add(utils.MinSqrtRatio, constants.One)
		} else {
			sqrtPriceLimitX64 = new(big.Int).Sub(utils.MaxSqrtRatio, constants.One)
		}
	}

	if zeroForOne {
		if sqrtPriceLimitX64.Cmp(utils.MinSqrtRatio) < 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX64TooLow
		}
		if sqrtPriceLimitX64.Cmp(p.SqrtRatioX64) >= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX64TooHigh
		}
	} else {
		if sqrtPriceLimitX64.Cmp(utils.MaxSqrtRatio) > 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX64TooHigh
		}
		if sqrtPriceLimitX64.Cmp(p.SqrtRatioX64) <= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX64TooLow
		}
	}

	exactInput := amountSpecified.Cmp(constants.Zero) >= 0

	// keep track of swap state

	state := struct {
		amountSpecifiedRemaining *big.Int
		amountCalculated         *big.Int
		sqrtPriceX64             *big.Int
		tick                     int
		liquidity                *big.Int
	}{
		amountSpecifiedRemaining: amountSpecified,
		amountCalculated:         constants.Zero,
		sqrtPriceX64:             p.SqrtRatioX64,
		tick:                     p.TickCurrent,
		liquidity:                p.Liquidity,
	}

	// start swap while loop
	for state.amountSpecifiedRemaining.Cmp(constants.Zero) != 0 && state.sqrtPriceX64.Cmp(sqrtPriceLimitX64) != 0 {
		var step StepComputations
		step.sqrtPriceStartX64 = state.sqrtPriceX64

		// because each iteration of the while loop rounds, we can't optimize this code (relative to the smart contract)
		// by simply traversing to the next available tick, we instead need to exactly replicate
		// tickBitmap.nextInitializedTickWithinOneWord
		step.tickNext, step.initialized = p.TickDataProvider.NextInitializedTickWithinOneWord(state.tick, zeroForOne, p.TickSpacing)

		if step.tickNext < utils.MinTick {
			step.tickNext = utils.MinTick
		} else if step.tickNext > utils.MaxTick {
			step.tickNext = utils.MaxTick
		}

		step.sqrtPriceNextX64, err = utils.GetSqrtRatioAtTick(step.tickNext)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		var targetValue *big.Int
		if zeroForOne {
			if step.sqrtPriceNextX64.Cmp(sqrtPriceLimitX64) < 0 {
				targetValue = sqrtPriceLimitX64
			} else {
				targetValue = step.sqrtPriceNextX64
			}
		} else {
			if step.sqrtPriceNextX64.Cmp(sqrtPriceLimitX64) > 0 {
				targetValue = sqrtPriceLimitX64
			} else {
				targetValue = step.sqrtPriceNextX64
			}
		}

		state.sqrtPriceX64, step.amountIn, step.amountOut, step.feeAmount, err = utils.ComputeSwapStep(state.sqrtPriceX64, targetValue, state.liquidity, state.amountSpecifiedRemaining, p.Fee)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		if exactInput {
			state.amountSpecifiedRemaining = new(big.Int).Sub(state.amountSpecifiedRemaining, new(big.Int).Add(step.amountIn, step.feeAmount))
			state.amountCalculated = new(big.Int).Sub(state.amountCalculated, step.amountOut)
		} else {
			state.amountSpecifiedRemaining = new(big.Int).Add(state.amountSpecifiedRemaining, step.amountOut)
			state.amountCalculated = new(big.Int).Add(state.amountCalculated, new(big.Int).Add(step.amountIn, step.feeAmount))
		}

		// TODO
		if state.sqrtPriceX64.Cmp(step.sqrtPriceNextX64) == 0 {
			// if the tick is initialized, run the tick transition
			if step.initialized {
				liquidityNet := p.TickDataProvider.GetTick(step.tickNext).LiquidityNet
				// if we're moving leftward, we interpret liquidityNet as the opposite sign
				// safe because liquidityNet cannot be type(int128).min
				if zeroForOne {
					liquidityNet = new(big.Int).Mul(liquidityNet, constants.NegativeOne)
				}
				state.liquidity = utils.AddDelta(state.liquidity, liquidityNet)
			}
			if zeroForOne {
				state.tick = step.tickNext - 1
			} else {
				state.tick = step.tickNext
			}
		} else if state.sqrtPriceX64.Cmp(step.sqrtPriceStartX64) != 0 {
			// recompute unless we're on a lower tick boundary (i.e. already transitioned ticks), and haven't moved
			state.tick, err = utils.GetTickAtSqrtRatio(state.sqrtPriceX64)
			if err != nil {
				return nil, nil, nil, 0, err
			}
		}
	}
	return state.amountCalculated, state.sqrtPriceX64, state.liquidity, state.tick, nil
}
