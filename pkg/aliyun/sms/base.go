// This file is auto-generated, don't edit it. Thanks.
package sms

import (
	"ipfast_server/pkg/util/log"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Client = dysmsapi20170525.Client

func NewClient(key, secret *string) (*dysmsapi20170525.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     key,
		AccessKeySecret: secret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	config.ConnectTimeout = tea.Int(5000)
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SendSms(phone, signName, templateCode, templateParam string, key, secret *string) (err error) {
	client, err := NewClient(key, secret)
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	}
	runtime := &util.RuntimeOptions{}
	err = func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}
		return nil
	}()
	if err != nil {
		log.Error("send sms failed: %v", err)
		return err
	}
	return nil
}
