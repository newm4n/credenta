package credenta

import (
	"github.com/SermoDigital/jose/crypto"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestToken_GenerateTokenPair(t *testing.T) {
	issuer := "TheIssuer"
	subject := "TheSubject"
	audience := []string{"audience1", "audience2"}
	additional := map[string]interface{}{"map1": "value1"}
	issuedAt := time.Now()
	accessTokenAge := time.Minute * 5
	refreshTokenAge := time.Hour * 24 * 30
	privateKey := DefaultPrivateKey()
	signMethod := crypto.SigningMethodRS256
	at, rt, err := GenerateNewJWTTokenPair(issuer, subject, audience, additional, issuedAt, accessTokenAge, refreshTokenAge, privateKey, signMethod)
	assert.NoError(t, err)

	t.Logf("at: %v", at)
	t.Logf("rt: %v", rt)

	certificate := DefaultCertificate()

	aissuer2, subject2, audience2, tokenType2, additional2, err := ReadJWTToken(at, certificate, signMethod)
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
	privateKey := DefaultPrivateKey()
	signMethod := crypto.SigningMethodRS256
	at, err := GenerateJWTToken(issuer, subject, audience, AccessTokenType, additional, issuedAt, issuedAt, issuedAt.Add(accessTokenAge), privateKey, signMethod)
	assert.NoError(t, err)

	certificate := DefaultCertificate()
	_, _, _, _, _, err = ReadJWTToken(at, certificate, signMethod)
	assert.Error(t, err)
	log.Println(err.Error())
}
