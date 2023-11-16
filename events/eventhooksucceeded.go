package events

import (
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/hooks"
	"github.com/captainhook-go/captainhook/hooks/app"
)

type HookSucceeded struct {
	Context *app.Context
	Config  *configuration.Hook
	Log     *hooks.ActionLog
}

func NewHookSucceededEvent(context *app.Context, hook *configuration.Hook, log *hooks.ActionLog) *HookSucceeded {
	e := HookSucceeded{
		Context: context,
		Config:  hook,
		Log:     log,
	}
	return &e
}

type HookSucceededSubscriber interface {
	Handle(event *HookSucceeded) error
}
