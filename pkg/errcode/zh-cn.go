package errcode

var zhCN = map[int]string{
	ErrNotLogin:               "未登录",
	ErrInvalidParam:           "参数无效",
	ErrInvalidPageInfo:        "无效的分页信息",
	ErrDB:                     "数据库错误",
	ErrServer:                 "服务器错误",
	ErrNoRight:                "没有权限",
	ErrRegisterParamInvalid:   "注册参数错误",
	ErrRegisterUsernameExists: "用户名已经注册",
	ErrUnknown:                "发生未知错误，请联系管理员",
	ErrInvalidToken:           "登录令牌无效",
	ErrTokenExpired:           "登录令牌已失效",
	ErrUserNotExists:          "用户不存在",
	ErrLoginWrongPassword:     "用户名或者密码错误",
	ErrSGiftSenderOrReceiver:  "送礼失败，参与人信息错误",
	ErrSGiftSendToSelf:        "送礼失败，不可送给自己",
	ErrSGiftInvalidGift:       "送礼失败，礼物信息错误",
	ErrSGiftInvalidUser:       "送礼失败，无效的参与人",
	ErrAnchorNotExists:        "主播不存在",
	ErrAnchorInvalid:          "主播信息无效",
}
