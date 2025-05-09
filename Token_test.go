package credenta

import (
	"crypto/rsa"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestToken_JWTWithOriginalKey(t *testing.T) {
	var theClaim jws.Claims
	t.Run("Test_MakeClaims", func(t *testing.T) {
		claims := jws.Claims{}
		claims.SetIssuer("TheIssuer")
		claims.SetSubject("TheSubject")
		claims.SetAudience("Audience1", "Audience2")
		claims.SetIssuedAt(time.Now())
		claims.SetNotBefore(time.Now())
		claims.SetExpiration(time.Now().Add(24 * time.Hour))
		assert.NotNil(t, claims)
		theClaim = claims
	})

	var rsaPrivate *rsa.PrivateKey
	var rsaPublic *rsa.PublicKey
	t.Run("Test_GetKey", func(t *testing.T) {
		privbytes, err := os.ReadFile("./data/crypto/exam.priv")
		assert.NoError(t, err)
		rPriv, err := crypto.ParseRSAPrivateKeyFromPEM(privbytes)
		assert.NoError(t, err)
		rsaPrivate = rPriv

		pubbytes, err := os.ReadFile("./data/crypto/exam.pub")
		assert.NoError(t, err)
		rPub, err := crypto.ParseRSAPublicKeyFromPEM(pubbytes)
		assert.NoError(t, err)
		rsaPublic = rPub
	})

	var tokenByte []byte
	t.Run("Test_MakeToken", func(t *testing.T) {
		jwtBytes := jws.NewJWT(theClaim, crypto.SigningMethodRS256)
		byts, err := jwtBytes.Serialize(rsaPrivate)
		assert.NoError(t, err)
		tokenByte = byts
	})

	t.Run("Test_ValidateToken", func(t *testing.T) {
		jwt, err := jws.ParseJWT(tokenByte)
		assert.NoError(t, err)

		err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256)
		assert.NoError(t, err)
	})

}

func TestToken_GenerateTokenPair(t *testing.T) {
	issuer := "TheIssuer"
	subject := "TheSubject"
	audience := []string{"audience1", "audience2"}
	additional := map[string]interface{}{"map1": "value1"}
	issuedAt := time.Now()
	accessTokenAge := time.Minute * 5
	refreshTokenAge := time.Hour * 24 * 30
	privateKey := GetDefaultPrivateKey()
	signMethod := crypto.SigningMethodRS256
	at, rt, err := GenerateNewJWTTokenPair(issuer, subject, audience, additional, issuedAt, accessTokenAge, refreshTokenAge, privateKey, signMethod)
	assert.NoError(t, err)

	at, err = GenerateJWTToken(issuer, subject, audience, AccessTokenType, additional, issuedAt, issuedAt, issuedAt.Add(accessTokenAge), privateKey, signMethod)

	t.Log(at)
	t.Log(rt)

	publicKey := GetDefaultPublicKey()

	aissuer2, subject2, audience2, tokenType2, additional2, err := ReadJWTToken(at, publicKey, signMethod)
	assert.NoError(t, err)

	assert.Equal(t, issuer, aissuer2)
	assert.Equal(t, subject, subject2)
	assert.Equal(t, audience, audience2)
	assert.Equal(t, string(AccessTokenType), string(tokenType2))
	assert.Equal(t, additional, additional2)
}

func TestToken_ExpireTokenRead(t *testing.T) {
	issuer := "TheIssuer"
	subject := "TheSubject"
	audience := []string{"audience1", "audience2"}
	additional := map[string]interface{}{"map1": "value1"}
	issuedAt := time.Now().Add(time.Hour * 24 * 365 * -1)
	accessTokenAge := time.Minute * 5
	privateKey := GetDefaultPrivateKey()
	signMethod := crypto.SigningMethodRS256
	at, err := GenerateJWTToken(issuer, subject, audience, AccessTokenType, additional, issuedAt, issuedAt, issuedAt.Add(accessTokenAge), privateKey, signMethod)
	assert.NoError(t, err)

	publicKey := GetDefaultPublicKey()
	_, _, _, _, _, err = ReadJWTToken(at, publicKey, signMethod)
	assert.Error(t, err)
	log.Println(err.Error())
}
