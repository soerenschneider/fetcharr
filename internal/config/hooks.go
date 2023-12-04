package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	WebhookType = "webhook"
	CmdHookType = "cmd"
)

type HookConfig interface {
	GetType() string
	GetName() string
	GetStage() Stage
}

type WebhookConfig struct {
	Name            string         `yaml:"name" validate:"required"`
	Stage           Stage          `yaml:"stage"`
	Endpoint        string         `yaml:"endpoint" validate:"url"`
	Verb            string         `yaml:"verb" validate:"omitempty,oneof=GET get POST post PATCH patch HEAD head"`
	Data            map[string]any `yaml:"data,omitempty" validate:"excluded_with=EncodedData"`
	EncodedDataFile string         `yaml:"encoded_data_file,omitempty" validate:"excluded_with=Data,filepath"`
}

func (w *WebhookConfig) GetName() string {
	return w.Name
}

func (w *WebhookConfig) GetStage() Stage {
	return w.Stage
}

func (w *WebhookConfig) GetType() string {
	return WebhookType
}

type CmdHookConfig struct {
	Name        string   `yaml:"name" validate:"required"`
	Stage       Stage    `yaml:"stage" validate:"required,oneof=PRE POST_SUCCESS POST_SUCCESS_TRANSFER POST_FAILURE"`
	Cmds        []string `yaml:"commands" validate:"required"`
	StopOnError *bool    `yaml:"stop_on_error"`
}

func (w *CmdHookConfig) GetStage() Stage {
	return w.Stage
}

func (w *CmdHookConfig) GetName() string {
	return w.Name
}

func (w *CmdHookConfig) GetType() string {
	return CmdHookType
}

type HookConfigContainer struct {
	HookConfig HookConfig
}

func (c *HookConfigContainer) UnmarshalYAML(node *yaml.Node) error {
	type inner struct {
		Type string `yaml:"type"`
	}

	hookType := &inner{}
	if err := node.Decode(hookType); err != nil {
		return err
	}

	switch hookType.Type {
	case WebhookType:
		hook := &WebhookConfig{}
		if err := node.Decode(hook); err != nil {
			return err
		}
		c.HookConfig = hook
	case CmdHookType:
		hook := getDefaultCmdHookConfig()
		if err := node.Decode(hook); err != nil {
			return err
		}
		c.HookConfig = hook
	default:
		return fmt.Errorf("unknown hook type %q", hookType.Type)
	}

	if err := Validate(c.HookConfig); err != nil {
		return err
	}

	return nil
}

func getDefaultCmdHookConfig() *CmdHookConfig {
	stopOnError := true
	return &CmdHookConfig{
		StopOnError: &stopOnError,
	}
}
