package credenta

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"os"
	"strings"
	"time"
)

const (
	RefreshTokenType TokenType = "rt+JWT"
	AccessTokenType  TokenType = "at+JWT"

	defaultPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKKsvsXtgZ2qrj
/aMnLDhmCCy3FiMnNpgBY//zrbce6ScQC33HxyF39Kx6PsZctxSF+VwCJcpDKtMi
7Tt69lbun9jUuAOH8KuukIHF/99cKB4MnjH+Z2CW59mMVwgvnj2DXez/HOF+bnIA
kw6w5xVECWfGxOWRkhQLmFKpPO2ie3ostPtornYqg9KddageONNi/PsKxccloskf
pYYf5jzg2lLgs7UqwGZAUYECR6XTlhJcHSUTdrU/hwKBur0pfudRt6frkY0KEaeK
E3h9l3/ZXWhWzMixWjYj16ACYYQrKZgktDRDBHRpDrUzbqeQYAw0l2B74pe4gC0a
rrMFwoC9AgMBAAECggEAYcYL5M/D5NEgB+6hiu70gcgfVBa1PqBFKJsD7QaNQzpQ
hY6BMO7qDVk8V0zn42w51UeRi4paRVy/Syt/skrUJUkSdWJfds3bQiwqTyeeDzRp
wAF8PUUi7ijISnrG/zyhFFkHJySBHAPvR77Xgo/n9YU09uk/+8SxxB/Rjn7kkMkZ
67sG4Nmjp1ek8XndWJnvuJ+xo/Shgy4jhyjy0CJaN0xbFXrbflczpvEMx7qUMTdG
YcYLz1OgieBQFWJBVW6ptIYXkgGL+lpMnCHzyjqGWx1f+nO47xoaAcltkbZMTfyV
2gyLTRkOcSQbeKmSYDzPuwCVVwJLtPdcbObJeT9dswKBgQDf5W/cZdMHAXSEWb7l
OOxNOXNlr22q06uoS1Enjbt+6bERgU7ZTR0ia328grYoUcV/+8LLVYepnd37tuwd
W+MB4zpOwhbOMN+mOfJN4HjpjwuI/itIqVtlu8guF00/L6X9imi07KULXrCKOejR
V+mmN/68x+vGOa4IkSzGIH4xewKBgQDnJ8Agch+xE2UvQ7pCkRyPo9Rt9OpTgsfY
OmnD5gP4ZnrXNCSUU9QvwVLysAbMldA0GEeZdhYrupg/R7iWvEGAwpmA4hp1F2SZ
TtEKZttLR4gMFdd6ilbiHKKYKDKkcMdfxAl1TP4Zm2hYA1skaPyai1gVGe7wCCjA
cUyp3pe1JwKBgFCT4BgvxSzGR0rCicMxI0n/nRpBcnSCTUr6IDDd/1aDgChOozPt
XsjeapgHastD8pJG5yoKlBJlMFjA0mUWhrJNNtTVYSO/zx2hySRh3uIfiwU4hBdY
a/5HAJol5LUSzuhagahralKXN23nvXRp8TrS+Ci0wpPKemm25ahAVWo1AoGBAJwc
34fKK5cm1y5tmky8vkJQTfaY8uzFpXxmLuoL3WCUrZ/L6mx2lRZPhVq8AUuIXl3g
i/KbquuLkKkkIglDSSXRx2Qgz+eGjf1wGoPg5XfY7ovi8G0lIvqAhlsmwtUGCdCm
kBC1l+LpbzYJxjM36Gnjc/CEXDel+wfFPRZ4a5L/AoGAXMYQrveqYR5Sj64a7TR0
ZixGKIGFxtuA03oGfhtR79PGAl9mr1Okkg17Htnf0g7w58V6UJVPm+luuzOrK6Eu
BmYelaypxlWfOjaamOVXpE/qr8IHowqCIohAaiSykzbkm5zfjhZijRkNFQ0AoaCI
lKDDtUkmCt+pbV0aucTF/VQ=
-----END PRIVATE KEY-----`

	defaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyirL7F7YGdqq4/2jJyw4
ZggstxYjJzaYAWP/8623HuknEAt9x8chd/Ssej7GXLcUhflcAiXKQyrTIu07evZW
7p/Y1LgDh/CrrpCBxf/fXCgeDJ4x/mdglufZjFcIL549g13s/xzhfm5yAJMOsOcV
RAlnxsTlkZIUC5hSqTztont6LLT7aK52KoPSnXWoHjjTYvz7CsXHJaLJH6WGH+Y8
4NpS4LO1KsBmQFGBAkel05YSXB0lE3a1P4cCgbq9KX7nUben65GNChGnihN4fZd/
2V1oVszIsVo2I9egAmGEKymYJLQ0QwR0aQ61M26nkGAMNJdge+KXuIAtGq6zBcKA
vQIDAQAB
-----END PUBLIC KEY-----`
)

type TokenType string

// GetDefaultPrivateKey get the default RSA Private Key based on the built in PrivateKey
func GetDefaultPrivateKey() *rsa.PrivateKey {
	rPriv, err := crypto.ParseRSAPrivateKeyFromPEM([]byte(defaultPrivateKey))
	if err != nil {
		return nil
	}
	return rPriv
}

// GetDefaultPublicKey get the default RSA Public Key based on the built in PublicKey
func GetDefaultPublicKey() *rsa.PublicKey {
	rPub, err := crypto.ParseRSAPublicKeyFromPEM([]byte(defaultPublicKey))
	if err != nil {
		return nil
	}
	return rPub
}

// LoadPublicKeyFromFile get the RSA PublicKey  based on the PEM file specified on path at
// publicKeyPath argument. If something goes wrong, it will return the default RSA PublicKey
// as returned by GetDefaultPublicKey function
func LoadPublicKeyFromFile(publicKeyPath string) *rsa.PublicKey {
	if strings.TrimSpace(publicKeyPath) == "" {
		return GetDefaultPublicKey()
	}
	if _, err := os.Stat(publicKeyPath); err != nil {
		return GetDefaultPublicKey()
	}
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return GetDefaultPublicKey()
	}
	rPub, err := crypto.ParseRSAPublicKeyFromPEM(publicKeyContent)
	if err != nil {
		return GetDefaultPublicKey()
	}
	return rPub
}

// LoadPrivateKeyFromFile get the RSA PrivateKey  based on the file specified on path at
// privateKeyPEMFilePath argument. If something goes wrong, it will return the default RSA PrivateKey
// as returned by GetDefaultPrivateKey function
func LoadPrivateKeyFromFile(privateKeyFilePath string) *rsa.PrivateKey {
	if strings.TrimSpace(privateKeyFilePath) == "" {
		return GetDefaultPrivateKey()
	}
	if _, err := os.Stat(privateKeyFilePath); err != nil {
		return GetDefaultPrivateKey()
	}
	privateKeyContent, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return GetDefaultPrivateKey()
	}
	rPriv, err := crypto.ParseRSAPrivateKeyFromPEM(privateKeyContent)
	if err != nil {
		return GetDefaultPrivateKey()
	}
	return rPriv
}

// RefreshNewAccessToken will generate a new AccessToken based on the  supplied valid refresh token string in refreshToken argument.
// The newly created access token will have an age of accessTokenAge starting of the time when this function is called.
// The supplied RefreshToken will be validated using RSA Private key in privateKey argument and
// the newly created AccessToken will be signed using x509 Certificate suplied in the certificate argument.
func RefreshNewAccessToken(refreshToken string, accessTokenAge time.Duration, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, signMethod *crypto.SigningMethodRSA) (accessToken string, err error) {
	issuer, subject, audiences, tokenType, additional, err := ReadJWTToken(refreshToken, publicKey, signMethod)
	if err != nil {
		return "", err
	}
	if tokenType != RefreshTokenType {
		return "", errors.New("refresh token does not match")
	}
	at, err := GenerateJWTToken(issuer, subject, audiences, AccessTokenType, additional, time.Now(), time.Now(), time.Now().Add(accessTokenAge), privateKey, signMethod)
	if err != nil {
		return "", err
	}
	return at, nil
}

// ReadJWTToken will read the supplied Token string suplied in the token argument.
// It will return all information pertaining the token. such as issuer, subject, audiences, tokenType, etc.
// If something wrong with the token, e.g. expired token or wrong certificate, it will return an error.
func ReadJWTToken(token string, publicKey *rsa.PublicKey, signMethod *crypto.SigningMethodRSA) (issuer, subject string, audiences []string, tokenType TokenType, additional map[string]interface{}, err error) {
	jwt, err := jws.ParseJWT([]byte(token))
	if err != nil {
		return "", "", nil, "", nil, fmt.Errorf("malformed jwt token")
	}

	if err := jwt.Validate(publicKey, signMethod); err != nil {
		return "", "", nil, "", nil, err
	}

	var ttype TokenType
	claims := jwt.Claims()
	additional = make(map[string]interface{})
	for k, v := range claims {
		kup := strings.ToUpper(k)
		if kup == "TYP" {
			ttypes := v.(string)
			ttype = TokenType(ttypes)
		} else if kup != "ISS" && kup != "AUD" && kup != "SUB" && kup != "IAT" && kup != "EXP" && kup != "NBF" {
			additional[k] = v
		}
	}

	issuer = ""
	if iss, ok := claims.Issuer(); ok {
		issuer = iss
	}
	subject = ""
	if sub, ok := claims.Subject(); ok {
		subject = sub
	}
	audience := make([]string, 0)
	if aud, ok := claims.Audience(); ok {
		audience = aud
	}

	return issuer, subject, audience, ttype, additional, nil
}

// GenerateNewJWTTokenPair will create new pair of JWT Access and Refresh Token strings. these string is ready for the
// web Authorization uses (AccessToken for access REST API, and RefreshToken for refreshing the access token if the
// access is expired).
// Both token will have the same informations such as, issuer, subject, audience, additional, etc.
// But they have different Age.
// Both token will be signed using the same privateKey.
func GenerateNewJWTTokenPair(issuer, subject string, audiences []string, additional map[string]interface{}, issuedAt time.Time, accessTokenAge, refreshTokenAge time.Duration, privateKey *rsa.PrivateKey, signMethod *crypto.SigningMethodRSA) (string, string, error) {
	at, err := GenerateJWTToken(issuer, subject, audiences, AccessTokenType, additional, issuedAt, issuedAt, issuedAt.Add(accessTokenAge), privateKey, signMethod)
	if err != nil {
		return "", "", err
	}
	ar, err := GenerateJWTToken(issuer, subject, audiences, RefreshTokenType, additional, issuedAt, issuedAt, issuedAt.Add(refreshTokenAge), privateKey, signMethod)
	if err != nil {
		return "", "", err
	}
	return at, ar, nil
}

// GenerateJWTToken will create a new JWT Token based on th supplied tokenType argument.
// It will have information pertaining the new Token such as issuer, subject, audiences, additionals, etc.
// The generated token will be signed using the supplied privateKey
func GenerateJWTToken(issuer, subject string, audiences []string, tokenType TokenType, additional map[string]interface{}, issuedAt, notBefore, expireAt time.Time, privateKey *rsa.PrivateKey, signMethod *crypto.SigningMethodRSA) (string, error) {
	claims := jws.Claims{}
	claims.SetIssuer(issuer)
	claims.SetSubject(subject)
	claims.SetAudience(audiences...)
	claims.SetIssuedAt(issuedAt)
	claims.SetNotBefore(notBefore)
	claims.SetExpiration(expireAt)
	claims["typ"] = tokenType
	if additional != nil {
		for k, v := range additional {
			claims[k] = v
		}
	}

	jwtBytes := jws.NewJWT(claims, signMethod)

	tokenByte, err := jwtBytes.Serialize(privateKey)
	if err != nil {
		return "", err
	}
	return string(tokenByte), nil
}
