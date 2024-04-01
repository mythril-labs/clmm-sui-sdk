package entities

import (
	"errors"
	"strings"
)

var (
	ErrDifferentChain = errors.New("different chain")
	ErrSameAddress    = errors.New("same address")
)

// Currency is any fungible financial instrument, including Ether, all ERC20 tokens, and other chain-native currencies
type Currency interface {
	IsNative() bool
	IsToken() bool
	ChainId() uint
	Decimals() uint
	Symbol() string
	Name() string
	Equal(other Currency) bool
	Wrapped() *Token
}

// baseCurrency is an abstract struct, do not use it directly
type baseCurrency struct {
	currency Currency
	isNative bool   // Returns whether the currency is native to the chain and must be wrapped (e.g. Ether)
	isToken  bool   // Returns whether the currency is a token that is usable in Uniswap without wrapping
	chainId  uint   // The chain ID on which this currency resides
	decimals uint   // The decimals used in representing currency amounts
	symbol   string // The symbol of the currency, i.e. a short textual non-unique identifier
	name     string // The name of the currency, i.e. a descriptive textual non-unique identifier
}

func (c *baseCurrency) IsNative() bool {
	return c.isNative
}

func (c *baseCurrency) IsToken() bool {
	return c.isToken
}

func (c *baseCurrency) ChainId() uint {
	return c.chainId
}

func (c *baseCurrency) Decimals() uint {
	return c.decimals
}

func (c *baseCurrency) Symbol() string {
	return c.symbol
}

func (c *baseCurrency) Name() string {
	return c.name
}

// Equal returns whether the currency is equal to the other currency
func (c *baseCurrency) Equal(other Currency) bool {
	panic("Equal method has to be overridden")
}

func (c *baseCurrency) Wrapped() *Token {
	panic("Wrapped method has to be overridden")
}

// Token represents an ERC20 token with a unique address and some metadata.
type Token struct {
	*baseCurrency
	Address string // The contract address on the chain on which this token lives
}

// NewToken creates a new token with the given currency and address.
func NewToken(chainID uint, address string, decimals uint, symbol string, name string) *Token {
	if decimals >= 255 {
		panic("Token currency decimals must be less than 255")
	}
	token := &Token{
		baseCurrency: &baseCurrency{
			isNative: false,
			isToken:  true,
			chainId:  chainID,
			decimals: decimals,
			symbol:   symbol,
			name:     name,
		},
		Address: address,
	}
	token.baseCurrency.currency = token
	return token
}

// Equal
/**
 * Returns true if the two tokens are equivalent, i.e. have the same chainId and address.
 * @param other token to compare
 */
func (t *Token) Equal(other Currency) bool {
	if other != nil {
		v, isToken := other.(*Token)
		if isToken {
			return v.isToken && t.chainId == v.chainId && t.Address == v.Address
		}
	}
	return false
}

// SortsBefore
/**
 * Returns true if the address of this token sorts before the address of the other token
 * @param other other token to compare
 * @throws if the tokens have the same address
 * @throws if the tokens are on different chains
 */
func (t *Token) SortsBefore(other *Token) (bool, error) {
	if t.chainId != other.chainId {
		return false, ErrDifferentChain
	}
	if t.Address == other.Address {
		return false, ErrSameAddress
	}
	return strings.ToLower(t.Address) < strings.ToLower(other.Address), nil
}

func (t *Token) Wrapped() *Token {
	return t
}
