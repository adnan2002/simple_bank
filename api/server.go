package api

import (
	db "example.com/db/sqlc"
	"example.com/db/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Store  db.Store
	Router *gin.Engine
}

func NewServer(store db.Store) *Server {
	r := gin.Default()
	server := &Server{
		Store:  store,
		Router: r,
	}

	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		value.RegisterValidation("currency", util.Currency)
	}

	server.Router.POST("/accounts", server.CreateAccount)
	server.Router.GET("/accounts/:id", server.GetAccount)
	server.Router.GET("/accounts", server.ListAccounts)
	server.Router.POST("/transfers", server.CreateTransfer)
	server.Router.POST("/users",server.CreateUser)
	server.Router.GET("/users/:username", server.GetUser)


	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
