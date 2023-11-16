package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Slack struct {
	Hook     *hook.Hook
	Action   *hook.SlackAction
	LogModel *model.ActionLogClient
}

func NewSlackAction(action *hook.SlackAction, hook *hook.Hook, log *model.ActionLogClient) *Slack {
	return &Slack{
		Hook:     hook,
		Action:   action,
		LogModel: log,
	}
}

func (s *Slack) Send(args map[string]string) {
	start := time.Now().UnixMilli()

	message := s.Action.Message
	for k, v := range args {
		message = strings.ReplaceAll(message, fmt.Sprintf("${%s}", k), v)
	}

	client := &http.Client{}
	body := map[string]string{
		"text": message,
	}
	if s.Action.Channel != "" {
		channel := s.Action.Channel
		if !strings.HasPrefix(channel, "#") {
			channel = "#" + channel
		}
		body["channel"] = channel
	}
	bodyJson, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", s.Action.WebhookUrl, bytes.NewReader(bodyJson))
	if err != nil {
		s.LogModel.AddErrorLog(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		s.LogModel.AddErrorLog(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.LogModel.AddErrorLog(fmt.Sprintf("Response code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	s.LogModel.AddDebugLog(fmt.Sprintf("Response body: %s", response))

	end := time.Now().UnixMilli()
	s.LogModel.AddLog("Send successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
}
