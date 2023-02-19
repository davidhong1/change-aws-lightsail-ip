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
			ctx := context.Background()
			err := doTelnet(ctx)
			if err != nil {
				glg.Error(err)
			}

			removeUnailableIP(ctx)
		}
	}()
}

func doTelnet(ctx context.Context) error {
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

	glg.Infof("telnet %d instance", len(instances))

	for _, instance := range instances {
		if *instance.State.Name == "stopped" {
			// filter stopped instance
			err = startInstance(ctx, *instance.Name)
			if err != nil {
				glg.Error(err)
			}
			continue
		}

		if instance.PublicIpAddress == nil || *instance.State.Name != "running" {
			continue
		}

		now := time.Now()
		defer func() {
			glg.Infof("telnet %s %s use time %f", *instance.Name, *instance.PublicIpAddress, time.Since(now).Seconds())
		}()

		_, err := telnet.TelnetWithRetry(ctx, *instance.PublicIpAddress, config.Conf.DefaultPort, config.Conf.TelnetRetryTime)
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
	err := stopInstance(ctx, instanceName)
	if err != nil {
		glg.Error(err)
		return err
	}

	time.Sleep(3 * time.Minute)

	err = startInstance(ctx, instanceName)
	if err != nil {
		glg.Error(err)
		return err
	}

	return nil
}

func stopInstance(ctx context.Context, instanceName string) error {
	sess, err := awsls.NewAwsSess(config.Conf.AwsDefaultRegion, config.Conf.AccessKeyID, config.Conf.AccessSecret)
	if err != nil {
		glg.Error(err)
		return err
	}
	err = awsls.StopInstance(ctx, sess, &instanceName)
	if err != nil {
		glg.Error(err)
	}
	return nil
}

func startInstance(ctx context.Context, instanceName string) error {
	sess, err := awsls.NewAwsSess(config.Conf.AwsDefaultRegion, config.Conf.AccessKeyID, config.Conf.AccessSecret)
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

func removeUnailableIP(ctx context.Context) {
	for _, v := range iplist.Get() {
		_, err := telnet.TelnetWithRetry(ctx, v, config.Conf.DefaultPort, config.Conf.TelnetRetryTime)
		if err != nil {
			glg.Error(err)

			// remove ip from array
			iplist.Remove(v)
		}
	}
}
