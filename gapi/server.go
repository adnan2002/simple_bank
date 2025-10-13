package gapi

import (
	db "example.com/db/sqlc"
	"example.com/db/util"
	"example.com/pb"
	"example.com/token"
)



type Server struct {
	Config     util.Config
	TokenMaker token.Maker
	Store      db.Store
	pb.UnimplementedUserServiceServer
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.Token)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Store:      store,
		Config:     config,
		TokenMaker: tokenMaker,
	}


	return server, nil
}

