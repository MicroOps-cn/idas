package ldap

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
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
	x.Host = options.Host
	x.ManagerDn = options.ManagerDn
	x.UserSearchBase = options.UserSearchBase
	x.UserSearchFilter = options.UserSearchFilter
	x.GroupSearchBase = options.GroupSearchBase
	x.GroupSearchFilter = options.GroupSearchFilter
	x.AttrEmail = options.AttrEmail
	x.AttrUsername = options.AttrUsername
	x.AttrUserDisplayName = options.AttrUserDisplayName
	x.AttrUserPhoneNo = options.AttrUserPhoneNo
	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbLdapOptions)(x))
}

// NewLdapOptions return a default option
// which host field point to nowhere.
func NewLdapOptions() *LdapOptions {
	return &LdapOptions{
		Host:                "127.0.0.1:389",
		ManagerDn:           "cn=admin,dc=example,dc=org",
		UserSearchBase:      "ou=users,dc=example,dc=org",
		GroupSearchBase:     "ou=groups,dc=example,dc=org",
		UserSearchFilter:    "(&(objectClass=inetOrgPerson)(uid={}))",
		GroupSearchFilter:   "(&(objectClass=groupOfNames)(uid={}))",
		AttrEmail:           "mail",
		AttrUsername:        "uid",
		AttrUserDisplayName: "cn",
		AttrUserPhoneNo:     "telephoneNumber",
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

func (x *LdapOptions) Valid() error {
	if x == nil {
		return fmt.Errorf("ldap options is null")
	}
	if govalidator.IsNull(x.Host) {
		return fmt.Errorf("ldap host option  is null")
	}
	if govalidator.IsNull(x.ManagerDn) {
		return fmt.Errorf("ldap manager_dn option is null")
	}
	if govalidator.IsNull(x.ManagerPassword) {
		return fmt.Errorf("ldap manager_password option is null")
	}
	if govalidator.IsNull(x.UserSearchBase) {
		return fmt.Errorf("ldap user_search_base option is null")
	}
	if govalidator.IsNull(x.UserSearchFilter) {
		return fmt.Errorf("ldap user_search_filter option is null")
	}
	if govalidator.IsNull(x.GroupSearchBase) {
		return fmt.Errorf("ldap group_search_base option is null")
	}
	if govalidator.IsNull(x.GroupSearchFilter) {
		return fmt.Errorf("ldap group_search_base option is null")
	}
	if govalidator.IsNull(x.AttrEmail) {
		return fmt.Errorf("ldap attr_email option is null")
	}
	if govalidator.IsNull(x.AttrUsername) {
		return fmt.Errorf("ldap attr_username option is null")
	}
	if govalidator.IsNull(x.AttrUserDisplayName) {
		return fmt.Errorf("ldap attr_user_display_name option is null")
	}
	return nil
}
