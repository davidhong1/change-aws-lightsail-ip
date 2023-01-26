package srv

import (
	"context"
	"time"

	"github.com/davidhong1/change-aws-lightsail-ip/config"
	"github.com/davidhong1/change-aws-lightsail-ip/internal/iplist"
	"github.com/davidhong1/change-aws-lightsail-ip/pkg/awsls"
	"github.com/davidhong1/change-aws-lightsail-ip/pkg/telnet"

	"github.com/kpango/glg"
)

func InitTelnetJob() {
	go func() {
		t := time.NewTicker(time.Second * time.Duration(config.Conf.TelnetInterval))
		defer func() {
			t.Stop()
		}()

		for {
			<-t.C
			err := doTelnet()
			if err != nil {
				glg.Error(err)
			}
		}
	}()
}

func doTelnet() error {
	ctx := context.Background()

	sess, err := awsls.NewAwsSess(config.Conf.AwsDefaultRegion, config.Conf.AccessKeyID, config.Conf.AccessSecret)
	if err != nil {
		glg.Error(err)
		return err
	}

	instances, err := awsls.ListInstances(ctx, sess)
	if err != nil {
		glg.Error(err)
		return err
	}
	for _, instance := range instances {
		now := time.Now()
		defer func() {
			glg.Info("telnet %s %s use time %f", *instance.Name, *instance.PublicIpAddress, time.Since(now).Seconds())
		}()

		_, err := telnet.Telnet(ctx, *instance.PublicIpAddress, config.Conf.DefaultPort)
		if err != nil {
			glg.Error(err)

			// remove ip from array
			iplist.Remove(*instance.PublicIpAddress)

			go stopAndStartInstance(ctx, *instance.Name)

			continue
		}

		// add can telnet ip
		iplist.Add(*instance.PublicIpAddress)
	}

	return nil
}

func stopAndStartInstance(ctx context.Context, instanceName string) error {
	sess, err := awsls.NewAwsSess(config.Conf.AwsDefaultRegion, config.Conf.AccessKeyID, config.Conf.AccessSecret)
	if err != nil {
		glg.Error(err)
		return err
	}
	err = awsls.StopInstance(ctx, sess, &instanceName)
	if err != nil {
		glg.Error(err)
		return err
	}

	time.Sleep(2 * time.Minute)

	sess, err = awsls.NewAwsSess(config.Conf.AwsDefaultRegion, config.Conf.AccessKeyID, config.Conf.AccessSecret)
	if err != nil {
		glg.Error(err)
		return err
	}
	err = awsls.StartInstance(ctx, sess, &instanceName)
	if err != nil {
		glg.Error(err)
		return err
	}
	return nil
}
