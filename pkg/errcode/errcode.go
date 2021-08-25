package errcode

type ErrCode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *ErrCode) Error() string {
	return zhCN[err.Code]
}

func (err *ErrCode) Text() string {
	return zhCN[err.Code]
}

func Text(code int) string {
	return zhCN[code]
}

func (err *ErrCode) WithData(data interface{}) *ErrCode {
	err.Data = data
	return err
}

func IsErrorCode(err error) (code *ErrCode, ok bool) {
	code, ok = err.(*ErrCode)
	return
}

// 3位数为基础错误，为原来的http错误码
const (
	ErrNotLogin = 401
	ErrNoRight  = 403
)

// 5位数错误码，其中`4XXXX`为客户端的错误，`5XXXX`为服务端错误
// 5位数中的第2-3位标记模块，4-5位标记业务错误码
const (
	ErrInvalidParam           = 40001
	ErrInvalidPageInfo        = 40002
	ErrRegisterParamInvalid   = 40101
	ErrRegisterUsernameExists = 40102
	ErrUserNotExists          = 40103
	ErrLoginWrongPassword     = 40104
	ErrInvalidToken           = 40201
	ErrTokenExpired           = 40202
	ErrSGiftSenderOrReceiver  = 40301
	ErrSGiftSendToSelf        = 40302
	ErrSGiftInvalidGift       = 40303
	ErrSGiftInvalidUser       = 40304
	ErrAnchorNotExists        = 40401
	ErrAnchorInvalid          = 40402

	ErrDB      = 50001
	ErrServer  = 50002
	ErrUnknown = 60001
)

func BuildWithData(code int, message string, data map[string]interface{}) *ErrCode {
	return build(code, message, data)
}

func BuildWithMessage(code int, message string) *ErrCode {
	return build(code, message, nil)
}

func Build(code int) *ErrCode {
	return build(code, Text(code), nil)
}

func build(code int, message string, data interface{}) *ErrCode {
	return &ErrCode{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
