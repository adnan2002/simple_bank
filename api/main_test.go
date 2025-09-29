package api

import (
	"os"
	"testing"
	"time"

	db "example.com/db/sqlc"
	"example.com/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)


func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		Token: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)

	require.NoError(t, err)
	return server
} 


func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)


	os.Exit(m.Run())
}