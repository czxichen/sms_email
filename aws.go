package sns

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// NewAwsSMS 创建Aws短信接口
func NewAwsSMS(region, key, sec string) (Sender, error) {
	s := session.New(&aws.Config{Region: &region, Credentials: credentials.NewStaticCredentials(key, sec, "")})
	snsClient := sns.New(s)
	_, err := snsClient.GetSMSAttributes(new(sns.GetSMSAttributesInput))
	return &awsSMS{client: snsClient, pool: new(sync.Pool)}, err
}

type awsSMS struct {
	client *sns.SNS
	pool   *sync.Pool
}

func (a *awsSMS) getPublishInput(contact, content string) *sns.PublishInput {
	input, ok := a.pool.Get().(*sns.PublishInput)
	if !ok {
		input = &sns.PublishInput{}
	}
	input.Message = &content
	input.PhoneNumber = &contact
	return input
}

func (a *awsSMS) putPublishInput(input *sns.PublishInput) {
	a.pool.Put(input)
}

func (a *awsSMS) Send(contact, extend, content string) error {
	input := a.getPublishInput(contact, content)
	_, err := a.client.Publish(input)
	a.putPublishInput(input)
	return err
}

// NewAwsEmail aws email 服务
func NewAwsEmail(host, from, key, sec string, isTLS bool) (Sender, error) {
	h := hmac.New(sha256.New, []byte(sec))
	h.Write([]byte("SendRawEmail"))
	buf := h.Sum(nil)
	version := []byte{0x02}
	version = append(version, buf...)
	passwd := base64.RawStdEncoding.EncodeToString(version)

	return SMTPEMail(host, from, key, passwd, isTLS)
}
