package itc

import (
	"crypto/ecdsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests
	// The token's expiration time, in Unix epoch time; tokens that expire more than 20 minutes in the future are not valid
	jwtExpirationInterval = 20 * time.Minute
)

type itcJWT struct {
	KeyID      string
	IssuerID   string
	PrivateKey *ecdsa.PrivateKey

	expiresAt time.Time
	encoded   string
}

func (j *itcJWT) Encode() (string, error) {
	if time.Now().After(j.expiresAt) {
		if err := j.createNew(); err != nil {
			return "", err
		}
	}

	return j.encoded, nil
}

func (j *itcJWT) createNew() error {
	expiresAt := time.Now().Add(jwtExpirationInterval)

	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": j.KeyID,
			"typ": "JWT",
		},
		Claims: jwt.MapClaims{
			"iss": j.IssuerID,
			"exp": expiresAt.Unix(),
			"aud": "appstoreconnect-v1",
		},
		Method: jwt.SigningMethodES256,
	}
	encoded, err := jwtToken.SignedString(j.PrivateKey)
	if err != nil {
		return err
	}

	j.encoded = encoded
	j.expiresAt = expiresAt

	return nil
}
