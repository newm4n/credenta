package credenta

import (
	"errors"
	"math"
)

// toUint64ByBit return the number of UInt64 array element to store a speciffic role flag and bitno is the n-th
// bit in each of the role bits.
//
// for example :
// uintSeq 0 means the first element in the array of roles.
// uintSeq 1 means the second element in the array of roles.
//
// bitno is the n-th bit in each of the role bits.
// for example :
// bitno 1 means the 1st bit. 00000001.
// bitno 2 means the 2nd bit. 00000010
// bitno 5 means the 5th bit. 00010000
// etc
func toUint64ByBit(roleID int) (uintSeq int, bitno int) {
	left := math.Mod(float64(roleID), 64.0)
	div := roleID / 64
	return div, int(left)
}

// isBitFlagOn return a boolean to check if a specific bit in an uint1 is on
// for example : if currentBit = 5 ... the bit is 00001001
// Means that : bit 0 0 0 0 1 0 0 1
// role ID          7 6 5 4 3 2 1 0
// so for currentBit 5 (00001001) only bit 0 and 3 is on.
// sBitFlagOn(5, 0) =? true
// sBitFlagOn(5, 1) =? false
// etc
func isBitFlagOn(currentBit uint64, roleID int) bool {
	flipper := uint64(1) << roleID
	return flipper == currentBit&flipper
}

// setBitFlagOn return a new uint64 on which the n-th bit (roleID) of the current uint64 is set set to 1
func setBitFlagOn(currentBit uint64, roleID int) uint64 {
	flipper := uint64(1) << roleID
	return currentBit | flipper
}

// setBitFlagOff return a new uint64 on which the n-th bit (roleID) of the current uint64 is set set to 0
func setBitFlagOff(currentBit uint64, roleID int) uint64 {
	flipper := uint64(1) << roleID
	notFlipper := 0xFFFFFFFF ^ flipper
	return currentBit & notFlipper
}

// IsHaveRole check if a specific role setting (roles), have a specific role ID (roleID) to ON
func IsHaveRole(roles []uint64, roleID int) bool {
	roleNo, bitNo := toUint64ByBit(roleID)
	return isBitFlagOn(roles[roleNo], bitNo)
}

// AddRole produce a new role setting on which on the new role will contain a specific role id.
func AddRole(roles []uint64, roleID int) ([]uint64, error) {
	roleNo, bitNo := toUint64ByBit(roleID)
	if roleNo > len(roles) {
		return nil, errors.New("role number out of bounds. role not large enough")
	}
	ret := roles
	ret[roleNo] = setBitFlagOn(ret[roleNo], bitNo)
	return ret, nil
}

// RemoveRole produce a new role setting on which on the new role will remove a specific role id.
func RemoveRole(roles []uint64, roleID int) ([]uint64, error) {
	roleNo, bitNo := toUint64ByBit(roleID)
	if roleNo > len(roles) {
		return nil, errors.New("role number out of bounds. role not large enough")
	}
	ret := roles
	ret[roleNo] = setBitFlagOff(ret[roleNo], bitNo)
	return ret, nil
}
