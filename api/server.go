package api

import (
	db "example.com/db/sqlc"
	"example.com/db/util"
	"example.com/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Config     util.Config
	TokenMaker token.Maker
	Store      db.Store
	Router     *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	r := gin.Default()
	tokenMaker, err := token.NewPasetoMaker(config.Token)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Store:      store,
		Router:     r,
		Config:     config,
		TokenMaker: tokenMaker,
	}

	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		value.RegisterValidation("currency", util.Currency)
	}

	authorized := server.Router.Group("/")

	authorized.Use(AuthMiddleware(server.TokenMaker))
	{
			authorized.GET("/users", server.GetUser)
			authorized.GET("/accounts", server.ListAccounts)
			authorized.POST("/accounts", server.CreateAccount)
			authorized.GET("/accounts/:id", server.GetAccount)
			authorized.POST("/transfers", server.CreateTransfer)

	}




	server.Router.POST("/users", server.CreateUser)
	server.Router.POST("/users/login", server.LoginUser)
	server.Router.POST("/users/auth/refresh", server.Refresh)

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
