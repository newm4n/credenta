package credenta

import "math"

func toUint64ByBit(roleSquence int) (uintSeq int, bitno int) {
	left := math.Mod(float64(roleSquence), 64.0)
	div := roleSquence / 64
	return div, int(left)
}

func isBitFlagOn(currentBit uint64, bitSequence int) bool {
	flipper := uint64(1) << bitSequence
	return flipper == currentBit&flipper
}

func setBitFlagOn(currentBit uint64, bitSequence int) uint64 {
	flipper := uint64(1) << bitSequence
	return currentBit | flipper
}

func setBitFlagOff(currentBit uint64, bitSequence int) uint64 {
	flipper := uint64(1) << bitSequence
	notFlipper := 0xFFFFFFFF ^ flipper
	return currentBit & notFlipper
}
