package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/mikosco4real/simple_bank/db/sqlc"
	"github.com/mikosco4real/simple_bank/token"
	"github.com/mikosco4real/simple_bank/util"
)

// This Server will serve all http request to our application
type Server struct {
	store *db.Store
	tokenMaker token.Maker
	config util.Config
	router *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{
		store: store,
		config: config,
		tokenMaker: tokenMaker,
	}
	
	router := gin.Default()

	// Create the routes
	router.POST("/users", server.createUser) // TODO: To create users
	router.POST("/users/login", server.loginUser) // TODO: To login users

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users", server.listAllUsers) // TODO: To get all users
	authRoutes.GET("/users/:id", server.getUserById) // TODO: To get user by id

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {

	return gin.H{"error": err.Error()}
}