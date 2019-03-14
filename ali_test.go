package sns

import (
	"testing"
)

func Test_Email(t *testing.T) {
	mail, err := NewAliEmail("smtpdm-ap-southeast-1.aliyun.com:80", "no-reply@work-stacks.com", "no-reply@work-stacks.com", "YBAIwOKFRUh6F7wg1", true)
	if err != nil {
		t.Fatalf("Init ali email error:%s\n", err.Error())
	}
	if err := mail.Send("dijielin@qq.com", "This is test mail", "Hello world"); err != nil {
		t.Fatalf("Send mail error:%s\n", err.Error())
	}
}
