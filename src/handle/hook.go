package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slices"
	"strings"
	"webhook/src"
	"webhook/src/hook"
	"webhook/src/hook/action"
	"webhook/src/model"
)

type Hook struct {
	MongoClient *mongo.Client
	Model       *model.WebhookClient
	EmailModel  *model.EmailClient
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
	if webhook.Id.IsZero() {
		h.Response.NotFound(c, "Webhook not found")
		return
	}
	_, err = h.Model.IncreaseCount(id, "callCount")
	if err != nil {
		h.Logger.Errorf("An exception occurred when counting the number of calls. %s", err.Error())
	}

	if !webhook.Enabled {
		h.Response.NotFound(c, "Webhook not found")
		return
	}

	matchedHook := hook.Hook{
		ID:                    webhook.Id.Hex(),
		Name:                  webhook.Name,
		TriggerRule:           webhook.Triggers,
		Actions:               webhook.Actions,
		PassArgumentsToAction: webhook.PassArgumentsToAction,
		Debug:                 webhook.Debug,
	}
	hookLogger := model.NewWebhookLogClient(h.LogModel, &matchedHook)

	auth := c.GetHeader("Authorization")
	if webhook.AuthToken != "" && fmt.Sprintf("hook %s", webhook.AuthToken) != auth {
		go hookLogger.AddWarnLog("Unauthorized call")
		h.Response.Unauthorized(c)
		return
	}

	req := &hook.Request{
		RawRequest:  c.Request,
		ContentType: c.ContentType(),
	}
	req.Body, _ = c.GetRawData()
	req.ParseHeaders(c.Request.Header)
	req.ParseQuery(c.Request.URL.Query())

	if matchedHook.Debug {
		saveReq := &hook.DebugRequest{}
		if slices.Contains(webhook.SaveRequest, "body") {
			saveReq.Body = string(req.Body)
		}
		if slices.Contains(webhook.SaveRequest, "header") {
			saveReq.Headers = c.Request.Header
		}
		if slices.Contains(webhook.SaveRequest, "query") {
			saveReq.Query = c.Request.URL.Query()
		}
		go hookLogger.AddDebugLog(saveReq.ToJson())
	}

	switch {
	case strings.Contains(req.ContentType, "json"):
		err := req.ParseJSONPayload()
		if err != nil {
			h.Logger.Errorf("Could not parse json payload [%s][%s] %s", matchedHook.ID, matchedHook.Name, err.Error())
			go hookLogger.AddErrorLog("Could not parse json payload: " + err.Error())
		}
	case strings.Contains(req.ContentType, "xml"):
		err := req.ParseXMLPayload()
		if err != nil {
			h.Logger.Errorf("Could not parse XML payload [%s][%s] %s", matchedHook.ID, matchedHook.Name, err.Error())
			go hookLogger.AddErrorLog("Could not parse XML payload: " + err.Error())
		}
	case strings.Contains(req.ContentType, "x-www-form-urlencoded"):
		err := req.ParseFormPayload()
		if err != nil {
			h.Logger.Errorf("Could not parse form payload [%s][%s] %s", matchedHook.ID, matchedHook.Name, err.Error())
			go hookLogger.AddErrorLog("Could not parse form payload: " + err.Error())
		}

	default:
		h.Logger.Errorf("[%s][%s] error parsing body payload due to unsupported content type header: %s", matchedHook.ID, matchedHook.Name, req.ContentType)
		go hookLogger.AddWarnLog("error parsing body payload due to unsupported content type header: " + req.ContentType)
	}

	var ok bool
	if matchedHook.TriggerRule == nil {
		ok = true
	} else {
		ok, err = matchedHook.TriggerRule.Evaluate(req)
		if err != nil {
			if !hook.IsParameterNodeError(err) {
				h.Logger.Errorf("Error occurred while evaluating hook rules. %s", err.Error())
				go hookLogger.AddErrorLog(err.Error())
				return
			}
			go hookLogger.AddErrorLog(fmt.Sprintf("Error occurred while evaluating hook rules. %s", err.Error()))
		}
	}

	if ok {
		h.Logger.Debugf("[%s] %s hook triggered successfully", matchedHook.ID, matchedHook.Name)
		go hookLogger.AddLog(fmt.Sprintf("%s triggered successfully", matchedHook.Name))

		envs, errors := matchedHook.ExtractArgumentsForEnv(req)
		if len(errors) > 0 {
			go hookLogger.AddWarnLog(fmt.Sprintf("Error occurred while extracting arguments %s", h.LogModel.ParseErrors(errors)))
		}
		args, errors := matchedHook.ExtractArgumentsAsMap(req)
		if len(errors) > 0 {
			go hookLogger.AddWarnLog(fmt.Sprintf("Error occurred while extracting arguments %s", h.LogModel.ParseErrors(errors)))
		}
		go hookLogger.AddDebugLog(fmt.Sprintf("Parsed envs: %s", envs))
		go hookLogger.AddDebugLog(fmt.Sprintf("Parsed args: %s", args))

		for _, act := range *matchedHook.Actions {
			actionLogger := model.NewActionLogClient(act.Driver, h.LogModel.GenerateLogId(10), h.LogModel, &matchedHook)

			switch act.Driver {
			case hook.ActionShellDriver:
				var actionShell hook.ShellAction
				err := mapstructure.Decode(act.Attributes, &actionShell)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %s", err.Error())
					go hookLogger.AddErrorLog(m)
					h.Logger.Error(m)
				}
				shell := action.NewShellAction(&actionShell, &matchedHook, h.Model, actionLogger)
				go shell.Exec(envs)

			case hook.ActionHttpDriver:
				var actionHttp hook.HttpAction
				err := mapstructure.Decode(act.Attributes, &actionHttp)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %s", err.Error())
					go hookLogger.AddErrorLog(m)
					h.Logger.Error(m)
				}
				h := action.NewHttpAction(&actionHttp, &matchedHook, h.Model, actionLogger)
				go h.Send(args)

			case hook.ActionDispatcherDriver:
				var dispatcher hook.DispatcherAction
				err := mapstructure.Decode(act.Attributes, &dispatcher)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %s", err.Error())
					go hookLogger.AddErrorLog(m)
					h.Logger.Error(m)
				}
				d := action.NewDispatcherAction(&dispatcher, &matchedHook, req, actionLogger)
				go d.Send(args)

			case hook.ActionEmailDriver:
				var email hook.EmailAction
				err := mapstructure.Decode(act.Attributes, &email)
				if err != nil {
					m := fmt.Sprintf("Could not convert action to struct: %s", err.Error())
					go hookLogger.AddErrorLog(m)
					h.Logger.Error(m)
				}
				profile, err := h.EmailModel.GetProfile(email.ProfileId)
				if err != nil {
					go hookLogger.AddErrorLog(fmt.Sprintf("Could not found smtp profile [%s] %s", email.ProfileId, err.Error()))
					return
				}
				e := action.NewEmailAction(&profile, &email, &matchedHook, h.EmailModel, actionLogger)
				go e.Send(args)

			default:
				go hookLogger.AddWarnLog(fmt.Sprintf("unsupported action: %s", act.Driver))
			}
		}

		_, err := h.Model.IncreaseCount(matchedHook.ID, "runCount")
		if err != nil {
			h.Logger.Errorf("An exception occurred when counting the number of runs. %s", err.Error())
		}
		h.Response.Success(c, nil, "OK")
		return
	}

	h.Response.Custom(c, http.StatusNoContent, "NO_CONTENT", nil, "")
}
