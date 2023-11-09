package action

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Dispatcher struct {
	Hook         *hook.Hook
	Action       *hook.DispatcherAction
	LogModel     *model.LogClient
	WebhookModel *model.WebhookClient
	Request      *hook.Request
}

func NewDispatcherAction(action *hook.DispatcherAction, hook *hook.Hook, log *model.LogClient, webhook *model.WebhookClient, req *hook.Request) *Dispatcher {
	return &Dispatcher{
		Action:       action,
		Hook:         hook,
		LogModel:     log,
		WebhookModel: webhook,
		Request:      req,
	}
}

func (d *Dispatcher) compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (d *Dispatcher) Send(args map[string]string) {
	ok := true
	for k, v := range d.Action.If {
		if val, o := args[k]; o {
			ok = d.compare(v, val)
			if !ok {
				break
			}
		} else {
			ok = false
		}
	}

	if !ok {
		return
	}
	d.LogModel.AddInfoLog(d.Hook, fmt.Sprintf("Send to %s", d.Action.Url))

	bodyReader := bytes.NewReader(d.Request.Body)
	req, err := http.NewRequest(d.Action.Method, d.Action.Url, bodyReader)
	if err != nil {
		d.LogModel.AddErrorLog(d.Hook, "[Dispatcher] "+err.Error())
		return
	}
	req.Header.Set("Content-Type", d.Request.ContentType)
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		d.LogModel.AddErrorLog(d.Hook, "[Dispatcher] "+err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	d.LogModel.AddDebugLog(d.Hook, "[Dispatcher] Response: "+string(body))
	d.LogModel.AddLog(d.Hook, "[Dispatcher] Send successfully")
}
