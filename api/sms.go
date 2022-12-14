package api

import (
	"context"
	"easy-go-iot/user-web/forms"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"easy-go-iot/user-web/global"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	zap.S().Info("SendSms here")
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing",
		// TODO key
		global.ServerConfig.AliSmsInfo.ApiKey,
		// TODO secret
		global.ServerConfig.AliSmsInfo.ApiSecret)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile //手机号
	// TODO 签名管理, 签名名称
	request.QueryParams["SignName"] = global.ServerConfig.AliSmsInfo.SignName //阿里云验证过的项目名 自己设置
	// TODO 模板管理, 模板CODE
	request.QueryParams["TemplateCode"] = global.ServerConfig.AliSmsInfo.TemplateCode //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"               //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	zap.S().Info(response)
	zap.S().Info(client.DoAction(request, response))
	if err != nil {
		zap.S().Error(err)
	}
	//将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	zap.S().Infof("redis key %s", sendSmsForm.Mobile)
	cmd := rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	result, err := cmd.Result()
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(result)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
