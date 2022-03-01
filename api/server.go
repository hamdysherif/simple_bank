package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hamdysherif/simplebank/db/sqlc"
)

type Server struct {
	db     db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{db: store}
	router := gin.Default()

	// register the custom currency validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/transfers", server.transferAmount)

	server.router = router
	return server
}

func (server *Server) Start(address string) {
	server.router.Run(address)
}

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
