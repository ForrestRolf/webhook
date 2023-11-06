package action

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Http struct {
	Hook         *hook.Hook
	Action       *hook.HttpAction
	LogModel     *model.LogClient
	WebhookModel *model.WebhookClient
}

func NewHttpAction(action *hook.HttpAction, hook *hook.Hook, log *model.LogClient, webhook *model.WebhookClient) *Http {
	return &Http{
		Action:       action,
		Hook:         hook,
		LogModel:     log,
		WebhookModel: webhook,
	}
}

func (h *Http) Send(args map[string]string) {
	_, err := url.ParseRequestURI(h.Action.Url)
	if err != nil {
		h.LogModel.AddErrorLog(h.Hook, "[HTTP] "+err.Error())
		return
	}
	h.LogModel.AddInfoLog(h.Hook, fmt.Sprintf("[HTTP] Send request to: [%s] %s, withAuth: %t", h.Action.Method, h.Action.Url, h.Action.AuthToken != ""))

	payload := h.Action.Payload
	for k, v := range args {
		payload = strings.ReplaceAll(payload, fmt.Sprintf("${%s}", k), v)
	}
	bodyReader := bytes.NewReader([]byte(payload))

	req, err := http.NewRequest(h.Action.Method, h.Action.Url, bodyReader)
	if err != nil {
		h.LogModel.AddErrorLog(h.Hook, "[HTTP] "+err.Error())
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
		h.LogModel.AddErrorLog(h.Hook, "[HTTP] "+err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if h.Action.SaveResponse {
		h.LogModel.AddDebugLog(h.Hook, "[HTTP] Response: "+string(body))
	}
	h.LogModel.AddLog(h.Hook, "[HTTP] Send successfully")
}
