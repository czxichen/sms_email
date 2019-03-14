package sns

// Template Template
type Template struct {
	Application string `gorm:"primary_key" json:"application"`
	Templates   string `gorm:"type(text)" json:"templates"`
}

// AliSMSRegisterBase Ali SMS Register Base
type AliSMSRegisterBase struct {
	Provider  string `gorm:"primary_key" json:"provider"`
	Host      string `json:"string"`
	AccessID  string `json:"access_id"`
	SecretKey string `json:"secret_key"`
	SignName  string `json:"sign_name"`
}

// IsValid 检测数据是否有效
func (a AliSMSRegisterBase) IsValid() bool {
	return a.Host != "" && a.AccessID != "" && a.SecretKey != "" && a.SignName != ""
}

// AwsSMSRegisterBase aws SMS Register Base
type AwsSMSRegisterBase struct {
	Provider  string `gorm:"primary_key" json:"provider"`
	Region    string `json:"region"`
	Key       string `json:"key"`
	SecretKey string `json:"secret_key"`
}

// IsValid 检测数据是否有效
func (a AwsSMSRegisterBase) IsValid() bool {
	return a.Region != "" && a.Key != "" && a.SecretKey != ""
}

// EmailRegisterBase Email Register Base
type EmailRegisterBase struct {
	Provider string `gorm:"primary_key" json:"provider"`
	Host     string `json:"host"`
	From     string `json:"from"`
	User     string `json:"user"`
	Passwd   string `json:"passwd"`
	IsTLS    bool   `json:"is_tls"`
}

// IsValid 检测数据是否有效
func (e EmailRegisterBase) IsValid() bool {
	return e.Host != "" && e.From != "" && e.User != "" && e.Passwd != ""
}
