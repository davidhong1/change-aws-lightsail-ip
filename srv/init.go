package srv

import (
	"fmt"
	"io/ioutil"
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
		// don't return, try to restart
	} else {
		glg.Debug(resp)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// 2. restart lightsail
	// get instance localIPv4
	resp, err = http.DefaultClient.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		glg.Error(err)
		return err
	}
	glg.Debug(resp)
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glg.Error(err)
		return err
	}
	localIPv4 := string(s)
	glg.Debug(localIPv4)

	// self reboot
	sess, err := newAwsSess()
	if err != nil {
		glg.Error(err)
		return err
	}
	// list Lightsail instance to found instanceName which match localIPv4
	instances, err := listLightsailInstance(sess)
	if err != nil {
		glg.Error(err)
		return err
	}
	var instanceName string
	for _, instance := range instances {
		if *instance.PrivateIpAddress == localIPv4 {
			instanceName = *instance.Name
			break
		}
	}
	if instanceName == "" {
		err = fmt.Errorf("instanceName is empty")
		glg.Error(err)
		return err
	}

	lightsailClient := lightsail.New(sess)
	rebootInstanceOutput, err := lightsailClient.RebootInstance(&lightsail.RebootInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	if err != nil {
		glg.Error(err)
		return err
	}
	glg.Info(rebootInstanceOutput)

	return nil
}

func listLightsailInstance(sess *session.Session) ([]*lightsail.Instance, error) {
	lightsailClient := lightsail.New(sess)

	var err error
	var nextToken *string
	var resp *lightsail.GetInstancesOutput
	var out []*lightsail.Instance
	for {
		if nextToken == nil {
			resp, err = lightsailClient.GetInstances(&lightsail.GetInstancesInput{})
		} else {
			resp, err = lightsailClient.GetInstances(&lightsail.GetInstancesInput{
				PageToken: nextToken,
			})
		}
		if err != nil {
			return nil, err
		}
		out = append(out, resp.Instances...)
		if resp.NextPageToken == nil {
			break
		}
		nextToken = resp.NextPageToken
	}

	return out, nil
}

func newAwsSess() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Conf.AwsDefaultRegion),
		Credentials: credentials.NewStaticCredentials(config.Conf.AccessKeyID, config.Conf.AccessSecret, ""),
	})
	return sess, err
}
