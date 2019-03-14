package sns

import (
	"testing"
)

func Test_RegisterTemplate(t *testing.T) {
	var config = &TemplateRegisterBase{
		Application: "Work-stacks",
		Temps: []*TemplateConfig{
			&TemplateConfig{Lang: "en", Type: "sms", Content: `{"code":"{{.Code}}"}`, Extend: "SMS-1789"},
			&TemplateConfig{Lang: "cn", Type: "sms", Content: `{"code":"{{.Code}}"}`, Extend: "SMS-1790"},
			&TemplateConfig{Lang: "cn", Type: "email", Content: `你的验证码为: {{.Code}}`, Extend: `验证码为: {{.Code}}`},
			&TemplateConfig{Lang: "en", Type: "email", Content: `Hi: {{.Name}}, your's verification code: {{.Code}}`, Extend: `verification code: {{.Code}}`},
		}}
	if err := TemplateRegister(config); err != nil {
		t.Fatal(err)
	}

	var cfg = &SendCodeConfig{
		Application: "Work-stacks",
		Type:        "sms", Lang: "en",
		Data: map[string]interface{}{"Code": "123456"},
	}

	s, c, e := TemplateToMsg(cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("\nExtend: %s\nContent: %s\n", s, c)
	cfg.Lang = "cn"
	s, c, e = TemplateToMsg(cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("\nExtend: %s\nContent: %s\n", s, c)
	cfg.Type = "email"
	s, c, e = TemplateToMsg(cfg)
	if e != nil {
		t.Fatal(e)
	}
	cfg.Lang = "en"
	cfg.Data["Name"] = "czxichen"
	t.Logf("\nExtend: %s\nContent: %s\n", s, c)
	s, c, e = TemplateToMsg(cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("\nExtend: %s\nContent: %s\n", s, c)
}
