package credenta

import (
	"fmt"
	"strings"
)

type PassphrasePolicy struct {
	WordCount               int  `json:"wordCount"`
	LetterCountPerWord      int  `json:"letterCountPerWord"`
	LetterCountMinimumTotal int  `json:"letterCountMinimumTotal"`
	MustHaveUpperAlphabet   bool `json:"mustHaveUpperAlphabet"`
	MustHaveNumeric         bool `json:"mustHaveNumeric"`
	MustHaveSymbol          bool `json:"mustHaveSymbol"`
}

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
