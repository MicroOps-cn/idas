package email

import (
	"bytes"
	"html/template"
)

type EmailTemplate struct {
	Subject      string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	TemplateFile string `protobuf:"bytes,2,opt,name=template_file,json=templateFile,proto3" json:"template_file,omitempty"`
	Topic        string `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	Set          string `protobuf:"bytes,4,opt,name=set,proto3" json:"set,omitempty"`
}

func (t EmailTemplate) Marshal() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
func (t *EmailTemplate) MarshalTo(data []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}
func (t *EmailTemplate) Unmarshal(data []byte) error {
	//TODO implement me
	panic("implement me")
}
func (t *EmailTemplate) Size() int {
	//TODO implement me
	panic("implement me")
}

func (t EmailTemplate) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
func (t *EmailTemplate) UnmarshalJSON(data []byte) error {
	//TODO implement me
	panic("implement me")
}

// only required if the compare option is set
func (t EmailTemplate) Compare(other EmailTemplate) int {
	//TODO implement me
	panic("implement me")
}

// only required if the equal option is set
func (t EmailTemplate) Equal(other EmailTemplate) bool {
	//TODO implement me
	panic("implement me")
}

// https://github.com/gogo/protobuf/blob/master/custom_types.md
// only required if populate option is set
func NewEmailTemplate(r interface{}) *EmailTemplate {
	//TODO implement me
	panic("implement me")

}
func (x *SmtpOptions) getTemplate(topic string, sets ...string) *EmailTemplate {
	if len(sets) == 0 {
		sets = append(sets, "")
	}
	for _, set := range sets {
		for _, tmpl := range x.Template {
			if tmpl.Topic == topic && tmpl.Set == set {
				return &tmpl
			}
		}
	}
	for _, tmpl := range x.Template {
		if tmpl.Topic == topic && (tmpl.Set == "" || tmpl.Set == "__default__") {
			return &tmpl
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
