package sns

import (
	"fmt"
	"time"
)

const (
	// SenderTypeSMS 发送短信
	SenderTypeSMS = "sms"
	// SenderTypeEmail 发送邮件
	SenderTypeEmail = "email"
)

// SendCodeConfig 验证码发送
type SendCodeConfig struct {
	Application string                 `json:"application"` // 应用名称
	Provider    string                 `json:"provider"`    // 服务商
	Lang        string                 `json:"lang"`        // 发送语言
	Tag         string                 `json:"tag"`         // 发送模版
	Type        string                 `json:"type"`        // 发送类型: 短信,Email
	From        string                 `json:"from"`        // 发送地址: Email有效
	Contact     string                 `json:"contact"`     // 联系方式
	Expired     int                    `json:"expired"`     // 超时时间,单位:秒
	EnableSame  int                    `json:"enable_same"` // 验证码未超过此时间则使用相同的code,并刷新过期时间
	MaxTry      int                    `json:"max_try"`     // 重试次数
	Data        map[string]interface{} `json:"data"`        // 要发送的数据,Code key 会被覆盖
}

// ContenName 内容名称
func (t *SendCodeConfig) ContenName() string {
	return fmt.Sprintf("%s_%s_%s_content", t.Lang, t.Tag, t.Type)
}

// ExtendName 扩展名称
func (t *SendCodeConfig) ExtendName() string {
	return fmt.Sprintf("%s_%s_%s_extend", t.Lang, t.Tag, t.Type)
}

// CheckCodeConfig 校验Code基础配置
type CheckCodeConfig struct {
	ID   int    `json:"id"`   // 验证码对应的ID
	Code string `json:"code"` // 接收到的Code
}

// TemplateRegisterBase 注册模版基础信息
type TemplateRegisterBase struct {
	Application string            `json:"application"`
	Temps       []*TemplateConfig `json:"temps"`
}

// TemplateConfig 模版注册配置
type TemplateConfig struct {
	Lang    string `json:"lang"`
	Tag     string `json:"tag"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Extend  string `json:"extend"`
}

// ContenName 内容名称
func (t *TemplateConfig) ContenName() string {
	return fmt.Sprintf("%s_%s_%s_content", t.Lang, t.Tag, t.Type)
}

// ExtendName 扩展名称
func (t *TemplateConfig) ExtendName() string {
	return fmt.Sprintf("%s_%s_%s_extend", t.Lang, t.Tag, t.Type)
}

// CodeInfo Code infomation model
type CodeInfo struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Application string    `sql:"index" json:"application"`                         // 应用名称
	Tag         string    `json:"tag"`                                             // 标签
	Contact     string    `gorm:"type:varchar(128);not null;index" json:"contact"` // 邮箱或者手机号
	Code        string    `gorm:"type:varchar(6);not null" json:"code"`            // 验证码根据需求调整长度
	Expiry      time.Time `sql:"index" json:"expiry"`                              // 过期时间
	EnableSame  time.Time `json:"enable_same"`                                     // 验证码未超过此时间则使用相同的code,并刷新过期时间
	TryNum      int       `sql:"index" json:"try_num"`                             // 已尝试次数
	MaxTry      int       `json:"max_try"`                                         // 最大允许尝试次数
}

// MailContentType mail主题类型
type MailContentType string

const (
	// MailContentPlain 文本格式
	MailContentPlain MailContentType = "plain"
	// MailContentHTML HTML格式
	MailContentHTML = "html"
)
