package main

import (
	"fmt"
	//"os"

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
	k1 := "嘿嘿/New Folder/"
	k2 := "go-eco.jpg"
	k3 := "sss"

	del := &s3.Delete{}
	del.Objects = make([]*s3.ObjectIdentifier, 3)
	del.Objects[0] = &s3.ObjectIdentifier{
		Key: &k1,
	}
	del.Objects[1] = &s3.ObjectIdentifier{
		Key: &k2,
	}
	del.Objects[2] = &s3.ObjectIdentifier{
		Key: &k3,
	}

	output, e := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: del,
	})

	if e != nil {
		fmt.Println("delete objects", e)
		return
	}

	fmt.Printf("delte objects output:%s", output.String())
}
