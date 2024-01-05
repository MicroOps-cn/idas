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

package ldap

import (
	"bytes"
	"fmt"
	"strings"

	fuck_tls "github.com/MicroOps-cn/fuck/clients/tls"
	"github.com/asaskevich/govalidator"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/pkg/errors"
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
	x.AppSearchBase = options.AppSearchBase
	x.AppSearchFilter = options.AppSearchFilter
	x.AttrEmail = options.AttrEmail
	x.AttrUsername = options.AttrUsername
	x.AttrUserDisplayName = options.AttrUserDisplayName
	x.AttrUserPhoneNo = options.AttrUserPhoneNo
	x.StartTLS = options.StartTLS
	x.IsTLS = options.IsTLS
	x.TLS = options.TLS
	err := unmarshaller.Unmarshal(bytes.NewReader(b), (*pbLdapOptions)(x))
	if err != nil {
		return err
	}
	if x.AppObjectClass != "groupOfUniqueNames" && x.AppObjectClass != "groupOfNames" {
		return fmt.Errorf("the ldap.app_object_class config can only be groupOfUniqueNames or groupOfNames")
	}
	return nil
}

// NewLdapOptions return a default option
// which host field point to nowhere.
func NewLdapOptions() *LdapOptions {
	return &LdapOptions{
		Host:                "127.0.0.1:389",
		ManagerDn:           "cn=admin,dc=example,dc=org",
		UserSearchBase:      "ou=users,dc=example,dc=org",
		AppSearchBase:       "ou=groups,dc=example,dc=org",
		UserSearchFilter:    "(&(objectClass=inetOrgPerson)(uid={}))",
		AppSearchFilter:     "(&(|(objectclass=groupOfNames)(objectclass=groupOfUniqueNames))(cn={}))",
		AppObjectClass:      "groupOfNames",
		AttrEmail:           "mail",
		AttrUsername:        "uid",
		AttrUserDisplayName: "cn",
		AttrUserPhoneNo:     "telephoneNumber",
		TLS: &fuck_tls.TLSOptions{
			MinVersion: "TLS12",
		},
	}
}

func (x *LdapOptions) ParseUserSearchFilter(username ...string) string {
	if len(username) == 0 {
		username = []string{"*"}
	}
	return strings.ReplaceAll(x.UserSearchFilter, "{}", username[0])
}

func (x *LdapOptions) Valid() error {
	if x == nil {
		return errors.New("ldap options is null")
	}
	if govalidator.IsNull(x.Host) {
		return errors.New("ldap host option is null")
	}
	if govalidator.IsNull(x.ManagerDn) {
		return errors.New("ldap manager_dn option is null")
	}
	passwd, err := x.ManagerPassword.UnsafeString()
	if err != nil {
		return err
	}
	if govalidator.IsNull(passwd) {
		return errors.New("ldap manager_password option is null")
	}
	if govalidator.IsNull(x.UserSearchBase) {
		return errors.New("ldap user_search_base option is null")
	}
	if govalidator.IsNull(x.UserSearchFilter) {
		return errors.New("ldap user_search_filter option is null")
	}
	if !strings.Contains(x.UserSearchFilter, "{}") {
		return errors.New("ldap user_search_filter option is invalid: does not contain {}")
	}
	if govalidator.IsNull(x.AppSearchBase) {
		return errors.New("ldap group_search_base option is null")
	}
	if govalidator.IsNull(x.AppSearchFilter) {
		return errors.New("ldap group_search_filter option is null")
	}
	if !strings.Contains(x.AppSearchFilter, "{}") {
		return errors.New("ldap group_search_filter option is invalid: does not contain {}")
	}
	if govalidator.IsNull(x.AttrEmail) {
		return errors.New("ldap attr_email option is null")
	}
	if govalidator.IsNull(x.AttrUsername) {
		return errors.New("ldap attr_username option is null")
	}
	if govalidator.IsNull(x.AttrUserDisplayName) {
		return errors.New("ldap attr_user_display_name option is null")
	}
	if x.StartTLS && x.IsTLS {
		return errors.New("ldap start_tls and is_tls cannot be both true")
	}
	if _, ok := fuck_tls.Versions[x.TLS.MinVersion]; ok {
		return nil
	}
	return fmt.Errorf("unknown TLS version: %s", x.TLS.MinVersion)
}
