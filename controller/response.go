package controller

import (
	"cc"
	"net/http"

	"boframe/pkg/errcode"
)

const (
	SuccessCode    = 1
	SuccessMessage = "成功"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func responseJSON(c *cc.Context, httpCode int, payload *ResponseData) {
	c.JSON(httpCode, payload)
}

func ResponseErrorCode(c *cc.Context, code *errcode.ErrCode) {
	res := &ResponseData{
		Code:    code.Code,
		Message: code.Message,
		Data:    code.Data,
	}
	responseJSON(c, http.StatusOK, res)
}

func ResponseErrorCodeWithHTTPCode(c *cc.Context, httpCode int, code *errcode.ErrCode) {
	res := &ResponseData{
		Code:    code.Code,
		Message: code.Message,
		Data:    code.Data,
	}
	responseJSON(c, httpCode, res)
}

func ResponseSuccess(c *cc.Context) {
	responseJSON(c, http.StatusOK, nil)
}

func ResponseSuccessWithData(c *cc.Context, payload interface{}) {
	res := &ResponseData{
		Code:    SuccessCode,
		Message: SuccessMessage,
		Data:    payload,
	}
	responseJSON(c, http.StatusOK, res)
}
