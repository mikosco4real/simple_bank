package token

import (
	"testing"
	"time"

	"github.com/mikosco4real/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T){
	maker, err := NewPasetoMaker(util.GenericRandomString(32))
	require.NoError(t, err)

	email := util.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// Test that every field in the payload is what you expect
	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T){
	maker, err := NewPasetoMaker(util.GenericRandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomEmail(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
