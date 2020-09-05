package config

import "time"

type Config struct {
	Trace       bool
	Debug       bool
	Port        int `default:"80"`
	BanTime time.Duration `default:"2m"`
	RequestTime time.Duration `default:"1m"`
	RequestLimit int `default:"100"`
	ZkHosts []string `required:"true" split_words:"true"`
	ZkPrefix string `default:"/limiter" split_words:"true"`
	ZkRequestNode string `default:"req" split_words:"true"`
	ZkBanNode string `default:"ban" split_words:"true"`
}