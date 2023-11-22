package action

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Dispatcher struct {
	Hook     *hook.Hook
	Action   *hook.DispatcherAction
	LogModel *model.ActionLogClient
	Request  *hook.Request
}

func NewDispatcherAction(action *hook.DispatcherAction, hook *hook.Hook, req *hook.Request, log *model.ActionLogClient) *Dispatcher {
	return &Dispatcher{
		Action:   action,
		Hook:     hook,
		LogModel: log,
		Request:  req,
	}
}

func (d *Dispatcher) compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (d *Dispatcher) Send(args map[string]string) {
	start := time.Now().UnixMilli()

	ok := true
	for k, v := range d.Action.If {
		if val, o := args[k]; o {
			switch d.Action.Compare {
			case "eq":
				ok = d.compare(v, val)
			case "neq":
				ok = !d.compare(v, val)
			default:
				ok = d.compare(v, val)
			}
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
	d.LogModel.AddInfoLog(fmt.Sprintf("Send to %s (%s), method: %s", d.Action.WebhookName, d.Action.Url, d.Action.Method))
	d.LogModel.AddDebugLog("Send body: " + string(d.Request.Body))

	bodyReader := bytes.NewReader(d.Request.Body)
	req, err := http.NewRequest(d.Action.Method, d.Action.Url, bodyReader)
	if err != nil {
		d.LogModel.AddErrorLog(err.Error())
		return
	}
	req.Header.Set("Content-Type", d.Request.ContentType)
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		d.LogModel.AddErrorLog(err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	d.LogModel.AddDebugLog(fmt.Sprintf("Response status: %s, body: %s", res.Status, string(body)))

	end := time.Now().UnixMilli()
	d.LogModel.AddLog("Send successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
}
