package jwt

import (
	"errors"
	"time"
	
	"boframe/settings"
	"github.com/dgrijalva/jwt-go"
)

var Secret = []byte("whatasecret")
var ErrInvalidToken = errors.New("invalid token")

// jwt过期时间, 按照实际环境设置
const expiration = 2 * time.Hour

type Claims struct {
	// 自定义字段, 可以存在用户名, 用户ID, 用户角色等等
	UserId   int64
	Username string
	// jwt.StandardClaims包含了官方定义的字段
	jwt.StandardClaims
}

func GenToken(uid int64, username string) (string, error) {
	t := time.Now()
	// 创建声明
	a := &Claims{
		UserId:   uid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t.Add(expiration).Unix(), // 过期时间
			IssuedAt:  t.Unix(),                 // 签发时间
			Issuer:    settings.Conf.Name,       // 签发者
		},
	}
	
	// 用指定的哈希方法创建签名对象
	tt := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	// 用上面的声明和签名对象签名字符串token
	// 1. 先对Header和PayLoad进行Base64URL转换
	// 2. Header和PayLoadBase64URL转换后的字符串用.拼接在一起
	// 3. 用secret对拼接在一起之后的字符串进行HASH加密
	// 4. 连在一起返回
	return tt.SignedString(Secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	var c = new(Claims)
	// 第三个参数: 提供一个回调函数用于提供要选择的秘钥, 回调函数里面的token参数,是已经解析但未验证的,可以根据token里面的值做一些逻辑, 如`kid`的判断
	token, err := jwt.ParseWithClaims(tokenStr, c,
		func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return c, nil
	}
	return nil, ErrInvalidToken
}

func IsErrInvalidToken(err error) bool {
	return errors.Is(err, ErrInvalidToken)
}
