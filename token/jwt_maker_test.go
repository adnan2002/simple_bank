package token

import (
	"testing"
	"time"

	"example.com/db/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)



func TestJWTMAker(t *testing.T){
	maker, err := NewJwtToken(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, paylod,  err := maker.CreateToken(username, duration)
	require.NotEmpty(t, paylod)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}


func TestExpiredJWTToken(t *testing.T){
	maker, err := NewJwtToken(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NotEmpty(t, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)

}


func TestInvalidJWTTokenAlgNone(t *testing.T){
	maker, err := NewJwtToken(util.RandomString(32))
	require.NoError(t, err)

	jwtStruct := maker.(*JWToken)

	token, err := CreateFakeToken(util.RandomOwner(), time.Minute, jwtStruct.Secret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, payload)



}



func CreateFakeToken(username string, duration time.Duration, secret string)  (string, error){
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	return jwtToken.SignedString([]byte(secret))

}