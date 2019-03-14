package sns

import (
	"testing"
)

const (
	region = "us-west-2"
	key    = `XKIAIJXB6X3QUOSRPYUQ`
	sec    = `e6q5c3a1PC2hHhdSoLa28zXAdGMmSeuT4tNxcMhU`
	phone  = "008617000000000"
)

func Test_aws_sms(t *testing.T) {
	s, err := NewAwsSMS(region, key, sec)
	if err != nil {
		t.Fatalf("New aws sms error:%s\n", err.Error())
	}
	if err = s.Send(phone, "", "Hello world"); err != nil {
		t.Fatalf("Sender error:%s\n", err.Error())
	}
}

const (
	from      = "no-repley@work-stacks.com"
	to        = "dijielin@qq.com"
	user      = "XKIAIJXB6X3QUOSRPYUQ"
	secPasswd = "e6q5c3a1PC2hHhdSoLa28zXAdGMmSeuT4tNxcMhU"
)

func Test_aws_email(t *testing.T) {
	s, err := NewAwsEmail("email-smtp.us-west-2.amazonaws.com:25", from, user, secPasswd, true)
	if err != nil {
		t.Fatalf("init aws email error:%s\n", err.Error())
	}
	if err = s.Send(to, "This is test mail", "Test"); err != nil {
		t.Fatalf("panic error：%s\n", err.Error())
	}
}
