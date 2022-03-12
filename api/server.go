package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/token"
	"github.com/hamdysherif/simplebank/util"
)

type Server struct {
	db         db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer generate a new server
func NewServer(store db.Store, config util.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.SemmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create the tokenmaker: %w", err)
	}
	server := &Server{db: store, tokenMaker: tokenMaker}

	registerCustomValidators()

	server.SetupRouter()
	return server, nil
}

//registerCustomValidators register any custom validator
func registerCustomValidators() {
	// register the custom currency validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
}

// SetupRouter setup the routers and urls for the server
func (server *Server) SetupRouter() {
	router := gin.Default()

	authorized := router.Group("/")
	{
		authorized.Use(Authentication(server.tokenMaker))
		authorized.POST("/accounts", server.createAccount)
		authorized.GET("/accounts", server.listAccounts)
		authorized.GET("/accounts/:id", server.getAccount)
		authorized.POST("/transfers", server.transferAmount)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

func (server *Server) Start(address string) {
	server.router.Run(address)
}

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
