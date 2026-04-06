package error_handler

import (
	"daos_core/internal/domain/dto/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpErrHandler(ctx *gin.Context, e error) {
	switch {
	case errors.Is(e, ErrCantParseJson):
		sendJSON(ctx, e, http.StatusBadRequest)
		return
	case errors.Is(e, ErrCantParseDTO):
		sendJSON(ctx, e, http.StatusInternalServerError)
		return
	default:
		sendJSON(ctx, e, http.StatusBadRequest)
		return
	}
}

func sendJSON(ctx *gin.Context, e error, code int) {
	ctx.JSON(code, common.RegularResponseDTO[any]{
		Ok:          false,
		Description: e.Error(),
	})
}

// когда - нибудь...
var (
	ErrCantParseToJson = errors.New("can't parse to json")
	ErrCantParseJson   = errors.New("can't parse json")
	ErrCantParseDTO    = errors.New("can't parse dto")
)
