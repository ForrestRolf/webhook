package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"webhook/src"
	"webhook/src/hook"
	"webhook/src/hook/action"
	"webhook/src/model"
)

type Hook struct {
	MongoClient *mongo.Client
	Model       *model.WebhookClient
	Response    *src.Response
	Logger      *logrus.Logger
	LogModel    *model.LogClient
}

func (h *Hook) HandleHook(c *gin.Context) {
	id := c.Param("id")
	webhook, err := h.Model.GetWebhook(id)
	if err != nil {
		h.Response.Fail(c, err.Error(), nil)
		return
	}
	if webhook.Id.IsZero() || !webhook.Enabled {
		h.Response.NotFound(c, "Webhook not found")
		return
	}

	matchedHook := hook.Hook{
		ID:                    webhook.Id.Hex(),
		Name:                  webhook.Name,
		TriggerRule:           webhook.Triggers,
		Actions:               webhook.Actions,
		PassArgumentsToAction: webhook.PassArgumentsToAction,
	}

	req := &hook.Request{
		RawRequest:  c.Request,
		ContentType: c.ContentType(),
	}
	req.Body, _ = c.GetRawData()

	var headers map[string][]string
	_ = c.ShouldBindHeader(&headers)
	req.ParseHeaders(headers)

	var query map[string][]string
	_ = c.ShouldBindQuery(&query)
	req.ParseQuery(query)

	switch {
	case strings.Contains(req.ContentType, "json"):
		err := req.ParseJSONPayload()
		if err != nil {
			h.Logger.Errorf("Could not parse json payload [%s][%s] %w", matchedHook.ID, matchedHook.Name, err)
		}

	case strings.Contains(req.ContentType, "x-www-form-urlencoded"):
		err := req.ParseFormPayload()
		if err != nil {
			h.Logger.Errorf("Could not parse form payload [%s][%s] %w", matchedHook.ID, matchedHook.Name, err)
		}

	default:
		h.Logger.Errorf("[%s][%s] error parsing body payload due to unsupported content type header: %s", matchedHook.ID, matchedHook.Name, req.ContentType)
	}

	var ok bool
	if matchedHook.TriggerRule == nil {
		ok = true
	} else {
		ok, err = matchedHook.TriggerRule.Evaluate(req)
		if err != nil {
			if !hook.IsParameterNodeError(err) {
				h.Logger.Errorf("Error occurred while evaluating hook rules. %w", err)
				return
			}
			go h.LogModel.AddErrorLog(&matchedHook, fmt.Sprintf("Error occurred while evaluating hook rules. %w", err))
		}
	}

	if ok {
		h.Logger.Debugf("[%s] %s hook triggered successfully", req.ID, matchedHook.ID)
		go h.LogModel.AddLog(&matchedHook, fmt.Sprintf("%s triggered successfully", matchedHook.Name))

		envs, errors := matchedHook.ExtractArgumentsForEnv(req)
		if len(errors) > 0 {
			go h.LogModel.AddWarnLog(&matchedHook, fmt.Sprintf("Error occurred while extracting arguments %s", h.LogModel.ParseErrors(errors)))
		}
		args, errors := matchedHook.ExtractArgumentsAsMap(req)
		if len(errors) > 0 {
			go h.LogModel.AddWarnLog(&matchedHook, fmt.Sprintf("Error occurred while extracting arguments %s", h.LogModel.ParseErrors(errors)))
		}

		for _, act := range *matchedHook.Actions {
			switch act.Driver {
			case hook.ActionShellDriver:
				var actionShell hook.ShellAction
				err := mapstructure.Decode(act.Attributes, &actionShell)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %w", err)
					go h.LogModel.AddErrorLog(&matchedHook, m)
					h.Logger.Error(m)
				}
				shell := action.NewShellAction(&actionShell, &matchedHook, h.LogModel, h.Model)
				go shell.Exec(envs)
			case hook.ActionHttpDriver:
				var actionHttp hook.HttpAction
				err := mapstructure.Decode(act.Attributes, &actionHttp)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %w", err)
					go h.LogModel.AddErrorLog(&matchedHook, m)
					h.Logger.Error(m)
				}
				h := action.NewHttpAction(&actionHttp, &matchedHook, h.LogModel, h.Model)
				go h.Send(args)
			default:
				go h.LogModel.AddWarnLog(&matchedHook, fmt.Sprintf("unsupported action: %s", act.Driver))
			}
		}

		h.Response.Success(c, nil, "OK")
		return
	}

	h.Response.Custom(c, http.StatusNoContent, "NO_CONTENT", nil, "")
}