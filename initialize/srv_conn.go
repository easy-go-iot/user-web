package initialize

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"

	"easy-go-iot/user-web/global"
	proto "easy-go-iot/user-web/proto"
)

// TODO 初始化grpc连接
func InitSrvConn() {
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
		grpc.WithInsecure())
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient

	test()
}

func test() {
	go func() {
		time.Sleep(50 * time.Second)
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
		time.Sleep(60 * time.Second)
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
}
