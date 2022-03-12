package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hamdysherif/simplebank/token"
)

func Authentication(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.Request.Header.Get("authorization")

		tokenParts := strings.Fields(token)

		if len(tokenParts) != 2 {
			ctx.JSON(http.StatusUnauthorized, responseError(errors.New("invalid token")))
			ctx.Abort()
			return
		}

		tokenType := tokenParts[0]
		if strings.ToLower(tokenType) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, responseError(errors.New("unsupported token type")))
			ctx.Abort()
			return
		}

		tokenPayload := tokenParts[1]

		payload, err := tokenMaker.VerifyToken(tokenPayload)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, responseError(errors.New("invalid token")))
			ctx.Abort()
			return
		}

		ctx.Set("authpayload", payload)

		ctx.Next()
	}
}
