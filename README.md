# CREDENTA

## About

A simple library for storing user/group credential management. Created using Golang programming language.

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

This library also help you to work with JWT. It uses `github.com/SermoDigital/jose` to work
sith JWT. You will need to have `openssl` tooling to create RSA private and public key
to generate the required keys so your JWT will be secured.

## A note for JOSE library

To add JOSE, you have to maksure JOSE is in your `go.mod` as follows.

```text
require (
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc
	...
}

exclude github.com/SermoDigital/jose v0.9.1
```

first, you call the following command to add JOSE

```shell
$ go get github.com/SermoDigital/jose
$ go get github.com/SermoDigital/jose@v0.9.2-0.20180104203859-803625baeddc
```

And the, you can edit your `go.mod` file like the above.

## How to work with JWT Token

1. Create your private key
2. Create public key from your private key

### 1. Create your private key

```shell
$ openssl genrsa -out sample_key.priv 2048
```

To load the saved keys, 

```go
import (
    "github.com/newm4n/credenta"
)

privateKey := credenta.LoadPrivateKeyFromFile("path/to/sample_key.priv")
```

You may notice that `LoadPrivateKeyFromFile` function does not return an `error`
instance. Its because the function will automatically return default PrivateKey if
it founds an error. Bellow shows function that create a default Private key.

```go
import (
    "github.com/newm4n/credenta"
)

privateKey := credenta.GetDefaultPrivateKey()
```

### 2. Create public key from your private key

```shell
$ openssl rsa -in sample_key.priv -pubout > sample_key.pub
```

To load the saved keys,

```go
import (
    "github.com/newm4n/credenta"
)

publicKey := credenta.LoadPublicKeyFromFile("path/to/sample_key.pub")
```

You may notice that `LoadPublicKeyFromFile` function does not return an `error`
instance. Its because the function will automatically return default PublicKey if
it founds an error. Bellow shows function that create a default Public key.

```go
import (
    "github.com/newm4n/credenta"
)

publicKey := credenta.GetDefaultPublicKey()
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