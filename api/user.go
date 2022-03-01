package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	hashedPassword, err := util.GenerateHashedPassowrd(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
	}

	user, err := server.db.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
