package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Config struct {
	SyncerImpl string `yaml:"syncer_impl" validate:"required,oneof=rsync"`
	Rsync      struct {
		// Mandatory options
		Host      string `yaml:"host" validate:"required"`
		LocalDir  string `yaml:"local_dir" validate:"required,endswith=/"`
		RemoteDir string `yaml:"remote_dir" validate:"required,endswith=/"`

		// Advanced options
		BandwidthLimit    string `yaml:"bwlimit" validate:"excludes= "`
		ExcludePattern    string `yaml:"exclude"`
		RemoveSourceFiles bool   `yaml:"remove_source_files"`
		RemoteShell       string `yaml:"remote_shell"`
	} `yaml:"rsync"`

	Hooks []HookConfigContainer `yaml:"hooks"`

	EventSourceImpl []string `yaml:"events_impl" validate:"omitempty,dive,oneof=kafka webhook_server ticker"`
	Kafka           struct {
		// Mandatory options
		Brokers []string `yaml:"brokers" validate:"dive,required_if=EventSourceImpl kafka"`
		Topic   string   `yaml:"topic" validate:"required_if=EventSourceImpl kafka"`
		GroupId string   `yaml:"group_id" validate:"required_if=EventSourceImpl kafka"`

		// Advanced options
		Partition   int    `yaml:"partition" validate:"gte=0"`
		TlsCertFile string `yaml:"tls_cert_file" validate:"omitempty,file"`
		TlsKeyFile  string `yaml:"tls_key_file" validate:"omitempty,file"`
	} `yaml:"kafka"`

	Webhook struct {
		// Mandatory options
		Address string `yaml:"address" validate:"required_if=EventSourceImpl webhook_server"`

		// Advanced options
		Path        string `yaml:"path" validate:"omitempty,startswith=/"`
		TlsCertFile string `yaml:"tls_cert_file" validate:"required_unless=TlsKeyFile '',omitempty,filepath"`
		TlsKeyFile  string `yaml:"tls_key_file" validate:"required_unless=TlsCertFile '',omitempty,filepath"`
	} `yaml:"webhook_server"`

	Ticker struct {
		IntervalSeconds int `yaml:"interval_s" validate:"required_if=EventSourceImpl ticker,gt=300"`
	} `yaml:"ticker"`

	MetricsAddr string `yaml:"metrics_addr" validate:"omitempty,tcp_addr"`
}

func Validate(s any) error {
	return validate.Struct(s)
}

func Read(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(content, conf)
	return conf, err
}
