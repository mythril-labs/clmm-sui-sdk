package utils

import (
	"errors"
	"math/big"

	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
)

const (
	MinTick = -443636  // The minimum tick that can be used on any pool.
	MaxTick = -MinTick // The maximum tick that can be used on any pool.
)

var (
	Q32 = big.NewInt(1 << 32)
	// The sqrt ratio corresponding to the minimum tick that could be used on any pool.
	MinSqrtRatio = big.NewInt(4295048016)
	// The sqrt ratio corresponding to the maximum tick that could be used on any pool.
	MaxSqrtRatio, _ = new(big.Int).SetString("79226673515401279992447579055", 10)

	Mask256 = new(big.Int).Sub(entities.MaxUint256, constants.One)
)

var (
	ErrInvalidTick      = errors.New("invalid tick")
	ErrInvalidSqrtRatio = errors.New("invalid sqrt ratio")
)

func signedShiftRight(val *big.Int, mulBy *big.Int, shiftBy uint) *big.Int {
	val.Mul(val, mulBy)
	val.Rsh(val, shiftBy)
	return val
}

/**
 * Returns the sqrt ratio as a Q64.64 for the given tick. The sqrt ratio is computed as sqrt(1.0001)^tick
 * @param tick the tick for which to compute the sqrt ratio
 */
func GetSqrtRatioAtTick(tick int) (*big.Int, error) {
	if tick < MinTick || tick > MaxTick {
		return nil, ErrInvalidTick
	}
	if tick > 0 {
		return getSqrtRatioAtTickPositive(tick)
	}
	return getSqrtRatioAtTickNegative(tick)
}

var (
	sqrtPositive1, _  = new(big.Int).SetString("79232123823359799118286999567", 10)
	sqrtPositive2, _  = new(big.Int).SetString("79228162514264337593543950336", 10)
	sqrtPositive3, _  = new(big.Int).SetString("79236085330515764027303304731", 10)
	sqrtPositive4, _  = new(big.Int).SetString("79244008939048815603706035061", 10)
	sqrtPositive5, _  = new(big.Int).SetString("79259858533276714757314932305", 10)
	sqrtPositive6, _  = new(big.Int).SetString("79291567232598584799939703904", 10)
	sqrtPositive7, _  = new(big.Int).SetString("79355022692464371645785046466", 10)
	sqrtPositive8, _  = new(big.Int).SetString("79482085999252804386437311141", 10)
	sqrtPositive9, _  = new(big.Int).SetString("79736823300114093921829183326", 10)
	sqrtPositive10, _ = new(big.Int).SetString("80248749790819932309965073892", 10)
	sqrtPositive11, _ = new(big.Int).SetString("81282483887344747381513967011", 10)
	sqrtPositive12, _ = new(big.Int).SetString("83390072131320151908154831281", 10)
	sqrtPositive13, _ = new(big.Int).SetString("87770609709833776024991924138", 10)
	sqrtPositive14, _ = new(big.Int).SetString("87770609709833776024991924138", 10)
	sqrtPositive15, _ = new(big.Int).SetString("119332217159966728226237229890", 10)
	sqrtPositive16, _ = new(big.Int).SetString("179736315981702064433883588727", 10)
	sqrtPositive17, _ = new(big.Int).SetString("407748233172238350107850275304", 10)
	sqrtPositive18, _ = new(big.Int).SetString("2098478828474011932436660412517", 10)
	sqrtPositive19, _ = new(big.Int).SetString("55581415166113811149459800483533", 10)
	sqrtPositive20, _ = new(big.Int).SetString("38992368544603139932233054999993551", 10)
)

func getSqrtRatioAtTickPositive(tick int) (*big.Int, error) {
	var ratio *big.Int
	if tick&1 != 0 {
		ratio = new(big.Int).Set(sqrtPositive1)
	} else {
		ratio = new(big.Int).Set(sqrtPositive2)
	}

	if (tick & 2) != 0 {
		signedShiftRight(ratio, sqrtPositive3, 96)
	}
	if (tick & 4) != 0 {
		signedShiftRight(ratio, sqrtPositive4, 96)
	}
	if (tick & 8) != 0 {
		signedShiftRight(ratio, sqrtPositive5, 96)
	}
	if (tick & 16) != 0 {
		signedShiftRight(ratio, sqrtPositive6, 96)
	}
	if (tick & 32) != 0 {
		signedShiftRight(ratio, sqrtPositive7, 96)
	}
	if (tick & 64) != 0 {
		signedShiftRight(ratio, sqrtPositive8, 96)
	}
	if (tick & 128) != 0 {
		signedShiftRight(ratio, sqrtPositive9, 96)
	}
	if (tick & 256) != 0 {
		signedShiftRight(ratio, sqrtPositive10, 96)
	}
	if (tick & 512) != 0 {
		signedShiftRight(ratio, sqrtPositive11, 96)
	}
	if (tick & 1024) != 0 {
		signedShiftRight(ratio, sqrtPositive12, 96)
	}
	if (tick & 2048) != 0 {
		signedShiftRight(ratio, sqrtPositive13, 96)
	}
	if (tick & 4096) != 0 {
		signedShiftRight(ratio, sqrtPositive14, 96)
	}
	if (tick & 8192) != 0 {
		signedShiftRight(ratio, sqrtPositive15, 96)
	}
	if (tick & 16384) != 0 {
		signedShiftRight(ratio, sqrtPositive16, 96)
	}
	if (tick & 32768) != 0 {
		signedShiftRight(ratio, sqrtPositive17, 96)
	}
	if (tick & 65536) != 0 {
		signedShiftRight(ratio, sqrtPositive18, 96)
	}
	if (tick & 131072) != 0 {
		signedShiftRight(ratio, sqrtPositive19, 96)
	}
	if (tick & 262144) != 0 {
		signedShiftRight(ratio, sqrtPositive20, 96)
	}

	ratio.Rsh(ratio, 32)
	ratio.And(ratio, Mask256)

	return ratio, nil
}

var (
	sqrtNegative1, _  = new(big.Int).SetString("18445821805675392311", 10)
	sqrtNegative2, _  = new(big.Int).SetString("18446744073709551616", 10)
	sqrtNegative3, _  = new(big.Int).SetString("18444899583751176498", 10)
	sqrtNegative4, _  = new(big.Int).SetString("18443055278223354162", 10)
	sqrtNegative5, _  = new(big.Int).SetString("18439367220385604838", 10)
	sqrtNegative6, _  = new(big.Int).SetString("18431993317065449817", 10)
	sqrtNegative7, _  = new(big.Int).SetString("18417254355718160513", 10)
	sqrtNegative8, _  = new(big.Int).SetString("18387811781193591352", 10)
	sqrtNegative9, _  = new(big.Int).SetString("18329067761203520168", 10)
	sqrtNegative10, _ = new(big.Int).SetString("18212142134806087854", 10)
	sqrtNegative11, _ = new(big.Int).SetString("17980523815641551639", 10)
	sqrtNegative12, _ = new(big.Int).SetString("17526086738831147013", 10)
	sqrtNegative13, _ = new(big.Int).SetString("16651378430235024244", 10)
	sqrtNegative14, _ = new(big.Int).SetString("15030750278693429944", 10)
	sqrtNegative15, _ = new(big.Int).SetString("12247334978882834399", 10)
	sqrtNegative16, _ = new(big.Int).SetString("8131365268884726200", 10)
	sqrtNegative17, _ = new(big.Int).SetString("3584323654723342297", 10)
	sqrtNegative18, _ = new(big.Int).SetString("696457651847595233", 10)
	sqrtNegative19, _ = new(big.Int).SetString("26294789957452057", 10)
	sqrtNegative20, _ = new(big.Int).SetString("37481735321082", 10)
)

func getSqrtRatioAtTickNegative(tick int) (*big.Int, error) {
	tick = -tick
	var ratio *big.Int
	if tick&1 != 0 {
		ratio = new(big.Int).Set(sqrtNegative1)
	} else {
		ratio = new(big.Int).Set(sqrtNegative2)
	}

	if (tick & 2) != 0 {
		signedShiftRight(ratio, sqrtNegative3, 64)
	}
	if (tick & 4) != 0 {
		signedShiftRight(ratio, sqrtNegative4, 64)
	}
	if (tick & 8) != 0 {
		signedShiftRight(ratio, sqrtNegative5, 64)
	}
	if (tick & 16) != 0 {
		signedShiftRight(ratio, sqrtNegative6, 64)
	}
	if (tick & 32) != 0 {
		signedShiftRight(ratio, sqrtNegative7, 64)
	}
	if (tick & 64) != 0 {
		signedShiftRight(ratio, sqrtNegative8, 64)
	}
	if (tick & 128) != 0 {
		signedShiftRight(ratio, sqrtNegative9, 64)
	}
	if (tick & 256) != 0 {
		signedShiftRight(ratio, sqrtNegative10, 64)
	}
	if (tick & 512) != 0 {
		signedShiftRight(ratio, sqrtNegative11, 64)
	}
	if (tick & 1024) != 0 {
		signedShiftRight(ratio, sqrtNegative12, 64)
	}
	if (tick & 2048) != 0 {
		signedShiftRight(ratio, sqrtNegative13, 64)
	}
	if (tick & 4096) != 0 {
		signedShiftRight(ratio, sqrtNegative14, 64)
	}
	if (tick & 8192) != 0 {
		signedShiftRight(ratio, sqrtNegative15, 64)
	}
	if (tick & 16384) != 0 {
		signedShiftRight(ratio, sqrtNegative16, 64)
	}
	if (tick & 32768) != 0 {
		signedShiftRight(ratio, sqrtNegative17, 64)
	}
	if (tick & 65536) != 0 {
		signedShiftRight(ratio, sqrtNegative18, 64)
	}
	if (tick & 131072) != 0 {
		signedShiftRight(ratio, sqrtNegative19, 64)
	}
	if (tick & 262144) != 0 {
		signedShiftRight(ratio, sqrtNegative20, 64)
	}

	return ratio, nil
}

var (
	magicSqrt10001, _ = new(big.Int).SetString("255738958999603826347141", 10)
	magicTickLow, _   = new(big.Int).SetString("3402992956809132418596140100660247210", 10)
	magicTickHigh, _  = new(big.Int).SetString("291339464771989622907027621153398088495", 10)
)

/**
 * Returns the tick corresponding to a given sqrt ratio, s.t. #getSqrtRatioAtTick(tick) <= sqrtRatioX64
 * and #getSqrtRatioAtTick(tick + 1) > sqrtRatioX64
 * @param sqrtRatioX64 the sqrt ratio as a Q64.64 for which to compute the tick
 */
func GetTickAtSqrtRatio(sqrtRatioX64 *big.Int) (int, error) {
	if sqrtRatioX64.Cmp(MinSqrtRatio) < 0 || sqrtRatioX64.Cmp(MaxSqrtRatio) >= 0 {
		return 0, ErrInvalidSqrtRatio
	}
	sqrtRatioX128 := new(big.Int).Lsh(sqrtRatioX64, 32)
	msb, err := MostSignificantBit(sqrtRatioX128)
	if err != nil {
		return 0, err
	}
	var r *big.Int
	if big.NewInt(msb).Cmp(big.NewInt(128)) >= 0 {
		r = new(big.Int).Rsh(sqrtRatioX128, uint(msb-127))
	} else {
		r = new(big.Int).Lsh(sqrtRatioX128, uint(127-msb))
	}

	log2 := new(big.Int).Lsh(new(big.Int).Sub(big.NewInt(msb), big.NewInt(128)), 64)

	for i := 0; i < 14; i++ {
		r = new(big.Int).Rsh(new(big.Int).Mul(r, r), 127)
		f := new(big.Int).Rsh(r, 128)
		log2 = new(big.Int).Or(log2, new(big.Int).Lsh(f, uint(63-i)))
		r = new(big.Int).Rsh(r, uint(f.Int64()))
	}

	logSqrt10001 := new(big.Int).Mul(log2, magicSqrt10001)

	tickLow := new(big.Int).Rsh(new(big.Int).Sub(logSqrt10001, magicTickLow), 128).Int64()
	tickHigh := new(big.Int).Rsh(new(big.Int).Add(logSqrt10001, magicTickHigh), 128).Int64()

	if tickLow == tickHigh {
		return int(tickLow), nil
	}

	sqrtRatio, err := GetSqrtRatioAtTick(int(tickHigh))
	if err != nil {
		return 0, err
	}
	if sqrtRatio.Cmp(sqrtRatioX64) <= 0 {
		return int(tickHigh), nil
	} else {
		return int(tickLow), nil
	}
}
