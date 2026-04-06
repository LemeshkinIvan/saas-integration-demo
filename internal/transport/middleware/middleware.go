package middleware

import (
	"daos_core/internal/domain/dto/common"
	"daos_core/internal/domain/services/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	service auth.Service
}

func NewMiddleware(a auth.Service) Middleware {
	return Middleware{service: a}
}

// set to http context accountID from bearer token string
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.GetHeader("Authorization")

		if value == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: "bearer token is empty",
			})
			return
		}

		token, err := parseBearer(value)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: err.Error(),
			})
			return
		}

		accountID, err := m.service.ValidateAccess(ctx, token)
		fmt.Println(accountID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: err.Error(),
			})
			return
		}

		ctx.Set("accountId", accountID)
		ctx.Next()
	}
}

// return token payload
func parseBearer(value string) (string, error) {
	str := strings.SplitN(value, " ", 2)
	if len(str) != 2 && str[0] != "Bearer" {
		return "", fmt.Errorf("its not bearer")
	}
	return str[1], nil
}

// no set to http context accountID from bearer token string
func (m *Middleware) AuthMiddlewareWithoutAccountIDParsing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.GetHeader("Authorization")

		if value == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: "bearer token is empty",
			})
			return
		}

		token, err := parseBearer(value)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: err.Error(),
			})
			return
		}

		_, err = m.service.ValidateAccess(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.RegularResponseDTO[any]{
				Ok:          false,
				Description: err.Error(),
			})
			return
		}

		ctx.Next()
	}
}

func (m *Middleware) Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
