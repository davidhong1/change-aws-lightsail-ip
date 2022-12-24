package srv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidhong1/change-aws-lightsail-ip/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/kpango/glg"
)

func InitTelnetJob() {
	if config.Conf.CNDefaultIP == "" {
		glg.Info("CNDefaultIP is empty, don't init telnet job")
		return
	}

	go func() {
		t := time.NewTicker(time.Second * time.Duration(config.Conf.TelnetMeInterval))
		defer func() {
			t.Stop()
		}()

		for {
			<-t.C
			glg.Infof("start doTelnetMe %s", config.Conf.CNDefaultIP)
			ok, err := doTelnetMe()
			if err != nil {
				glg.Error(err)
				callChangeIP()
			}
			if ok {
				glg.Info("doTelnetMe ok")
			} else {
				glg.Info("doTelnetMe !ok")
				callChangeIP()
			}
		}
	}()
}

func doTelnetMe() (bool, error) {
	if config.Conf.CNDefaultIP == "" {
		// TODO
		glg.Warn("config.CNDefaultIP is empty")
		return true, nil
	}

	resp, err := http.DefaultClient.Get(
		fmt.Sprintf("http://%s%s/v1/changeawslightsailipsrv/telnetme", config.Conf.CNDefaultIP, config.Conf.ListenPort),
	)
	if err != nil {
		glg.Error(err)
		return false, err
	}
	glg.Debug(resp)

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}

func callChangeIP() error {
	if config.Conf.CNDefaultIP == "" {
		// TODO
		glg.Warn("config.CNDefaultIP is empty")
		return nil
	}

	// 1. reset ip
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://%s%s/v1/changeawslightsailipsrv/ip", config.Conf.CNDefaultIP, config.Conf.ListenPort),
		nil,
	)
	if err != nil {
		glg.Error(err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		glg.Error(err)
		return err
	}
	glg.Debug(resp)

	// 2. restart lightsail
	// get instance id
	resp, err = http.DefaultClient.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		glg.Error(err)
		return err
	}
	glg.Debug(resp)
	body := make([]byte, 0)
	_, err = resp.Body.Read(body)
	if err != nil {
		glg.Error(err)
		return err
	}
	instanceID := string(body)
	// self reboot
	sess, err := newAwsSess()
	if err != nil {
		glg.Error(err)
		return err
	}
	lightsailClient := lightsail.New(sess)
	rebootInstanceOutput, err := lightsailClient.RebootInstance(&lightsail.RebootInstanceInput{
		InstanceName: aws.String(instanceID),
	})
	if err != nil {
		glg.Error(err)
		return err
	}
	glg.Info(rebootInstanceOutput)

	return nil
}

func newAwsSess() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Conf.AwsDefaultRegion),
		Credentials: credentials.NewStaticCredentials(config.Conf.AccessKeyID, config.Conf.AccessSecret, ""),
	})
	return sess, err
}
