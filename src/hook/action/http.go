package action

import (
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

func (s *Http) Send(args map[string]string) {

}
