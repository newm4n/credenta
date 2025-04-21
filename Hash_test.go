package credenta

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeVerification(t *testing.T) {
	pass := "the big brown fox jumps over the lazy dog"

	hash, err := MakeVerification(VerificationMethodPLAIN, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodPLAIN, pass, hash))

	hash, err = MakeVerification(VerificationMethodMD5, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodMD5, pass, hash))

	hash, err = MakeVerification(VerificationMethodSHA1, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodSHA1, pass, hash))

	hash, err = MakeVerification(VerificationMethodSHA256, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodSHA256, pass, hash))

	hash, err = MakeVerification(VerificationMethodSHA512, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodSHA512, pass, hash))

	hash, err = MakeVerification(VerificationMethodARGON, pass)
	assert.NoError(t, err)
	fmt.Println(hash)
	assert.True(t, MatchVerification(VerificationMethodARGON, pass, hash))
}
