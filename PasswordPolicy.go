package credenta

import (
	"fmt"
	"strings"
)

// SimplePasswordPolicy will create a new passphrase validation policy.
// The supplied password must be be 1 word thus it must NOT be separated by space
// and the total password must be more than 8 characters.
// Any alphanumeric, symbol or numbers may present or not.
func SimplePasswordPolicy() *PassphrasePolicy {
	return &PassphrasePolicy{
		WordCount:               1,
		LetterCountPerWord:      8,
		LetterCountMinimumTotal: 8,
		MustHaveUpperAlphabet:   false,
		MustHaveNumeric:         false,
		MustHaveSymbol:          false,
	}
}

// StrongPasswordPolicy will create a new passphrase validation policy.
// the supplied passphrase must be minimum of 3 words where each word must be separated by a space.
// Each of the word must be minimum 5 letters and the overall passphrase is
// have a minimum 12 letters, including space.
// Any alphanumeric, symbol or numbers may present or not.
func StrongPasswordPolicy() *PassphrasePolicy {
	return &PassphrasePolicy{
		WordCount:               3,
		LetterCountPerWord:      5,
		LetterCountMinimumTotal: 12,
		MustHaveUpperAlphabet:   false,
		MustHaveNumeric:         false,
		MustHaveSymbol:          false,
	}
}

// ClassicPasswordPolicy will create a new passphrase validation policy.
// The supplied password must be be 1 word thus it must NOT be separated by space
// and the total password must be more than 8 characters.
// The passprase must contains minimal 1 numbers, 1 uppler case letter and 1 symbol.
func ClassicPasswordPolicy() *PassphrasePolicy {
	return &PassphrasePolicy{
		WordCount:               1,
		LetterCountPerWord:      8,
		LetterCountMinimumTotal: 8,
		MustHaveUpperAlphabet:   true,
		MustHaveNumeric:         true,
		MustHaveSymbol:          true,
	}
}

// PassphrasePolicy is the rule of passphrase.
type PassphrasePolicy struct {
	WordCount               int  `json:"wordCount"`
	LetterCountPerWord      int  `json:"letterCountPerWord"`
	LetterCountMinimumTotal int  `json:"letterCountMinimumTotal"`
	MustHaveUpperAlphabet   bool `json:"mustHaveUpperAlphabet"`
	MustHaveNumeric         bool `json:"mustHaveNumeric"`
	MustHaveSymbol          bool `json:"mustHaveSymbol"`
}

// IsPasswordValid test the supplied pass argument if valid according to the rules specified by the Policy.
func (policy *PassphrasePolicy) IsPasswordValid(pass string) (bool, error) {
	realPass := strings.TrimSpace(pass)
	if len(realPass) != len(pass) {
		return false, fmt.Errorf("contain leading and trailing spaces")
	}
	words := strings.Split(realPass, " ")
	if len(words) != policy.WordCount {
		return false, fmt.Errorf("different word count (%d != %d)", len(words), policy.WordCount)
	}
	for _, w := range words {
		if len(w) < policy.LetterCountPerWord {
			return false, fmt.Errorf("need more letter in passphrase word")
		}
	}
	if policy.LetterCountMinimumTotal < len(realPass) {
		return false, fmt.Errorf("passphrase needs minimum %d letters", policy.LetterCountMinimumTotal)
	}
	if policy.MustHaveUpperAlphabet && !strings.ContainsAny(pass, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false, fmt.Errorf("passphrase requires upper alphabet")
	}
	if policy.MustHaveNumeric && !strings.ContainsAny(pass, "0123456789") {
		return false, fmt.Errorf("passphrase requires number")
	}
	if policy.MustHaveSymbol && !strings.ContainsAny(pass, "`'\"\\[]{},./;':!@#$%^&*()_+-=") {
		return false, fmt.Errorf("passphrase requires number")
	}
	return true, nil
}
