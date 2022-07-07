package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"html/template"
)

type Template struct {
	*OriginalTemplate
	tmpl *template.Template
}

func (t Template) Marshal() ([]byte, error) {
	return proto.Marshal(t.OriginalTemplate)
}

func (t *Template) Unmarshal(data []byte) (err error) {
	if t.OriginalTemplate == nil {
		t.OriginalTemplate = &OriginalTemplate{}
	}
	if err = proto.Unmarshal(data, t.OriginalTemplate); err != nil {
		return err
	}
	t.tmpl, err = template.ParseFiles(t.TemplateFile)
	if err != nil {
		return err
	}
	return nil
}

func (t Template) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.OriginalTemplate)
}
func (t *Template) UnmarshalJSON(data []byte) (err error) {
	if t.OriginalTemplate == nil {
		t.OriginalTemplate = &OriginalTemplate{}
	}
	if err = json.Unmarshal(data, &t.OriginalTemplate); err != nil {
		return err
	}
	t.tmpl, err = template.ParseFiles(t.TemplateFile)
	if err != nil {
		return err
	}
	return
}

func (m *SmtpOptions) getTemplate(topic string, sets ...string) *Template {
	if len(sets) == 0 {
		sets = append(sets, "")
	}
	for _, set := range sets {
		for _, tmpl := range m.Template {
			if tmpl.Topic == topic && tmpl.Set == set {
				return &tmpl
			}
		}
	}
	for _, tmpl := range m.Template {
		if tmpl.Topic == topic && (tmpl.Set == "" || tmpl.Set == "__default__") {
			return &tmpl
		}
	}
	return nil
}
func (m *SmtpOptions) GetSubjectAndBody(data interface{}, topic string, sets ...string) (subject, body string, err error) {
	t := m.getTemplate(topic, sets...)
	buffer := new(bytes.Buffer)
	if t.tmpl == nil {
		return "", "", fmt.Errorf("template is nil")
	}
	if err = t.tmpl.Execute(buffer, data); err != nil {
		return "", "", err
	}
	return t.GetSubject(), buffer.String(), nil
}

func NewSmtpOptions() *SmtpOptions {
	return &SmtpOptions{}
}
