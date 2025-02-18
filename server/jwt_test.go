package server_test

import (
	"github.com/astroniumm/go-asyncapi/config"
	"github.com/astroniumm/go-asyncapi/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJWTManager(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	JWTMgr := server.NewJWTManager(conf)
	userID := uuid.New()
	tokenPair, err := JWTMgr.GenerateTokenPair(userID)
	require.NoError(t, err)

	accessTokenSubject, err := tokenPair.AccessToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userID.String(), accessTokenSubject)

	accessTokenIssuer, err := tokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+conf.ServerHost+":"+conf.ServerPort, accessTokenIssuer)

	refreshTokenSubject, err := tokenPair.RefreshToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userID.String(), refreshTokenSubject)

	refreshTokenIssuer, err := tokenPair.RefreshToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+conf.ServerHost+":"+conf.ServerPort, refreshTokenIssuer)

}
