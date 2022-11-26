package main

import (
	"context"
	proto "easy-go-iot/user-web/proto"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"easy-go-iot/user-web/global"
	"easy-go-iot/user-web/initialize"
	myvalidator "easy-go-iot/user-web/validator"
)

func main() {
	//1. 初始化logger
	initialize.InitLogger()

	//2. 初始化配置文件
	initialize.InitConfig()

	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//5. 初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info ， warn， error， fetal
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

	go func() {
		time.Sleep(10 * time.Second)
		userConn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", "userrpc-service", 50001),
			grpc.WithInsecure())
		if err != nil {
			zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		}
		userClient := proto.NewUserClient(userConn)
		rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
			Id: int32(1),
		})
		if err != nil {
			zap.S().Error(err)
		}
		zap.S().Info(rsp.Mobile, rsp.NickName, rsp.Password)
	}()

	go func() {
		time.Sleep(20 * time.Second)
		userConn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
			grpc.WithInsecure())
		if err != nil {
			zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		}
		userClient := proto.NewUserClient(userConn)
		rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
			Id: int32(1),
		})
		if err != nil {
			zap.S().Error(err)
		}
		zap.S().Info(rsp.Mobile, rsp.NickName, rsp.Password)
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
