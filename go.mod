module boframe

go 1.16

require (
	cc v0.0.0
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-redis/redis/v8 v8.11.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/spf13/viper v1.8.1
	go.mongodb.org/mongo-driver v1.7.1 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace cc => ./pkg/cc
