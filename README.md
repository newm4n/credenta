# CREDENTA

## About

A simple library created using Golang programming language.

## Vision n Mission

Aiming to be a simple library that handles user, role and group management in
a simplest manner as possible. One could simply add this library to their
project using `git add` command  and it simply have a simple yet good breed
of user management for their application. 

We foresee that the system could comfortably handles like 10k user accounts. More 
accounts might need better improvement in the future.

## How to add

```shell
git add github.com/newm4n/credenta
```

## Todo

- Datastorage issue, SQLITE, POSTGRES, MongoDB or custom made file store
- Authentication management
- Token management using JWT

# JWT Token

## How to work with JWT Token

1. Create your private key
2. Create public key from your private key

### 1. Create your private key

```shell
openssl genrsa -out sample_key.priv 2048
```

### 2. Create public key from your private key

```shell
openssl rsa -in sample_key.priv -pubout > sample_key.pub
```

### 3. Create JWT Token

```go
import (
	"github.com/newm4n/credenta"
)

// Claim related informations
issuer := "TheIssuer"
subject := "TheSubject"
audience := []string{"audience1", "audience2"}
additional := map[string]interface{}{"map1": "value1"}
issuedAt := time.Now()
accessTokenAge := time.Minute * 5

// Encryption related information
privateKey := credenta.GetDefaultPrivateKey()
signMethod := crypto.SigningMethodRS256

// Function that generate the token
at, err := credenta.GenerateJWTToken(issuer,subject,audience,AccessTokenType,additional,issuedAt,issuedAt,issuedAt.Add(accessTokenAge),privateKey,signMethod)
```

### 4. Read and Validate JWT Token

```go
import (
	"github.com/newm4n/credenta"
)

// Get the public key for validation.
publicKey := credenta.GetDefaultPublicKey()

// err will not nil IF token is not valid, e.g expired or the signature not match
issuer, subject, audience, tokenType, additional, err := credenta.ReadJWTToken(at, publicKey, signMethod)
```