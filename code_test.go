package sns

import "testing"

func Test_Code(t *testing.T) {
	err := InitMysql("root:dijielin@tcp(192.168.0.128:3306)/verification?charset=utf8&parseTime=True&loc=Local", true)
	if err != nil {
		t.Fatalf("Init mysql error:%s\n", err.Error())
	}
	defer CloseMysql()

	AutoMigrate()
	registerSender()
	registerTemplate()
	var cfg = &SendCodeConfig{
		Application: "work-stacks", Provider: "ali", Lang: "cn", Contact: "123456",
		Tag: "system", Type: "sms", MaxTry: 10, Expired: 300}
	u, err := SendCode(cfg)
	if err != nil {
		t.Fatalf("send code error:%s\n", err.Error())
	}
	t.Logf("Code id:%d\n", u)

	ok, err := CheckCode(&CheckCodeConfig{ID: 8, Code: "123456"})
	if err != nil {
		t.Fatalf("check code:%s\n", err.Error())
	}
	t.Logf("Check result:%v\n", ok)
}

type sender struct{}

func (sender) Send(contact, extend, content string) error {
	println(contact, extend, content)
	return nil
}

func registerSender() {
	SenderRegister("ali", "sms", sender{})
}

func registerTemplate() {
	var config = &TemplateRegisterBase{
		Application: "work-stacks",
		Temps: []*TemplateConfig{
			&TemplateConfig{Lang: "en", Tag: "system", Type: "sms", Content: `{"code":"{{.Code}}"}`, Extend: "SMS-1789"},
			&TemplateConfig{Lang: "cn", Tag: "system", Type: "sms", Content: `{"code":"{{.Code}}"}`, Extend: "SMS-1790"},
			&TemplateConfig{Lang: "cn", Tag: "system", Type: "email", Content: `你的验证码为: {{.Code}}`, Extend: `验证码为: {{.Code}}`},
			&TemplateConfig{Lang: "en", Tag: "system", Type: "email", Content: `Hi: {{.Name}}, your's verification code: {{.Code}}`, Extend: `verification code: {{.Code}}`},
		}}
	if err := TemplateRegister(config); err != nil {
		panic(err)
	}
}
