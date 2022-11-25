package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"easy-go-iot/user-web/global"
)

func InitConfig() {

	configFileName := fmt.Sprintf("conf/server-config.yaml")

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.ServerConfig)
}
