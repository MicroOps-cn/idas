package email

import (
	"bytes"
	"html/template"
)

func (x *SmtpOptions) getTemplate(topic string, sets ...string) *Template {
	if len(sets) == 0 {
		sets = append(sets, "")
	}
	for _, set := range sets {
		for _, tmpl := range x.Template {
			if tmpl.Topic == topic && tmpl.Set == set {
				return tmpl
			}
		}
	}
	for _, tmpl := range x.Template {
		if tmpl.Topic == topic && (tmpl.Set == "" || tmpl.Set == "__default__") {
			return tmpl
		}
	}
	return nil
}
func (x *SmtpOptions) GetBody(data interface{}, topic string, sets ...string) (body string, err error) {
	t := x.getTemplate(topic, sets...)
	buffer := new(bytes.Buffer)
	if tmpl, err := template.ParseFiles(t.TemplateFile); err != nil {
		return "", err
	} else {
		if err := tmpl.Execute(buffer, data); err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
}

func NewSmtpOptions() *SmtpOptions {
	return &SmtpOptions{}
}
