package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	ResponseOk                   = "OK"
	ResponseServerError          = "SERVER_ERROR"
	ResponseBadRequestError      = "BAD_REQUEST"
	ResponseTimedOut             = "TIMED_OUT"
	ResponseUnauthorizedError    = "UNAUTHORIZED"
	ResponseUnauthenticatedError = "UNAUTHENTICATED"
)

// BaseResponse represents base http response
type BaseResponse struct {
	Status       string      `json:"status"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data"`
}

func getTimedOutRespBody() string {
	timeOutBody := BaseResponse{
		Status: ResponseTimedOut,
	}
	marshal, _ := json.Marshal(timeOutBody)

	return string(marshal)
}

// ===== Response wrapper using gin =====

func WriteServerErrorResponse(ctx *gin.Context, functionName string, err error) {
	if err == nil {
		WriteNotOkResponse(ctx, http.StatusInternalServerError, ResponseServerError)
	}

	errMessage := err.Error()
	log.Println(fmt.Sprintf("[ERROR] [%s] %s", functionName, errMessage))
	WriteNotOkResponseWithErrMsg(ctx, http.StatusInternalServerError, ResponseServerError, errMessage)
}

func WriteBadRequestResponse(ctx *gin.Context, status string) {
	WriteNotOkResponse(ctx, http.StatusBadRequest, status)
}

func WriteBadRequestResponseWithErrMsg(ctx *gin.Context, status string, err error) {
	if err == nil {
		WriteBadRequestResponse(ctx, status)
	}

	WriteNotOkResponseWithErrMsg(ctx, http.StatusBadRequest, status, err.Error())
}

func WriteNotFoundResponse(ctx *gin.Context, status string) {
	WriteNotOkResponse(ctx, http.StatusNotFound, status)
}

func WriteUnauthorizedResponse(ctx *gin.Context) {
	WriteNotOkResponse(ctx, http.StatusUnauthorized, ResponseUnauthorizedError)
}

func WriteUnauthenticatedResponse(ctx *gin.Context) {
	WriteNotOkResponse(ctx, http.StatusUnauthorized, ResponseUnauthenticatedError)
}

func WriteTimedOutResponse(ctx *gin.Context) {
	WriteNotOkResponse(ctx, http.StatusGatewayTimeout, ResponseTimedOut)
}

// WriteOkResponse writes 200 response using gin.
func WriteOkResponse(ctx *gin.Context, data interface{}) {
	resp := BaseResponse{
		Status: ResponseOk,
		Data:   data,
	}
	WriteResponse(ctx, resp, http.StatusOK)
}

// WriteNotOkResponse writes non 200 response.
func WriteNotOkResponse(ctx *gin.Context, statusCode int, status string) {
	resp := BaseResponse{
		Status: status,
	}
	WriteResponse(ctx, resp, statusCode)
}

// WriteNotOkResponseWithErrMsg writes non 200 responses with error message.
func WriteNotOkResponseWithErrMsg(ctx *gin.Context, statusCode int, status, errMsg string) {
	resp := BaseResponse{
		Status:       status,
		ErrorMessage: errMsg,
	}

	WriteResponse(ctx, resp, statusCode)
}

func WriteResponse(ctx *gin.Context, resp BaseResponse, statusCode int) {
	ctx.JSON(statusCode, resp)
	return
}
