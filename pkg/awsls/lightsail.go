package awsls

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/kpango/glg"
)

func StopInstance(ctx context.Context, sess *session.Session, instanceName *string) error {
	lightsailClient := lightsail.New(sess)

	output, err := lightsailClient.StopInstance(&lightsail.StopInstanceInput{
		Force:        aws.Bool(true),
		InstanceName: instanceName,
	})
	if err != nil {
		glg.Error(err)
		return nil
	}
	glg.Info(output.Operations)

	return nil
}

func StartInstance(ctx context.Context, sess *session.Session, instanceName *string) error {
	lightsailClient := lightsail.New(sess)

	output, err := lightsailClient.StartInstance(&lightsail.StartInstanceInput{
		InstanceName: instanceName,
	})
	if err != nil {
		glg.Error(err)
		return nil
	}
	glg.Info(output.Operations)

	return nil
}

func ListInstances(ctx context.Context, sess *session.Session) ([]*lightsail.Instance, error) {
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

func NewAwsSess(region, accessKeyID, accessSecret string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, accessSecret, ""),
	})
	return sess, err
}
