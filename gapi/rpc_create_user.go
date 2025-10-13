package gapi

import (
	"context"

	db "example.com/db/sqlc"
	"example.com/pb"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user, err := server.Store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.GetUsername(),
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return nil, pgErr
			}
		}
		return nil, err
	}


	return &pb.CreateUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
	}, nil

}


func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.Store.GetUser(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "user does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "invalid username or password")
	}


	return &pb.LoginUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
	}, nil

}