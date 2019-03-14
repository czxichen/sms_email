package sns

import "fmt"

// Handler handler
type Handler struct{}

// SendCode 发送验证码
func (h *Handler) SendCode(config SendCodeConfig, resp *Response) error {
	resp.ID, resp.Error = SendCode(&config)
	return nil
}

// CheckCode 校验验证码
func (h *Handler) CheckCode(config CheckCodeConfig, resp *Response) error {
	resp.Status, resp.Error = CheckCode(&config)
	return nil
}

// RegisterTemplate 注册模版
func (h *Handler) RegisterTemplate(config TemplateRegisterBase, resp *Response) error {
	if err := TemplateRegister(&config); err != nil {
		resp.Error = err
	}
	return nil
}

// RegisterAliSMSSender Register ali sms Sender
func (h *Handler) RegisterAliSMSSender(config AliSMSRegisterBase, resp *Response) error {
	if config.IsValid() {
		SenderRegister(config.Provider, SenderTypeSMS, NewAliSMS(config.Host, config.AccessID, config.SecretKey, config.SignName))
	} else {
		resp.Error = fmt.Errorf("invalid data:%+v", config)
	}
	return nil
}

// RegisterAwsSMSSender Register aws sms Sender
func (h *Handler) RegisterAwsSMSSender(config AwsSMSRegisterBase, resp *Response) error {
	if config.IsValid() {
		sender, err := NewAwsSMS(config.Region, config.Key, config.SecretKey)
		if err == nil {
			SenderRegister(config.Provider, SenderTypeSMS, sender)
		}
		resp.Error = err
	} else {
		resp.Error = fmt.Errorf("invalid data:%+v", config)
	}
	return nil
}

// RegisterAliEMailSender Register ali email Sender
func (h *Handler) RegisterAliEMailSender(config EmailRegisterBase, resp *Response) error {
	return h.registerEmailSender(config, resp, NewAliEmail)
}

// RegisterAwsEMailSender Register aws email Sender
func (h *Handler) RegisterAwsEMailSender(config EmailRegisterBase, resp *Response) error {
	return h.registerEmailSender(config, resp, NewAwsEmail)
}

func (h *Handler) registerEmailSender(config EmailRegisterBase, resp *Response, newSender func(string, string, string, string, bool) (Sender, error)) error {
	if config.IsValid() {
		sender, err := newSender(config.Host, config.From, config.User, config.Passwd, config.IsTLS)
		if err == nil {
			SenderRegister(config.Provider, SenderTypeEmail, sender)
		}
		resp.Error = err
	} else {
		resp.Error = fmt.Errorf("invalid data:%+v", config)
	}
	return nil
}

// Response 响应数据
type Response struct {
	ID     uint  `json:"id,omitempty"`
	Status bool  `json:"status,omitempty"`
	Error  error `json:"error,omitempty"`
}
