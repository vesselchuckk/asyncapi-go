package server_test

import (
	"github.com/astroniumm/go-asyncapi/config"
	"github.com/astroniumm/go-asyncapi/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func testConfig() *config.Config {
	return &config.Config{
		JwtSecret:  "mysecret",
		ServerHost: "localhost",
		ServerPort: "8080",
	}
}

func TestJWTManager(t *testing.T) {
	conf := testConfig()

	JWTMgr := server.NewJWTManager(conf)
	userID := uuid.New()
	tokenPair, err := JWTMgr.GenerateTokenPair(userID)
	require.NoError(t, err)

	require.True(t, JWTMgr.IsAccessToken(tokenPair.AccessToken))
	require.False(t, JWTMgr.IsAccessToken(tokenPair.RefreshToken))

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

	parsedAccessToken, err := JWTMgr.Parse(tokenPair.AccessToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.AccessToken, parsedAccessToken)

	parsedRefreshToken, err := JWTMgr.Parse(tokenPair.RefreshToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.RefreshToken, parsedRefreshToken)

}
