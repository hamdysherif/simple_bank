package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hamdysherif/simplebank/db/sqlc"
)

type Server struct {
	db     db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{db: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) {
	server.router.Run(address)
}

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
