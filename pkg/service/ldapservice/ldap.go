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

package ldapservice

import (
	"context"
	"strings"

	"github.com/MicroOps-cn/fuck/sets"
	goldap "github.com/go-ldap/ldap/v3"

	"github.com/MicroOps-cn/idas/pkg/client/ldap"
)

const (
	ClassIdasCore         = "idasCore"
	ClassIdasApp          = "idasApp"
	ClassExtensibleObject = "extensibleObject"
)

func NewUserAndAppService(ctx context.Context, name string, client *ldap.Client) *UserAndAppService {
	classes, err := ldap.GetAvailableObjectClass(client.Session(ctx))
	if err != nil {
		return nil
	}
	var hasIDASClass bool
	if classes.HasAll("idasCore", "idasApp") {
		hasIDASClass = true
	}
	var appObjectClass []string
	var memberAttr string
	if hasIDASClass {
		appObjectClass = []string{ClassIdasCore, ClassIdasApp}
	} else {
		appObjectClass = []string{ClassExtensibleObject}
	}
	appMemberClass := client.Options().GetAppObjectClass()
	if len(appMemberClass) == 0 || appMemberClass == "groupOfUniqueNames" {
		appObjectClass = append(appObjectClass, "groupOfUniqueNames", "top")
		memberAttr = "uniqueMember"
	} else {
		appObjectClass = append(appObjectClass, "groupOfNames", "top")
		memberAttr = "member"
	}
	return &UserAndAppService{name: name, Client: client, hasIDASClass: hasIDASClass, appObjectClass: appObjectClass, memberAttr: memberAttr}
}

type UserAndAppService struct {
	*ldap.Client
	name           string
	hasIDASClass   bool
	appObjectClass []string
	memberAttr     string
}

func (s UserAndAppService) GetUserClass() sets.Set[string] {
	if s.hasIDASClass {
		return sets.New[string](ClassIdasCore)
	}
	return sets.New[string](ClassExtensibleObject)
}

func (s UserAndAppService) GetAppClass() []string {
	return s.appObjectClass
}

func (s UserAndAppService) GetMemberAttr() string {
	return s.memberAttr
}

func (s UserAndAppService) Name() string {
	return s.name
}

func (s UserAndAppService) GetAppSearchFilter(appName ...string) string {
	if len(appName) == 0 {
		appName = []string{"*"}
	}
	return strings.ReplaceAll(s.Options().GetAppSearchFilter(), "{}", appName[0])
}

func (s UserAndAppService) AutoCreateOrganizationalUnit(ctx context.Context, name string) error {
	session := s.Session(ctx)
	_, err := session.Search(goldap.NewSearchRequest(
		name,
		goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)",
		nil,
		nil))
	if ldap.IsLdapError(err, 32) {
		if _, suffix, found := strings.Cut(name, ","); found && len(suffix) > 0 {
			if err = s.AutoCreateOrganizationalUnit(ctx, suffix); err != nil {
				return err
			}
		}
		addReq := goldap.NewAddRequest(name, nil)
		addReq.Attributes = append(addReq.Attributes, goldap.Attribute{
			Type: "objectClass", Vals: []string{"top", "organizationalUnit"},
		})
		err = session.Add(addReq)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s UserAndAppService) AutoMigrate(ctx context.Context) error {
	err := s.AutoCreateOrganizationalUnit(ctx, s.Options().UserSearchBase)
	if err != nil {
		return err
	}
	return s.AutoCreateOrganizationalUnit(ctx, s.Options().AppSearchBase)
}

type ldapUpdateColumn struct {
	columnName     string
	ldapColumnName string
	val            []string
	oriVals        []string
}
