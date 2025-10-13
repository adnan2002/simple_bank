package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



const minCharsForSecret = 32

type JWToken struct {
	Secret string
}


func NewJwtToken(secret string) (Maker, error) {
	if minCharsForSecret > len(secret){
		return nil, errors.New("invalid Duration")
	}


	return &JWToken{
		Secret: secret,
	}, nil
}



func (maker *JWToken)  CreateToken(username string, duration time.Duration)  (string, *Payload, error){
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.Secret))

	return token, payload, err

}

func (maker *JWToken) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != "HS256" {
			return nil, errors.New("invalid Token")
		}
		return []byte(maker.Secret), nil
	}

	parsedToken, err := jwt.ParseWithClaims(token, &Payload{} ,keyFunc)

	if err != nil {
		return nil, err
	}


	payload, ok := parsedToken.Claims.(*Payload)
	if !ok{
		return nil, errors.New("invalid Claim")
	}

	return payload, nil
}