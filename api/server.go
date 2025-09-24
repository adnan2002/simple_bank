package api

import (
	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Store  *db.Store
	Router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	r := gin.Default()
	server := &Server{
		Store:  store,
		Router: r,
	}

	server.Router.POST("/accounts", server.CreateAccount)
	server.Router.GET("/accounts/:id", server.GetAccount)
	server.Router.GET("/accounts", server.ListAccounts)

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
