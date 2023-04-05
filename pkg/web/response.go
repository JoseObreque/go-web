package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
The ErrorResponse struct represents the response from the server when an error occurs.

	Status (int): HTTP Status Code as an integer. Example: 200.
	Code (string): HTTP Status Code as a string. Example: "OK".
	Message (string): Error message.
*/
type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

/*
The Response struct represents a successful response from the server.
*/
type Response struct {
	Data interface{} `json:"data"`
}

/*
The Success function emits a successful response to the client.

	Status (int): HTTP Status Code as an integer. Example: 200.
	Data (string): Any data required in the response to the client.
*/
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		data,
	})
}

/*
The Failure function emits a failed response to the client.

	Status (int): HTTP Status Code as an integer. Example: 200.
	err (error): The error associated to the failed response to the client.
*/
func Failure(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{
		Status:  status,
		Code:    http.StatusText(status),
		Message: err.Error(),
	})
}
