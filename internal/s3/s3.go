package s3

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(cfg)
}

func Upload(fileName string, file io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String("go-image-util"),
		Key: aws.String(fileName),
		Body: file,
	}

	_, err := client.PutObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
