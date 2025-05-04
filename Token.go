package credenta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

	DefaultCertificatePem = `-----BEGIN CERTIFICATE-----
MIIDszCCApsCFEconmI48Jyia8stpnpUkWY4k1fKMA0GCSqGSIb3DQEBCwUAMIGV
MQswCQYDVQQGEwJpZDEQMA4GA1UECAwHSmFrYXJ0YTEQMA4GA1UEBwwHSmFrYXJ0
YTESMBAGA1UECgwJSHlwZXJqdW1wMRQwEgYDVQQLDAtEZXZlbG9wbWVudDEWMBQG
A1UEAwwNaHlwZXJqdW1wLmRldjEgMB4GCSqGSIb3DQEJARYRZGV2QGh5cGVyanVt
cC5kZXYwHhcNMjIxMDI0MTEyMzUzWhcNMjMxMDI0MTEyMzUzWjCBlTELMAkGA1UE
BhMCaWQxEDAOBgNVBAgMB0pha2FydGExEDAOBgNVBAcMB0pha2FydGExEjAQBgNV
BAoMCUh5cGVyanVtcDEUMBIGA1UECwwLRGV2ZWxvcG1lbnQxFjAUBgNVBAMMDWh5
cGVyanVtcC5kZXYxIDAeBgkqhkiG9w0BCQEWEWRldkBoeXBlcmp1bXAuZGV2MIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxyG4ghFvv3BNarBBIuw3gIbs
zBVxPdtn5gZyCGY2aUhlxuBwRGLU3JFy4z+a1daIgIgsSGmW0CjYE2PAxpWPp8Br
9yU4n9nA9DWgjd6qObgwxag3NV5igvWOg17r4AifvroPnMejUznuJDJqS4o9Qw58
EvNIyaQIzOSrsOGd56HBganVo9/Rn1I1+A0AHASCdW7yTryoSKxcpQ9x/yMGpkZp
hfd1s8bacfB9uB14jA0OMy1te6DcPf6ncvl3LCziMxOnWw/nSHPr3Q0mozUaOezr
VW2iXpPF+4c/8gnQN8+LKxkbY7Z5823VUlyBYy2H/+bq2n5Ztjn4b3RdX4otywID
AQABMA0GCSqGSIb3DQEBCwUAA4IBAQBEpu5DjbmqnGpOhojYnlsN47RxDH473o2x
Kbh5Uv6BNdw7umrbUxW6mJy6RHXa/4rXxWHz9vkuMRTl6f7YEzKIErEr4KOdaYN2
CXmsGDq2pXMP3LOfWqaDxv5X0XQCIDWoF9KAJ7blbpj5twIGcBu+6i13BAixKRd6
K63nIo/a1gVPpdk7Gw1AXRYlrSifRA7z54LkXtEvvd+NaQrA43ROhzp4iocYfWTg
L0RgbIRuNaD08qQGx2y9XjZwPne9k0lYJ/xZUvDeK7kn15D71fNPL98iCczRTJV+
T3ylcQfgHAo2qbf3thIli5557WMeN9tTCad7PVRFxvp5CnUE17F9
-----END CERTIFICATE-----`

	DefaultPrivateKeyPem = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDHIbiCEW+/cE1q
sEEi7DeAhuzMFXE922fmBnIIZjZpSGXG4HBEYtTckXLjP5rV1oiAiCxIaZbQKNgT
Y8DGlY+nwGv3JTif2cD0NaCN3qo5uDDFqDc1XmKC9Y6DXuvgCJ++ug+cx6NTOe4k
MmpLij1DDnwS80jJpAjM5Kuw4Z3nocGBqdWj39GfUjX4DQAcBIJ1bvJOvKhIrFyl
D3H/IwamRmmF93Wzxtpx8H24HXiMDQ4zLW17oNw9/qdy+XcsLOIzE6dbD+dIc+vd
DSajNRo57OtVbaJek8X7hz/yCdA3z4srGRtjtnnzbdVSXIFjLYf/5uraflm2Ofhv
dF1fii3LAgMBAAECggEAV08HJWw8uUylfO4rmQLF8QH3iCGspYjx/I596AXcvHuk
ODcGsk089DqHi3DDVBe8gNJzzAoTuE/4MUQu7UL9GfmJvPJiad/hNJHgw+NJcoY6
SCeRkUJBYhcFNb9kHzXYyJiuqLh1eVOwgMlkBpvlcWXD69EkbFiSDTgDuUmq4Lgq
z3mUnu6cQkVXhE0rtX8tdDXEebLyKv+8YReUxuWWm8490YZqp8M0jGLbLPKIUkOr
nAvVlX9A4sykr54YvIrvVNw1wVrrPIjI0v4vcCqInb/0PfsvAzG5Ud67KjpFugR5
K0qEUL1Wsw9AcHxDCoQvMQQD2nK65p9jXbfDV5Cu+QKBgQDqHjIiBjVGtAIpgBV5
1whgllc1CCtt/0tR7RNAmZKih0en/OhdFBUfzAAF8gX/j+la3GVTXld1/9DzT7tZ
23GCIpoEyUimXc9gF+B0w05lH3xQYWCIoztMj2iRTDD7jmXQzqDmFa6OJgwFLY4e
wQexqNlvxHp8Td4W0Ojef5CtwwKBgQDZvmZZprKwOQ2UVPcy/6Q7CU6Ij31JJciN
zhWqVrIjj9XYFQK/GmInfGxsvvr3IsdM53Lt6ugrn+liS77dzP/sYqOwjUFhIMZK
Z9+YfB+YJD1DDUBZt+r9WACmLzxOZfXE7fe7nMgxyfz2I6ICMZ1DOLFhETCnC5pp
CZd26oXXWQKBgD00SazlbJYgRxRsXLDui00c4I2Hpjrqa9luHgNcYp5EuXHsRx7W
OjOG1Fa5j+Hg0IOlbIPf/QNnLkv9gyAZo1H/E76+lFSR373iYBaGXH9JPOmSm3b9
HWqFbzPU9FU/Q9TTv/KGpoyY27ma0DWwBv/mAXobpl3KyY2zbb2FIeCbAoGAbqvF
nb+KhuMYsdHVqwggUxlR3zr/NNSNcPXUMTXLaSPMTv2u3a7tQKCPA162dDIrFj11
PtPsmW+30YwqQNXXJjCkfjHtjw53eo39KaW88TlKIfB0SqWePJIkElNj1X0hQ6yo
A6WWYygE+J331CGfivEfxvRTxDOzkbucToa47FECgYEAwxpVi+TXzAR0uuIb4rJT
aYQ6tqxRYxwXJwS9oi9q5PBuN68yV2w3hsjRiebGJRwiqRcbqwvSw0C9n3U60I0+
aym09mHZDzLldUlH+lt4NbPjdiqNQXwBnB9MnJRpwKxdj97FY7nB2E2GgCHlkleN
GMoEuXHErktY7j5JAOYyAbo=
-----END PRIVATE KEY-----`
)

type TokenType string

// DefaultCertificate get the default x509 certificate based on the built in certificate PEM
func DefaultCertificate() *x509.Certificate {
	blk, _ := pem.Decode([]byte(DefaultCertificatePem))
	certificate, _ := x509.ParseCertificate(blk.Bytes)
	return certificate
}

// DefaultPrivateKey get the default RSA Private Key based on the built in PrivateKey PEM
func DefaultPrivateKey() *rsa.PrivateKey {
	blk, _ := pem.Decode([]byte(DefaultPrivateKeyPem))
	privateKey, _ := x509.ParsePKCS8PrivateKey(blk.Bytes)
	return privateKey.(*rsa.PrivateKey)
}

// LoadCertificateFromPEMFile get the x509 certificate  based on the PEM file specified on path at
// certificatePEMFilePath argument. If something goes wrong, it will return the default certificate
// as returned by DefaultCertificate function
func LoadCertificateFromPEMFile(certificatePEMFilePath string) *x509.Certificate {
	if strings.TrimSpace(certificatePEMFilePath) == "" {
		return DefaultCertificate()
	}
	if _, err := os.Stat(certificatePEMFilePath); err != nil {
		return DefaultCertificate()
	}
	certPemContent, err := os.ReadFile(certificatePEMFilePath)
	if err != nil {
		return DefaultCertificate()
	}
	blk, _ := pem.Decode(certPemContent)
	if blk == nil {
		return DefaultCertificate()
	}
	certificate, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return DefaultCertificate()
	}
	return certificate
}

// LoadPrivateKeyFromPEMFile get the RSA PrivateKey  based on the PEM file specified on path at
// privateKeyPEMFilePath argument. If something goes wrong, it will return the default RSA PrivateKey
// as returned by DefaultPrivateKey function
func LoadPrivateKeyFromPEMFile(privateKeyPEMFilePath string) *rsa.PrivateKey {
	if strings.TrimSpace(privateKeyPEMFilePath) == "" {
		return DefaultPrivateKey()
	}
	if _, err := os.Stat(privateKeyPEMFilePath); err != nil {
		return DefaultPrivateKey()
	}
	keyPemContent, err := os.ReadFile(privateKeyPEMFilePath)
	if err != nil {
		return DefaultPrivateKey()
	}
	blk, _ := pem.Decode(keyPemContent)
	if blk == nil {
		return DefaultPrivateKey()
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return DefaultPrivateKey()
	}
	return privateKey.(*rsa.PrivateKey)
}

// RefreshNewAccessToken will generate a new AccessToken based on the  supplied valid refresh token string in refreshToken argument.
// The newly created access token will have an age of accessTokenAge starting of the time when this function is called.
// The supplied RefreshToken will be validated using RSA Private key in privateKey argument and
// the newly created AccessToken will be signed using x509 Certificate suplied in the certificate argument.
func RefreshNewAccessToken(refreshToken string, accessTokenAge time.Duration, certificate *x509.Certificate, privateKey *rsa.PrivateKey, signMethod *crypto.SigningMethodRSA) (accessToken string, err error) {
	issuer, subject, audiences, tokenType, additional, err := ReadJWTToken(refreshToken, certificate, signMethod)
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
func ReadJWTToken(token string, certificate *x509.Certificate, signMethod *crypto.SigningMethodRSA) (issuer, subject string, audiences []string, tokenType TokenType, additional map[string]interface{}, err error) {
	jwt, err := jws.ParseJWT([]byte(token))
	if err != nil {
		return "", "", nil, "", nil, fmt.Errorf("malformed jwt token")
	}

	if err := jwt.Validate(certificate.PublicKey, signMethod); err != nil {
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
