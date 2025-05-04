package credenta

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommon_toUint64ByBit(t *testing.T) {
	seq, bitno := toUint64ByBit(0)
	assert.Equal(t, 0, seq)
	assert.Equal(t, 0, bitno)

	seq, bitno = toUint64ByBit(1)
	assert.Equal(t, 0, seq)
	assert.Equal(t, 1, bitno)

	seq, bitno = toUint64ByBit(64)
	assert.Equal(t, 1, seq)
	assert.Equal(t, 0, bitno)
}
func TestCommon_isBitFlagOn(t *testing.T) {
	assert.False(t, isBitFlagOn(0x00000000, 0))
	assert.True(t, isBitFlagOn(0x00000001, 0))

	assert.True(t, isBitFlagOn(0x00000003, 0))
	assert.True(t, isBitFlagOn(0x00000003, 1))

	assert.True(t, isBitFlagOn(0x00000005, 0))
	assert.False(t, isBitFlagOn(0x00000005, 1))
	assert.True(t, isBitFlagOn(0x00000005, 2))

	assert.True(t, isBitFlagOn(0x00000009, 0))
	assert.False(t, isBitFlagOn(0x00000009, 1))
	assert.False(t, isBitFlagOn(0x00000009, 2))
	assert.True(t, isBitFlagOn(0x00000009, 3))
}
func TestCommon_setBitFlagOn(t *testing.T) {
	assert.Equal(t, uint64(0x00000001), setBitFlagOn(0x00000000, 0))
	assert.Equal(t, uint64(0x00000002), setBitFlagOn(0x00000000, 1))
	assert.Equal(t, uint64(0x00000004), setBitFlagOn(0x00000000, 2))
	assert.Equal(t, uint64(0x00000008), setBitFlagOn(0x00000000, 3))
	assert.Equal(t, uint64(0x00000010), setBitFlagOn(0x00000000, 4))
}
func TestCommon_setBitFlagOff(t *testing.T) {
	assert.Equal(t, uint64(0x00000000), setBitFlagOff(0x00000001, 0))
	assert.Equal(t, uint64(0x00000000), setBitFlagOff(0x00000002, 1))
	assert.Equal(t, uint64(0x00000000), setBitFlagOff(0x00000004, 2))
	assert.Equal(t, uint64(0x00000000), setBitFlagOff(0x00000008, 3))
	assert.Equal(t, uint64(0x00000000), setBitFlagOff(0x00000010, 4))
}
func TestCommon_IsHaveRole(t *testing.T) {
	assert.False(t, IsHaveRole([]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0))
	assert.True(t, IsHaveRole([]uint64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0))

	assert.True(t, IsHaveRole([]uint64{0x00000003, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0))
	assert.True(t, IsHaveRole([]uint64{0x00000003, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1))

	assert.True(t, IsHaveRole([]uint64{0x00000005, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0))
	assert.False(t, IsHaveRole([]uint64{0x00000005, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1))
	assert.True(t, IsHaveRole([]uint64{0x00000005, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 2))

	assert.True(t, IsHaveRole([]uint64{0x00000009, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0))
	assert.False(t, IsHaveRole([]uint64{0x00000009, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1))
	assert.False(t, IsHaveRole([]uint64{0x00000009, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 2))
	assert.True(t, IsHaveRole([]uint64{0x00000009, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 3))
}
func TestCommon_AddRole(t *testing.T) {
	nRole, err := AddRole([]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 0)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)

	nRole, err = AddRole([]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 1)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{2, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)

	nRole, err = AddRole([]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 2)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{4, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)
}
func TestCommon_RemoveRole(t *testing.T) {
	nRole, err := RemoveRole([]uint64{1, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 0)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)

	nRole, err = RemoveRole([]uint64{2, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 1)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)

	nRole, err = RemoveRole([]uint64{4, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 2)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nRole)
}
