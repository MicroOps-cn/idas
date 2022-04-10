package ldap

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"strings"
)

type pbLdapOptions LdapOptions

func (p *pbLdapOptions) Reset() {
	(*LdapOptions)(p).Reset()
}

func (p *pbLdapOptions) String() string {
	return (*LdapOptions)(p).String()
}

func (p *pbLdapOptions) ProtoMessage() {
	(*LdapOptions)(p).Reset()
}

func (x *LdapOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewLdapOptions()
	x.UserSearchBase = options.UserSearchBase
	x.UserSearchFilter = options.UserSearchFilter
	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbLdapOptions)(x))
}

// NewLdapOptions return a default option
// which host field point to nowhere.
func NewLdapOptions() *LdapOptions {
	return &LdapOptions{
		Host:              "",
		ManagerDn:         "",
		UserSearchBase:    "",
		GroupSearchBase:   "",
		UserSearchFilter:  "(&(objectClass=inetOrgPerson)(uid={}))",
		GroupSearchFilter: "(&(objectClass=groupOfNames)(uid={}))",
	}
}

func (x *LdapOptions) ParseUserSearchFilter(username ...string) string {
	if len(username) == 0 {
		username = []string{"*"}
	}
	return strings.ReplaceAll(x.UserSearchFilter, "{}", username[0])
}

func (x *LdapOptions) ParseGroupSearchFilter(username ...string) string {
	if len(username) == 0 {
		username = []string{"*"}
	}
	return strings.ReplaceAll(x.GroupSearchFilter, "{}", username[0])
}
