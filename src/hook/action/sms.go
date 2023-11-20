package action

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Sms struct {
	Hook     *hook.Hook
	Action   *hook.SmsAction
	Profile  *model.SmsProfile
	LogModel *model.ActionLogClient
	Message  string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewSmsAction(action *hook.SmsAction, profile *model.SmsProfile, hook *hook.Hook, log *model.ActionLogClient) *Sms {
	return &Sms{
		Hook:     hook,
		Action:   action,
		Profile:  profile,
		LogModel: log,
	}
}

func (sms *Sms) Send(args map[string]string) {
	start := time.Now().UnixMilli()

	message := sms.Action.Content
	for k, v := range args {
		message = strings.ReplaceAll(message, fmt.Sprintf("${%s}", k), v)
	}
	sms.Message = message

	switch sms.Profile.Provider {
	case hook.ActionSmsTwilioDriver:
		sms.SendTwilioSms()
	case hook.ActionSmsBurstDriver:
		sms.SendBurstSms()
	case hook.ActionSmsPlivoDriver:
		sms.SendPlivoSms()
	case hook.ActionSmsSNSDriver:
		sms.sendSNS()
	default:
		sms.LogModel.AddWarnLog(fmt.Sprintf("Unsupported sms provider: %s", sms.Profile.Provider))
	}

	end := time.Now().UnixMilli()
	sms.LogModel.AddLog("Send successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
}

func (sms *Sms) SendTwilioSms() {
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sms.Profile.AK)
	message := map[string]string{
		"to":   sms.Action.To,
		"from": sms.Profile.From,
		"body": sms.Message,
	}
	auth := "Basic " + basicAuth(sms.Profile.AK, sms.Profile.SK)

	sms.request(url, message, auth)
}

func (sms *Sms) SendBurstSms() {
	url := "https://api.transmitsms.com/send-sms.json"
	message := map[string]string{
		"to":      sms.Action.To,
		"from":    sms.Profile.From,
		"message": sms.Message,
	}
	auth := "Basic " + basicAuth(sms.Profile.AK, sms.Profile.SK)

	sms.request(url, message, auth)
}

func (sms *Sms) SendPlivoSms() {
	url := fmt.Sprintf("https://api.plivo.com/v1/Account/%s/Message", sms.Profile.AK)
	message := map[string]string{
		"to":   sms.Action.To,
		"from": sms.Profile.From,
		"body": sms.Message,
	}
	auth := "Basic " + basicAuth(sms.Profile.AK, sms.Profile.SK)

	sms.request(url, message, auth)
}

func (sms *Sms) request(url string, message map[string]string, auth string) {
	client := &http.Client{}
	bodyJson, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		sms.LogModel.AddErrorLog(err.Error())
		return
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		sms.LogModel.AddErrorLog(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		sms.LogModel.AddErrorLog(fmt.Sprintf("Response code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	if sms.Hook.Debug {
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			sms.LogModel.AddErrorLog(err.Error())
		}
		sms.LogModel.AddDebugLog(fmt.Sprintf("Response body: %s", response))
	}
}

type SNSPublishAPI interface {
	Publish(ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func PublishSNSMessage(c context.Context, api SNSPublishAPI, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return api.Publish(c, input)
}
func (sms *Sms) sendSNS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sms.Profile.AK, sms.Profile.SK, "")), config.WithRegion(sms.Profile.Region))
	if err != nil {
		sms.LogModel.AddErrorLog(err.Error())
		return
	}

	client := sns.NewFromConfig(cfg)
	input := &sns.PublishInput{
		Message:  &sms.Message,
		TopicArn: &sms.Action.TopicArn,
	}

	result, err := PublishSNSMessage(context.TODO(), client, input)
	if err != nil {
		sms.LogModel.AddErrorLog(err.Error())
		return
	}

	if sms.Hook.Debug {
		sms.LogModel.AddDebugLog(fmt.Sprintf("Message ID: %s, Sequence Number: %s", result.MessageId, result.SequenceNumber))
	}
}
