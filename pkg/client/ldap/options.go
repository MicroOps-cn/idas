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

	"github.com/asaskevich/govalidator"
	"github.com/gogo/protobuf/jsonpb"
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
	x.AppRoleSearchFilter = options.AppRoleSearchFilter
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
		AppSearchBase:       "ou=groups,dc=example,dc=org",
		UserSearchFilter:    "(&(objectClass=inetOrgPerson)(uid={}))",
		AppSearchFilter:     "(&(|(objectClass=idasApp)(objectClass=extensibleObject))(objectClass=groupOfUniqueNames)(cn={}))",
		AppRoleSearchFilter: "(&(|(objectClass=idasRoleGroup)(objectClass=extensibleObject))(objectClass=groupOfNames)(cn={}))",
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
	if !strings.Contains(x.UserSearchFilter, "{}") {
		return fmt.Errorf("ldap user_search_filter option is invalid: does not contain {}")
	}
	if govalidator.IsNull(x.AppSearchBase) {
		return fmt.Errorf("ldap group_search_base option is null")
	}
	if govalidator.IsNull(x.AppSearchFilter) {
		return fmt.Errorf("ldap group_search_filter option is null")
	}
	if !strings.Contains(x.AppSearchFilter, "{}") {
		return fmt.Errorf("ldap group_search_filter option is invalid: does not contain {}")
	}
	if govalidator.IsNull(x.AppRoleSearchFilter) {
		return fmt.Errorf("ldap group_search_filter option is null")
	}
	if !strings.Contains(x.AppRoleSearchFilter, "{}") {
		return fmt.Errorf("ldap group_search_filter option is invalid: does not contain {}")
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
