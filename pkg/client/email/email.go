package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"gopkg.in/gomail.v2"
	"mime"
	"path"
)

type SmtpClient struct {
	*gomail.Dialer
	*gomail.Message
	from string
	to   []string
}

func (clt *SmtpClient) NewClient() *SmtpClient {
	return &SmtpClient{
		gomail.NewDialer(clt.Host, clt.Port, clt.Username, clt.Password),
		gomail.NewMessage(),
		clt.from,
		clt.to,
	}
}

func (clt *SmtpClient) SetTo(to []string) {
	clt.to = to
	clt.SetNativeHeader("To", clt.to...)
}
func (clt *SmtpClient) SetFrom(from string) {
	clt.from = from
	clt.SetAddressHeader("From", clt.from, "")
}
func (clt *SmtpClient) SetSubject(subject string) {
	clt.SetHeader("Subject", "=?UTF-8?B?"+base64.StdEncoding.EncodeToString([]byte(subject))+"?=")
}
func (clt *SmtpClient) Attach(filename string, settings ...gomail.FileSetting) {
	_, fname := path.Split(filename)
	clt.Message.Attach(filename, append(
		settings,
		gomail.Rename(fname),
		gomail.SetHeader(map[string][]string{
			"Content-Disposition": {
				fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", fname)),
			},
		}),
	)...)
}
func (clt *SmtpClient) Send() error {
	if len(clt.GetHeader("From")) == 0 {
		clt.SetAddressHeader("From", clt.from, "")
	}
	if len(clt.GetHeader("To")) == 0 {
		clt.SetNativeHeader("To", clt.to...)
	}
	return clt.DialAndSend(clt.Message)
}

func NewSmtpClient(_ context.Context, options *SmtpOptions) (*SmtpClient, error) {
	if options == nil {
		return nil, fmt.Errorf("smtp options is null")
	}
	if options.Host == "" || options.Username == "" || options.Password == "" {
		return nil, fmt.Errorf("smtp host/username/password is null")
	}
	dialer := gomail.NewDialer(options.Host, int(options.Port), options.Username, options.Password)
	if clt, err := dialer.Dial(); err != nil {
		return nil, err
	} else {
		if options.From == "" {
			options.From = options.Username
		}
		_ = clt.Close()
		return &SmtpClient{
			dialer,
			gomail.NewMessage(),
			options.From,
			options.To,
		}, nil
	}
}
