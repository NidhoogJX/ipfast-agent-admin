// This file is auto-generated, don't edit it. Thanks.
package captcha_v2

import (
	"fmt"
	"ipfast_server/pkg/util/log"

	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type Client = captcha20230305.Client

func NewClient(key, secret *string) (*captcha20230305.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     key,
		AccessKeySecret: secret,
	}
	config.Endpoint = tea.String("captcha.cn-shanghai.aliyuncs.com")
	config.ConnectTimeout = tea.Int(5000)
	config.ReadTimeout = tea.Int(5000)
	client, err := captcha20230305.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Verif(client *captcha20230305.Client, CaptchaVerifyParam string) (err error) {
	defer func() {
		if r := tea.Recover(recover()); r != nil {
			log.Error("阿里云验证码验证失败:%s", r.Error())
		}
	}()
	// 创建APi请求
	request := &captcha20230305.VerifyIntelligentCaptchaRequest{}
	request.SceneId = tea.String("vjyflq1t")
	request.CaptchaVerifyParam = tea.String(CaptchaVerifyParam)
	resp, err := client.VerifyIntelligentCaptcha(request)
	if err != nil {
		err = fmt.Errorf("阿里云验证码验证失败:%s", err)
		log.Error("阿里云验证码验证失败:%s", err.Error())
		return
	}
	log.Debug("阿里云验证码验证结果:%v,%v", *resp.Body.Result.VerifyResult, resp)
	if !*resp.Body.Result.VerifyResult {
		log.Error("阿里云验证码验证失败:%s", *resp.Body.Result.VerifyCode)
		err = fmt.Errorf("阿里云验证码验证失败%s", *resp.Body.Result.VerifyCode)
		return
	}
	return
}
