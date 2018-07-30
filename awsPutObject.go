package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	c "github.com/aws/aws-sdk-go/aws/credentials"
)

type Provider struct {
	C c.Value
}

func (p *Provider) Retrieve() (c.Value, error) {
	return p.C, nil
}

func (p *Provider) IsExpired() bool {
	return true
}

func NewProvider(accessKey, secretAccessKey, sessionToken, providerName string) *Provider {
	return &Provider{
		C: c.Value{
			AccessKeyID: accessKey,
			// AWS Secret Access Key
			SecretAccessKey: secretAccessKey,
			// AWS Session Token
			SessionToken: sessionToken,
			// Provider used to get credentials
			ProviderName: providerName,
		},
	}
}

func main() {
	crd := c.NewCredentials(NewProvider(
		"27qXHOY3xPjKlpNjL0Aiou0hQ+LRnJwooDjSIvqeswP1MdQQwKVvBg==",
		"",
		"my token",
		"wff provider name",
	))

	endpoint := "cn-bj-s3.ufileos.com"
	region := "cn-bj"
	logLevel := aws.LogDebug
	isPathStyle := true
	disableSSL := true
	//contentType := "application/x-www-form-urlencoded"

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      crd,
		Endpoint:         &endpoint,
		Region:           &region,
		LogLevel:         &logLevel,
		S3ForcePathStyle: &isPathStyle,
		DisableSSL:       &disableSSL,
	}))
	svc := s3.New(sess)

	bucket := "wff"
	key := "test_aws_put_object"

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("put object, open ", e)
		return
	}
	defer file.Close()

	s, e := file.Stat()
	if e != nil {
		fmt.Println("put object, stat ", e)
		return
	}

	length := s.Size()
	output, e := svc.PutObject(&s3.PutObjectInput{
		Body:          file,
		Bucket:        &bucket,
		ContentLength: &length,
		Key:           &key,
	})

	if e != nil {
		fmt.Println("put object ", e)
		return
	}

	fmt.Printf("put object output:%s", output.String())
}
