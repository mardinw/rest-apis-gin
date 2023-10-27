package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sethvargo/go-envconfig"
	"payuoge.com/configs"
	"payuoge.com/pkg/cognito"
	"payuoge.com/pkg/s3"
)

type AwsConnect struct {
	Cognito    *cognito.AwsCognito
	S3         *s3.AwsS3
	GoogleAuth *cognito.GoogleAuthenticator
}

func NewConnect() *AwsConnect {
	ctx := context.Background()

	var configs configs.AppConfiguration

	if err := envconfig.Process(ctx, &configs); err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(configs.AwsConf.AwsRegion),
		config.WithSharedConfigProfile(configs.AwsConf.AwsProfile),
	)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &AwsConnect{
		Cognito:    cognito.NewCognitoClient(&cfg, configs.AwsConf.ClientId, configs.AwsConf.ClientSecret, configs.AwsConf.UserPoolId),
		S3:         s3.NewS3Connect(&cfg),
		GoogleAuth: cognito.NewGoogleAuthenticator(&cfg),
	}
}
