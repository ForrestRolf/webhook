package action

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Http struct {
	Hook         *hook.Hook
	Action       *hook.HttpAction
	LogModel     *model.ActionLogClient
	WebhookModel *model.WebhookClient
}

func NewHttpAction(action *hook.HttpAction, hook *hook.Hook, webhook *model.WebhookClient, log *model.ActionLogClient) *Http {
	return &Http{
		Action:       action,
		Hook:         hook,
		LogModel:     log,
		WebhookModel: webhook,
	}
}

func (h *Http) Send(args map[string]string) {
	start := time.Now().UnixMilli()

	_, err := url.ParseRequestURI(h.Action.Url)
	if err != nil {
		h.LogModel.AddErrorLog(err.Error())
		return
	}
	h.LogModel.AddInfoLog(fmt.Sprintf("[HTTP] Send request to: [%s] %s, withAuth: %t", h.Action.Method, h.Action.Url, h.Action.AuthToken != ""))

	payload := h.Action.Payload
	for k, v := range args {
		payload = strings.ReplaceAll(payload, fmt.Sprintf("${%s}", k), v)
	}
	bodyReader := bytes.NewReader([]byte(payload))
	h.LogModel.AddDebugLog("Send body: " + payload)

	req, err := http.NewRequest(h.Action.Method, h.Action.Url, bodyReader)
	if err != nil {
		h.LogModel.AddErrorLog(err.Error())
		return
	}

	if h.Action.ContentType != "" {
		req.Header.Set("Content-Type", h.Action.ContentType)
	}
	if h.Action.AuthToken != "" {
		req.Header.Set("Authorization", h.Action.AuthToken)
	}

	client := http.Client{
		Timeout: time.Duration(h.Action.Timeout) * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		h.LogModel.AddErrorLog(err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if h.Action.SaveResponse {
		h.LogModel.AddDebugLog(fmt.Sprintf("[HTTP] Response status: %s, body: %s", res.Status, string(body)))
	}

	end := time.Now().UnixMilli()
	h.LogModel.AddLog("Send successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
}
