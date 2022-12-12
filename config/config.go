package config

import "github.com/kpango/glg"

var (
	Conf Config
)

type Config struct {
	ListenPort       string `default:":11095" usage:"e.g. :11095" json:"listen_port"`
	AwsDefaultRegion string `default:"ap-northeast-1" usage:"e.g. ap-northeast-1" json:"aws_default_region"`
	AccessKeyID      string `json:"access_key_id"`
	AccessSecret     string `json:"access_secret"`

	TelnetMeInterval int64 `default:"30" json:"telnet_me_interval"`

	CNDefaultIP string `json:"cn_default_ip"`

	// 可用 IP 列表, TODO 不是个配置项
	DefaultIP string `json:"default_ip"`
}

func (c Config) Print() {
	glg.Info("Init Config")
	glg.Info("ListenPort", c.ListenPort)
	glg.Info("AwsDefaultRegion", c.AwsDefaultRegion)
	glg.Info("TelnetMeInterval", c.TelnetMeInterval)
	glg.Info("CNDefaultIP", c.CNDefaultIP)
}
