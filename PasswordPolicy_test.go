package credenta

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassphrasePolicy_SimpleIsPasswordValid(t *testing.T) {
	valid, err := SimplePasswordPolicy().IsPasswordValid("Testing1pa$$word")
	assert.NoError(t, err)
	assert.True(t, valid)
}
