package credenta

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
)

const (
	// VerificationMethodPLAIN specifies if a password will be saved as it is, no change. Password will be seen and
	// recovered by the reader.
	VerificationMethodPLAIN VerificationMethod = "PLAIN"
	// VerificationMethodMD5 specifies if a password will be hashed using MD5 and this hash will be stored.
	VerificationMethodMD5 VerificationMethod = "MD5"
	// VerificationMethodSHA1 specifies if a password will be hashed using SHA1 and this hash will be stored.
	VerificationMethodSHA1 VerificationMethod = "SHA1"
	// VerificationMethodSHA256 specifies if a password will be hashed using SHA256 and this hash will be stored.
	VerificationMethodSHA256 VerificationMethod = "SHA256"
	// VerificationMethodSHA512 specifies if a password will be hashed using SHA512 and this hash will be stored.
	VerificationMethodSHA512 VerificationMethod = "SHA512"
	// VerificationMethodARGON specifies if a password will be hashed using ARGON and this hash will be stored.
	// ARGON is the current award wining hashing mechanism and once saved the password is NOT RECOVERABLE
	VerificationMethodARGON VerificationMethod = "ARGON"
)

// VerificationMethod specify on how a user's password were stored.
type VerificationMethod string

// MakeVerification will hash the supplied pass argument using the hashing mechanism.
func MakeVerification(method VerificationMethod, pass string) (string, error) {
	switch method {
	case VerificationMethodPLAIN:
		return makePlain(pass)
	case VerificationMethodMD5:
		return makeMD5(pass)
	case VerificationMethodSHA1:
		return makeSHA1(pass)
	case VerificationMethodSHA256:
		return makeSHA256(pass)
	case VerificationMethodSHA512:
		return makeSHA512(pass)
	case VerificationMethodARGON:
		return makeARGON(pass)
	default:
		return "", errors.New("unknown verification method")
	}
}

// MatchVerification will return true if the hash of password match to the hashed password, depends on the hasing method
// used when creating the password hash on the first place (MakeVerification)
func MatchVerification(method VerificationMethod, pass, hash string) bool {
	switch method {
	case VerificationMethodPLAIN:
		return matchPLAIN(pass, hash)
	case VerificationMethodMD5:
		return matchMD5(pass, hash)
	case VerificationMethodSHA1:
		return matchSHA1(pass, hash)
	case VerificationMethodSHA256:
		return matchSHA256(pass, hash)
	case VerificationMethodSHA512:
		return matchSHA512(pass, hash)
	case VerificationMethodARGON:
		return matchARGON(pass, hash)
	default:
		return false
	}
}

func makePlain(pass string) (string, error) {
	if len(pass) == 0 {
		return "", fmt.Errorf("password too short")
	}
	return pass, nil
}

func matchPLAIN(pass, hash string) bool {
	return pass == hash
}

func makeMD5(pass string) (string, error) {
	if len(pass) == 0 {
		return "", fmt.Errorf("password too short")
	}
	hasher := md5.New()
	result := hasher.Sum([]byte(pass))
	return hex.EncodeToString(result), nil
}

func matchMD5(pass, hash string) bool {
	hashed, err := makeMD5(pass)
	if err != nil {
		return false
	}
	return hashed == hash
}

func makeSHA1(pass string) (string, error) {
	if len(pass) == 0 {
		return "", fmt.Errorf("password too short")
	}
	hasher := sha1.New()
	result := hasher.Sum([]byte(pass))
	return hex.EncodeToString(result), nil
}

func matchSHA1(pass, hash string) bool {
	hashed, err := makeSHA1(pass)
	if err != nil {
		return false
	}
	return hashed == hash
}

func makeSHA256(pass string) (string, error) {
	if len(pass) == 0 {
		return "", fmt.Errorf("password too short")
	}
	hasher := sha256.New()
	result := hasher.Sum([]byte(pass))
	return hex.EncodeToString(result), nil
}

func matchSHA256(pass, hash string) bool {
	hashed, err := makeSHA256(pass)
	if err != nil {
		return false
	}
	return hashed == hash
}

func makeSHA512(pass string) (string, error) {
	if len(pass) == 0 {
		return "", fmt.Errorf("password too short")
	}
	hasher := sha512.New()
	result := hasher.Sum([]byte(pass))
	return hex.EncodeToString(result), nil
}

func matchSHA512(pass, hash string) bool {
	hashed, err := makeSHA512(pass)
	if err != nil {
		return false
	}
	return hashed == hash
}

func makeARGON(pass string) (string, error) {
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func matchARGON(pass, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(pass, hash)
	if err != nil {
		return false
	}
	return match
}
