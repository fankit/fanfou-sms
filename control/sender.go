package control

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
	"go.uber.org/zap"
	"tsmsrv/conf"
	"tsmsrv/utils"
)

type SendSmsSrv struct {
	scli     *sms.Client
}

func NewSendSmsSrv() *SendSmsSrv {

	credential := common.NewCredential(
		conf.GlobConfig.SmsSection().Key(`secretid`).String(),
		conf.GlobConfig.SmsSection().Key(`secretkey`).String(),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	return &SendSmsSrv{scli: client}
}

func (s *SendSmsSrv) SendSms(phns []string, cxt []string) (err error) {
	var (
		phones    []*string
		context   []*string
		sresponse *sms.SendSmsResponse
	)
	for _, p := range phns {
		ph := `86` + p
		phones = append(phones, &ph)
	}

	for _, ct := range cxt {
		context = append(context, &ct)
	}

	sign := conf.GlobConfig.SmsSection().Key(`sign`).String()
	sdkappid := conf.GlobConfig.SmsSection().Key(`appid`).String()
	tmpid := conf.GlobConfig.SmsSection().Key(`templateid`).String()


	smsRequest := sms.NewSendSmsRequest()
	smsRequest.PhoneNumberSet = phones
	smsRequest.Sign = &sign
	smsRequest.SmsSdkAppid = &sdkappid
	smsRequest.TemplateID = &tmpid
	smsRequest.TemplateParamSet = context

	if sresponse, err = s.scli.SendSms(smsRequest); err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return err.(*errors.TencentCloudSDKError)
		} else {
			return err
		}
	}

	//fmt.Printf("%s", sresponse.ToJsonString())
	utils.Logger.Log.Info(`SendSms`, zap.Any(`desc`, sresponse.ToJsonString()))

	return
}