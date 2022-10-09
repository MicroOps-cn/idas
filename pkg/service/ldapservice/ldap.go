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
	"github.com/MicroOps-cn/idas/pkg/client/ldap"
	"github.com/MicroOps-cn/idas/pkg/utils/sets"
	goldap "github.com/go-ldap/ldap"
	"strings"
)

const (
	ClassIdasCore         = "idasCore"
	ClassIdasApp          = "idasApp"
	ClassIdasRoleGroup    = "idasRoleGroup"
	ClassExtensibleObject = "extensibleObject"
)

func NewUserAndAppService(ctx context.Context, name string, client *ldap.Client) *UserAndAppService {
	classes, err := ldap.GetAvailableObjectClass(client.Session(ctx))
	if err != nil {
		return nil
	}
	var hasIDASClass bool
	if classes.HasAll("idasCore", "idasApp", "idasRoleGroup") {
		hasIDASClass = true
	}
	return &UserAndAppService{name: name, Client: client, hasIDASClass: hasIDASClass}
}

type UserAndAppService struct {
	*ldap.Client
	name         string
	hasIDASClass bool
}

func (s UserAndAppService) GetUserClass() sets.Set[string] {
	if s.hasIDASClass {
		return sets.New[string](ClassIdasCore)
	}
	return sets.New[string](ClassExtensibleObject)
}
func (s UserAndAppService) GetAppClass() sets.Set[string] {
	if s.hasIDASClass {
		return sets.New[string](ClassIdasCore, ClassIdasApp)
	}
	return sets.New[string](ClassExtensibleObject)
}

func (s UserAndAppService) GetAppRoleGroupClass() sets.Set[string] {
	if s.hasIDASClass {
		return sets.New[string](ClassIdasRoleGroup)
	}
	return sets.New[string](ClassExtensibleObject)
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

func (s UserAndAppService) GetAppRoleSearchFilter(roleName ...string) string {
	if len(roleName) == 0 {
		roleName = []string{"*"}
	}
	return strings.ReplaceAll(s.Options().GetAppRoleSearchFilter(), "{}", roleName[0])
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
}
