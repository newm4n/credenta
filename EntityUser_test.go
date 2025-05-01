package credenta

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToUint64ByBit(t *testing.T) {
	seq, bit := toUint64ByBit(0)
	assert.Equal(t, seq, 0)
	assert.Equal(t, bit, 0)

	seq, bit = toUint64ByBit(1)
	assert.Equal(t, seq, 0)
	assert.Equal(t, bit, 1)

	seq, bit = toUint64ByBit(63)
	assert.Equal(t, seq, 0)
	assert.Equal(t, bit, 63)

	seq, bit = toUint64ByBit(64)
	assert.Equal(t, seq, 1)
	assert.Equal(t, bit, 0)
}

func TestPrintHex(t *testing.T) {
	fmt.Printf("%x", uint64(0xFFFFFFFF))
}

func TestAddRole(t *testing.T) {
	for i := 0; i < 17; i++ {
		nInt := setBitFlagOn(0, i)
		fmt.Printf("%d -> %x .. %v %v\n", i, nInt, isBitFlagOn(nInt, i), isBitFlagOn(nInt, i+1))
		zint := setBitFlagOff(nInt, i)
		assert.Equal(t, uint64(0), zint)
	}
}

func TestBitFlagOn(t *testing.T) {
	assert.True(t, isBitFlagOn(1, 0))
	assert.True(t, isBitFlagOn(7, 0))
	assert.True(t, isBitFlagOn(7, 1))
	assert.True(t, isBitFlagOn(7, 2))
}

func TestBitShift(t *testing.T) {
	flipper := uint64(1) << 0
	t.Log(flipper)
	assert.True(t, isBitFlagOn(1, 0))
	assert.True(t, isBitFlagOn(2, 1))
	assert.True(t, isBitFlagOn(3, 1))

	flipper = uint64(1) << 1
	t.Log(flipper)

	assert.True(t, isBitFlagOn(7, 0))
	assert.True(t, isBitFlagOn(7, 1))
	assert.True(t, isBitFlagOn(7, 2))
}
