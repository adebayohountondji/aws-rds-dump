package internal

import (
	awsClient "github.com/aws/aws-sdk-go/aws"
	awsClientCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	awsClientSession "github.com/aws/aws-sdk-go/aws/session"
	awsClientS3 "github.com/aws/aws-sdk-go/service/s3"
	"os"
)

type S3 struct {
	bucketName string
	client     *awsClientS3.S3
}

func (s3 *S3) PutObject(objectKey string, file *os.File) (err error) {
	_, err = s3.client.PutObject(
		&awsClientS3.PutObjectInput{
			Bucket: &s3.bucketName,
			Key:    &objectKey,
			Body:   file,
		},
	)

	return err
}

type AwsConfig struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string
}

type Aws struct {
	clientSession *awsClientSession.Session
}

func NewAws(config AwsConfig) (aws Aws, err error) {
	var awsClientConfig []*awsClient.Config

	if config.AccessKeyId != "" || config.SecretAccessKey != "" {
		awsClientConfig = append(
			awsClientConfig,
			&awsClient.Config{
				Credentials: awsClientCredentials.NewStaticCredentials(
					config.AccessKeyId,
					config.SecretAccessKey,
					"",
				),
			},
		)
	}

	if config.Region != "" {
		awsClientConfig = append(
			awsClientConfig,
			&awsClient.Config{
				Region: &config.Region,
			},
		)
	}

	session, err := awsClientSession.NewSession(awsClientConfig...)

	if err == nil {
		aws = Aws{clientSession: session}
	}

	return aws, err
}

func (aws *Aws) NewS3(bucketName string) S3 {
	return S3{
		bucketName: bucketName,
		client:     awsClientS3.New(aws.clientSession),
	}
}
