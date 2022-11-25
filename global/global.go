package global

import (
	"easy-go-iot/user-web/config"
	proto "easy-go-iot/user-web/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
