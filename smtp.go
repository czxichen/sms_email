package sns

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net/smtp"
	"net/textproto"
	"strings"
)

// SMTPEMail SMTP EMail
func SMTPEMail(host, from, user, passwd string, isTLS bool) (Sender, error) {
	em := &email{Host: host, User: user, Passwd: passwd, From: from, IsTLS: isTLS}
	client, err := em.dial()
	if err != nil {
		return nil, err
	}
	client.Close()
	return em, nil
}

type email struct {
	Host   string // 邮箱smtp地址
	User   string
	Passwd string
	From   string
	IsTLS  bool // 表示是否启用tls连接
}

func (email *email) dial() (*smtp.Client, error) {
	client, err := smtp.Dial(email.Host)
	if err != nil {
		return nil, err
	}
	if email.IsTLS {
		if err = client.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			client.Close()
			return nil, err
		}
	}
	if err = client.Auth(smtp.PlainAuth("", email.User, email.Passwd, strings.Split(email.Host, ":")[0])); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

// 发送封装好的消息体
func (email *email) Send(contact, subject, content string) error {
	client, err := email.dial()
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Mail(email.From); err != nil {
		return err
	}

	if err = client.Rcpt(contact); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}
	// 默认使用的是HTML格式,如果是纯文本则使用:MailContentPlain
	info := &mailInfo{Type: MailContentHTML, To: []string{contact}, From: email.From, Subject: subject, Content: content}
	if err = info.Writer(wc); err != nil {
		wc.Close()
	}
	err = client.Quit()
	if err != nil {
		if strings.Index(err.Error(), "250") == 0 {
			return nil
		}
	}
	return err
}

type mailInfo struct {
	Type    MailContentType // 邮件文本类型
	To      []string        // 接收地址列表
	From    string          // 发送地址
	Subject string          // 邮件主题
	Content string          // 邮件内容
}

func (info *mailInfo) Headers() textproto.MIMEHeader {
	res := make(textproto.MIMEHeader)
	if _, ok := res["To"]; !ok && len(info.To) > 0 {
		res.Set("To", strings.Join(info.To, ","))
	}

	if _, ok := res["Subject"]; !ok && info.Subject != "" {
		res.Set("Subject", info.Subject)
	}

	if _, ok := res["From"]; !ok {
		res.Set("From", info.From)
	}
	return res
}

func (info *mailInfo) Writer(datawriter io.Writer) error {
	headers := info.Headers()
	headers.Set("Content-Type", fmt.Sprintf("text/%s; charset=UTF-8", info.Type))
	headers.Set("Content-Transfer-Encoding", "quoted-printable")
	headerToBytes(datawriter, headers)
	fmt.Fprintf(datawriter, "\r\n\r\n%s\r\n", info.Content)
	return nil
}

// 序列化头信息
func headerToBytes(w io.Writer, header textproto.MIMEHeader) {
	for field, vals := range header {
		for _, subval := range vals {
			io.WriteString(w, field)
			io.WriteString(w, ": ")
			switch {
			case field == "Content-Type" || field == "Content-Disposition":
				w.Write([]byte(subval))
			default:
				w.Write([]byte(mime.QEncoding.Encode("UTF-8", subval)))
			}
			io.WriteString(w, "\r\n")
		}
	}
}
