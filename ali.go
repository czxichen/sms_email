package sns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// NewAliSMS 创建一个sms对象,用来发送阿里云短信
func NewAliSMS(host, accessid, secretkey, signname string) Sender {
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	return &aliSMS{host: host, secretKey: secretkey + "&",
		parasPool: &sync.Pool{New: func() interface{} {
			return url.Values{
				"SignatureMethod":  {"HMAC-SHA1"},
				"SignatureVersion": {"1.0"},
				"Format":           {"JSON"},
				"Action":           {"SendSms"},
				"Version":          {"2017-05-25"},
				"RegionId":         {"cn-hangzhou"},
				"SignName":         {signname},
				"AccessKeyId":      {accessid},

				"TemplateCode":   {""}, // SMS_137656329
				"SignatureNonce": {""}, // uuid.New().String()
				"PhoneNumbers":   {""}, // 17000000000
				"Timestamp":      {""}, // time.Now().UTC().Format("2006-01-02T15:04:05Z")
				"TemplateParam":  {""}, // `{"code":"123456"}`
			}
		}}}
}

type aliSMS struct {
	host      string
	secretKey string
	parasPool *sync.Pool
}

// 请求参数签名
func (base *aliSMS) sign(str string) string {
	hash := hmac.New(sha1.New, []byte(base.secretKey))
	hash.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// 编码请求参数
func (base *aliSMS) encodeToURI(phoneNumbers, content, extend string) string {
	_paras := base.parasPool.Get().(url.Values)

	_paras["SignatureNonce"][0] = uuid.New().String()
	_paras["Timestamp"][0] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	_paras["TemplateParam"][0] = content
	_paras["PhoneNumbers"][0] = phoneNumbers
	_paras["TemplateCode"][0] = extend
	paramstr := _paras.Encode()
	base.parasPool.Put(_paras)

	stringToSign := fmt.Sprintf("GET&%s&%s", "%2F", url.QueryEscape(paramstr))
	return fmt.Sprintf("Signature=%s&%s", url.QueryEscape(base.sign(stringToSign)), paramstr)
}

// phoneNumbers:接收号码,国际号码: '00 + 国家代码 + 手机号'
// code:json格式的字符串,结合模版使用,`{"code":"123456"}`
// extend:使用的模版ID
func (base *aliSMS) Send(phoneNumbers, extend, content string) error {
	uri := base.encodeToURI(phoneNumbers, content, extend)

	resp, err := http.Get(base.host + "?" + uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respCode := &aliSMSReply{}
	if err := json.Unmarshal(body, respCode); err != nil {
		return err
	}
	if respCode.Code != "OK" {
		return errors.New(respCode.Code)
	}
	return nil
}

// aliSMSReply 响应格式
type aliSMSReply struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

// NewAliEmail 阿里发送接口
func NewAliEmail(host, from, user, passwd string, isTLS bool) (Sender, error) {
	return SMTPEMail(host, from, user, passwd, isTLS)
}
