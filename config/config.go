package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/kpango/glg"
)

var (
	Conf Config
)

type Config struct {
	ListenPort       string `default:":11095" usage:"e.g. :11095" json:"listen_port"`
	AwsDefaultRegion string `default:"ap-northeast-1" usage:"e.g. ap-northeast-1" json:"aws_default_region"`
	DefaultPort      int    `default:"10310" json:"default_port"`

	AccessKeyID  string `json:"access_key_id"`
	AccessSecret string `json:"access_secret"`

	TelnetInterval int64 `default:"300" json:"telnet_interval"`
}

func (c Config) Print() {
	glg.Info("Init Config")
	glg.Info("ListenPort", c.ListenPort)
	glg.Info("AwsDefaultRegion", c.AwsDefaultRegion)
	glg.Info("DefaultPort", c.DefaultPort)
	glg.Info("TelnetInterval", c.TelnetInterval)
}

func Init(configFilePath string) {
	loader := aconfig.LoaderFor(&Conf, aconfig.Config{
		SkipEnv:   true,
		SkipFlags: true,
		Files:     []string{configFilePath},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		glg.Fatalln(err)
		return
	}
	Conf.Print()
}
