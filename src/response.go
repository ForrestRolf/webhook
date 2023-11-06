package src

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Custom(context *gin.Context, httpStatus int, code string, data interface{}, msg string) {
	context.JSON(httpStatus, gin.H{
		"meta": gin.H{
			"code":    code,
			"message": msg,
		},
		"payload": data,
		"ts":      time.Now().UnixMilli(),
	})
}

func (r *Response) Success(context *gin.Context, data interface{}, msg string) {
	r.Custom(context, http.StatusOK, "OK", data, msg)
}

func (r *Response) Fail(context *gin.Context, msg string, data interface{}) {
	r.Custom(context, http.StatusInternalServerError, "SERVER_ERROR", data, msg)
}

func (r *Response) BadRequest(context *gin.Context, msg string, data interface{}) {
	r.Custom(context, http.StatusBadRequest, "BAD_REQUEST", data, msg)
}

func (r *Response) NotFound(context *gin.Context, msg string) {
	r.Custom(context, http.StatusNotFound, "NOT_FOUND", nil, msg)
}

func (r *Response) Unauthorized(context *gin.Context) {
	r.Custom(context, http.StatusUnauthorized, "UNAUTHORIZED", nil, "")
}
