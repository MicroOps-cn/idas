/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"path"

	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	*gomail.Dialer
	*gomail.Message
	from string
	to   []string
}

func (clt *SMTPClient) NewClient() *SMTPClient {
	return &SMTPClient{
		gomail.NewDialer(clt.Host, clt.Port, clt.Username, clt.Password),
		gomail.NewMessage(),
		clt.from,
		clt.to,
	}
}

func (clt *SMTPClient) SetTo(to []string) {
	clt.to = to
	clt.SetNativeHeader("To", clt.to...)
}

func (clt *SMTPClient) SetFrom(from string) {
	clt.from = from
	clt.SetAddressHeader("From", clt.from, "")
}

func (clt *SMTPClient) SetSubject(subject string) {
	clt.SetHeader("Subject", "=?UTF-8?B?"+base64.StdEncoding.EncodeToString([]byte(subject))+"?=")
}

func (clt *SMTPClient) Attach(filename string, settings ...gomail.FileSetting) {
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

func (clt *SMTPClient) Send() error {
	if len(clt.GetHeader("From")) == 0 {
		clt.SetAddressHeader("From", clt.from, "")
	}
	if len(clt.GetHeader("To")) == 0 {
		clt.SetNativeHeader("To", clt.to...)
	}
	return clt.DialAndSend(clt.Message)
}

func NewSMTPClient(_ context.Context, options *SmtpOptions) (*SMTPClient, error) {
	if options == nil {
		return nil, fmt.Errorf("smtp options is null")
	}
	if options.Host == "" || options.Username == "" || options.Password == "" {
		return nil, fmt.Errorf("smtp host/username/password is null")
	}
	dialer := gomail.NewDialer(options.Host, int(options.Port), options.Username, options.Password)
	clt, err := dialer.Dial()
	if err != nil {
		return nil, err
	}
	if options.From == "" {
		options.From = options.Username
	}
	_ = clt.Close()
	return &SMTPClient{
		dialer,
		gomail.NewMessage(),
		options.From,
		options.To,
	}, nil
}
