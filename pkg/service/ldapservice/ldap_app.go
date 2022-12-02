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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	goldap "github.com/go-ldap/ldap"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/pkg/client/ldap"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
	"github.com/MicroOps-cn/idas/pkg/utils/sets"
	w "github.com/MicroOps-cn/idas/pkg/utils/wrapper"
)

const GroupStatusName = "status"

// getAppDnByEntryUUID
//
//	@Description[en-US]: Search and obtain the application dn (LDAP distinguished name) through UUID under "app_search_base".
//	@Description[zh-CN]: 从“app_search_base”下通过UUID搜索并获取应用dn(LDAP distinguished name)。
//	@param ctx  context.Context
//	@param id   string
//	@return dn  string
//	@return err error
func (s UserAndAppService) getAppDnByEntryUUID(ctx context.Context, id string) (dn string, err error) {
	u, err := uuid.FromString(id)
	if err != nil {
		return "", errors.ParameterError(fmt.Sprintf("id <%s> format error", id))
	}
	conn := s.Session(ctx)
	defer conn.Close()
	if result, err := conn.Search(goldap.NewSearchRequest(
		s.Options().AppSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(entryUUID=%s)", u.String()),
		[]string{},
		nil,
	)); err != nil {
		return "", err
	} else if len(result.Entries) == 0 {
		return "", errors.StatusNotFound(id)
	} else {
		dn = result.Entries[0].DN
	}
	return dn, nil
}

// getAppDetailByDn
//
//	@Description[en-US]: Use the dn(LDAP distinguished name) or application name to search for application information from the LDAP directory specified by "app_search_base". The directory level of the search is 1.
//	@Description[zh-CN]: 使用dn(LDAP distinguished name)或应用名称从 app_search_base 指定的LDAP目录内搜索应用信息, 搜索的目录层级为1
//	@param ctx         context.Context
//	@param dn          string
//	@return appDetail  *models.App
//	@return err        error
func (s UserAndAppService) getAppDetailByDn(ctx context.Context, dn string) (*models.App, error) {
	searchReq := goldap.NewSearchRequest(
		dn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)
	return s.getAppDetailByReq(ctx, searchReq)
}

// getAppDetailByReq
//
//	@Description[en-US]: Use the <ldap.SearchRequest> to search for application information from the LDAP directory specified by "app_search_base". The directory level of the search is 1.
//	@Description[zh-CN]: 使用<ldap.SearchRequest>从 app_search_base 指定的LDAP目录内搜索应用信息, 搜索的目录层级为1
//	@param ctx         context.Context
//	@param searchReq   *ldap.SearchRequest
//	@return appDetail  *models.App
//	@return err        error
func (s UserAndAppService) getAppDetailByReq(ctx context.Context, searchReq *goldap.SearchRequest) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq.Attributes = []string{"entryUUID", "description", "cn", "avatar", "uniqueMember", "member", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName}

	ret, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(ret.Entries) == 0 {
		return nil, errors.StatusNotFound("App")
	}
	appEntry := ret.Entries[0]
	var users []*models.User
	var userMap map[string]*models.User
	var roles []*models.AppRole

	member := append(appEntry.GetAttributeValues("uniqueMember"), appEntry.GetAttributeValues("member")...)
	if len(member) > 0 {
		userMap = make(map[string]*models.User, len(member))
		for _, userDn := range member {
			userInfo, err := s.getUserDetailByDn(ctx, userDn)
			if err != nil {
				if ldap.IsLdapError(err, goldap.LDAPResultControlNotFound) {
					return nil, err
				}
				continue
			}
			userMap[userDn] = userInfo
			users = append(users, userInfo)
		}
		searchMemberGroupReq := goldap.NewSearchRequest(
			appEntry.DN,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
			s.GetAppRoleSearchFilter(),
			[]string{"entryUUID", "member", "cn", "isDefault"},
			nil,
		)
		searchMemberGroupRet, err := conn.Search(searchMemberGroupReq)
		if err != nil {
			if !ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) {
				return nil, err
			}
		} else {
			roles = make([]*models.AppRole, len(searchMemberGroupRet.Entries))
			for idx, entry := range searchMemberGroupRet.Entries {
				roleName := entry.GetAttributeValue("cn")
				roles[idx] = &models.AppRole{Model: models.Model{Id: entry.GetAttributeValue("entryUUID")}, Name: roleName}
				if strings.ToLower(entry.GetAttributeValue("isDefault")) == "true" {
					roles[idx].IsDefault = true
				}
				for _, m := range entry.GetAttributeValues("member") {
					for userDn, user := range userMap {
						if userDn == m {
							user.Role = roleName
							user.RoleId = entry.GetAttributeValue("entryUUID")
						}
					}
				}
			}
		}
	}

	return &models.App{
		Model: models.Model{
			Id:         appEntry.GetAttributeValue("entryUUID"),
			CreateTime: w.M[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: w.M[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Name:        appEntry.GetAttributeValue("cn"),
		Description: appEntry.GetAttributeValue("description"),
		Avatar:      appEntry.GetAttributeValue("avatar"),
		Status:      models.AppMeta_Status(w.M[int](httputil.NewValue(appEntry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
		GrantMode:   models.AppMeta_GrantMode(w.M[int](httputil.NewValue(appEntry.GetAttributeValue("grantMode")).Default("0").Int())),
		GrantType:   models.AppMeta_GrantType(w.M[int](httputil.NewValue(appEntry.GetAttributeValue("grantType")).Default("0").Int())),
		Storage:     s.name,
		Roles:       roles,
		Users:       users,
	}, nil
}

// getAppRoleByUserDnAndAppDn
//
//	@Description[en-US]: Get the permission of the specified user under the application
//	@Description[zh-CN]: 获取应用下指定用户的权限
//	@param ctx       context.Context
//	@param appDn     string
//	@param userDn    string
//	@return roleId   string
//	@return roleName string
//	@return err      error
func (s UserAndAppService) getAppRoleByUserDnAndAppDn(ctx context.Context, appDn string, userDn string) (roleId, roleName string, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq := goldap.NewSearchRequest(
		appDn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", []string{"cn", "member", "uniqueMember", "objectClass"}, nil,
	)
	ret, err := conn.Search(searchReq)
	if err != nil {
		return "", "", err
	}
	if len(ret.Entries) == 0 {
		return "", "", errors.StatusNotFound("app")
	}
	appEntry := ret.Entries[0]
	if !sets.New[string](append(appEntry.GetAttributeValues("uniqueMember"), appEntry.GetAttributeValues("member")...)...).Has(userDn) {
		return "", "", fmt.Errorf("%s is not authorized to user %s", appDn, userDn)
	}
	searchMemberGroupReq := goldap.NewSearchRequest(
		appEntry.DN,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(`(&%s(member=%s))`, s.GetAppRoleSearchFilter(), userDn),
		[]string{"entryUUID", "member", "cn", "isDefault"}, nil,
	)
	searchMemberGroupRet, err := conn.Search(searchMemberGroupReq)
	if err != nil {
		if ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) {
			return "", "", err
		}
	} else if len(searchMemberGroupRet.Entries) > 0 {
		entry := searchMemberGroupRet.Entries[0]
		if strings.ToLower(entry.GetAttributeValue("isDefault")) == "true" {
			roleName = entry.GetAttributeValue("cn")
			roleId = entry.GetAttributeValue("entryUUID")
		}
		return entry.GetAttributeValue("entryUUID"), entry.GetAttributeValue("cn"), nil
	}
	return roleId, roleName, nil
}

// GetApps
//
//	@Description[en-US]: Get the application list. The application information does not include agent, role, user and other information.
//	@Description[zh-CN]: 获取应用列表，应用信息中不包含代理、角色、用户等信息。
//	@param ctx       context.Context
//	@param keywords  string
//	@param current   int64
//	@param pageSize  int64
//	@return total    int64
//	@return apps     []*models.App
//	@return err      error
func (s UserAndAppService) GetApps(ctx context.Context, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{s.GetAppSearchFilter()}
	if len(keywords) > 0 {
		filters = append(filters, fmt.Sprintf("(|(cn=*%s*)(description=*%s*))", keywords, keywords))
	}
	var filter string
	if len(filters) >= 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(filters, ""))
	}
	req := goldap.NewSearchRequest(
		s.Options().AppSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"entryUUID", "description", "cn", "avatar", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
		nil,
	)
	ret, err := conn.Search(req)
	if err != nil {
		return 0, nil, err
	}
	total = int64(len(ret.Entries))
	entrys := ret.Entries
	if int((current-1)*pageSize) > len(entrys) {
		return
	} else if int(current*pageSize) < len(entrys) {
		entrys = ret.Entries[(current-1)*pageSize : current*pageSize]
	} else {
		entrys = ret.Entries[(current-1)*pageSize:]
	}
	for _, entry := range entrys {
		apps = append(apps, &models.App{
			Model: models.Model{
				Id:         entry.GetAttributeValue("entryUUID"),
				CreateTime: w.M[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: w.M[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Name:        entry.GetAttributeValue("cn"),
			Description: entry.GetAttributeValue("description"),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.AppMeta_Status(w.M[int](httputil.NewValue(entry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
			GrantMode:   models.AppMeta_GrantMode(w.M[int](httputil.NewValue(entry.GetAttributeValue("grantMode")).Default("0").Int())),
			GrantType:   models.AppMeta_GrantType(w.M[int](httputil.NewValue(entry.GetAttributeValue("grantType")).Default("0").Int())),
			Storage:     s.name,
		})
	}
	return total, apps, nil
}

func (s UserAndAppService) DeepDeleteEntry(ctx context.Context, dn string) (err error) {
	conn := s.Session(ctx)
	searchChildrenReq := goldap.NewSearchRequest(
		dn,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)", nil, nil,
	)
	children, err := conn.Search(searchChildrenReq)
	if err != nil {
		return err
	}
	for _, entry := range children.Entries {
		if entry.DN != dn {
			err = s.DeepDeleteEntry(ctx, entry.DN)
			if err != nil {
				return err
			}
		}
	}

	return conn.Del(goldap.NewDelRequest(dn, nil))
}

// DeleteApps
//
//	@Description[en-US]: Delete apps in batch.
//	@Description[zh-CN]: 批量删除应用。
//	@param ctx     context.Context
//	@param ids     []string         : ID List
//	@return total  int64            : The quantity has been deleted. Since go ldap does not support transactions temporarily, an error may be reported after deleting a part.
//	@return err    error
func (s UserAndAppService) DeleteApps(ctx context.Context, ids []string) (total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, id := range ids {
		var dn string
		if dn, err = s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id); err != nil {
			return total, err
		} else if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
			fmt.Println("delete: ", dn)
			err = s.DeepDeleteEntry(ctx, dn)
			if err != nil {
				return
			}
		}
		total++
	}
	return
}

// DeleteApp
//
//	@Description[en-US]: Delete an app.
//	@Description[zh-CN]: 删除应用。
//	@param ctx 	context.Context
//	@param id 	string
//	@return err	error
func (s UserAndAppService) DeleteApp(ctx context.Context, id string) (err error) {
	return w.Error[int64](s.DeleteApps(ctx, []string{id}))
}

// UpdateApp
//
//	@Description[en-US]: Update applies the value of the specified column. If no column is specified, all column information is updated.
//	@Description[zh-CN]: 更新应用指定列的值，如果未指定列，则表示更新所有列信息。
//	@param ctx           context.Context
//	@param app           *models.App
//	@param updateColumns ...string
//	@return newApp       *models.App
//	@return err          error
func (s UserAndAppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (newApp *models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	dn, err := s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), app.Id)
	if err != nil {
		return nil, err
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"name", "description", "avatar", "grant_type", "grant_mode", "status", "user"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(dn, nil)

	member := make([]string, len(app.Users))
	for idx, user := range app.Users {
		if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
			return nil, err
		}
	}

	replace := []ldapUpdateColumn{
		{columnName: "name", ldapColumnName: "cn", val: []string{app.Name}},
		{columnName: "description", ldapColumnName: "description", val: []string{app.Description}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{app.Avatar}},
		{columnName: "grant_type", ldapColumnName: "grantType", val: []string{strconv.Itoa(int(app.GrantType))}},
		{columnName: "grant_mode", ldapColumnName: "grantMode", val: []string{strconv.Itoa(int(app.GrantMode))}},
		{columnName: "status", ldapColumnName: "status", val: []string{strconv.Itoa(int(app.Status))}},
		{columnName: "user", ldapColumnName: "uniqueMember", val: member},
	}
	for _, value := range replace {
		if columns.Has(value.columnName) && len(value.val) > 0 && len(value.val[0]) > 0 {
			req.Replace(value.ldapColumnName, value.val)
		}
	}

	if len(req.Changes) > 0 {
		if err = conn.Modify(req); err != nil {
			return nil, err
		}
	}

	if len(app.Roles) > 0 {
		for _, role := range app.Roles {
			for _, user := range app.Users {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && user.Role == role.Name {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			roleDn := fmt.Sprintf("cn=%s,%s", role.Name, dn)
			for _, user := range role.Users {
				fmt.Println(role.Name, role.Id, ">>", user.Id, user.Username)
			}
			if err = s.PatchAppRole(context.WithValue(ctx, global.LDAPConnName, conn), roleDn, role); err != nil {
				return nil, err
			}
		}
	}

	newApp, err = s.GetAppInfo(context.WithValue(ctx, global.LDAPConnName, conn), app.Id, "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newApp == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newApp, nil
}

// GetAppInfo
//
//	@Description[en-US]: Use the ID or application name to search for application information from the LDAP directory specified by "app_search_base". The directory level of the search is 1.
//	@Description[zh-CN]: 使用ID或应用名称从 app_search_base 指定的LDAP目录内搜索应用信息, 搜索的目录层级为1
//	@param ctx  context.Context
//	@param id   string          : App ID
//	@param name string          : App Name
//	@return app *models.App     : App Details
//	@return err error
func (s UserAndAppService) GetAppInfo(ctx context.Context, id, name string) (appDetail *models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) > 0 {
		dn, err := s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id)
		if err != nil {
			return nil, err
		}
		if len(dn) > 0 {
			return s.getAppDetailByDn(ctx, dn)
		}
	}

	if len(name) > 0 {
		req := goldap.NewSearchRequest(
			s.Options().AppSearchBase,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
			s.GetAppSearchFilter(name),
			nil,
			nil,
		)
		return s.getAppDetailByReq(ctx, req)
	}
	return nil, errors.NotFoundError
}

// CreateApp
//
//	@Description[en-US]: Create an app.
//	@Description[zh-CN]: 创建应用
//	@param ctx        context.Context
//	@param app        *models.App
//	@return appDetail *models.App
//	@return error
func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	dn := fmt.Sprintf("cn=%s,%s", app.Name, s.Options().AppSearchBase)
	req := goldap.NewAddRequest(dn, nil)

	var err error
	member := make([]string, len(app.Users))
	for idx, user := range app.Users {
		if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
			return nil, err
		}
	}

	attrs := map[string][]string{
		"description":  {app.Description},
		"avatar":       {app.Avatar},
		"grantType":    {strconv.Itoa(int(app.GrantType))},
		"grantMode":    {strconv.Itoa(int(app.GrantMode))},
		"status":       {strconv.Itoa(int(app.Status))},
		"objectClass":  append(s.GetAppClass().List(), "groupOfUniqueNames", "top"),
		"uniqueMember": member,
	}

	if len(app.Id) > 0 {
		attrs["entryUUID"] = []string{app.Id}
	}

	for name, value := range attrs {
		if len(value) > 0 && len(value[0]) > 0 {
			req.Attribute(name, value)
		}
	}
	if err = conn.Add(req); err != nil {
		return nil, err
	}

	if len(app.Roles) > 0 {
		for _, role := range app.Roles {
			for _, user := range app.Users {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && string(user.Role) == role.Name {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			roleDn := fmt.Sprintf("cn=%s,%s", role.Name, dn)
			if err = s.PatchAppRole(context.WithValue(ctx, global.LDAPConnName, conn), roleDn, role); err != nil {
				return nil, err
			}
		}
	}
	newAppInfo, err := s.getAppDetailByDn(context.WithValue(ctx, global.LDAPConnName, conn), dn)
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed: "+err.Error())
	} else if newAppInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed. ")
	}
	return newAppInfo, nil
}

var ldapColumnMap = map[string]string{
	"appname":     "name",
	"description": "description",
	"avatar":      "avatar",
	"grant_type":  "grantType",
	"grant_mode":  "grantMode",
	"user":        "uniqueMember",
	"status":      "status",
}

// PatchApp
//
//	@Description[en-US]: Incremental update application.
//	@Description[zh-CN]: 增量更新应用。
//	@param ctx        context.Context
//	@param fields     map[string]interface{}
//	@return appDetail app *models.App
//	@return err       error
func (s UserAndAppService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
	id, ok := fields["id"].(string)
	if !ok {
		return nil, errors.ParameterError("unknown id")
	}
	if !strings.HasSuffix(id, s.Options().AppSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	req := goldap.NewModifyRequest(id, nil)

	for name, value := range fields {
		ldapColumnName, ok := ldapColumnMap[name]
		if !ok {
			return nil, errors.ParameterError("unsupported field name: " + name)
		}
		switch val := value.(type) {
		case float64:
			req.Replace(ldapColumnName, []string{strconv.Itoa(int(val))})
		case string:
			req.Replace(ldapColumnName, []string{val})
		default:
			return nil, errors.ParameterError(fmt.Sprintf("unsupported field value type: name=%s,type=%T", name, value))
		}
	}

	conn := s.Session(ctx)
	defer conn.Close()
	if err := conn.Modify(req); err != nil {
		return nil, err
	}
	newAppInfo, err := s.GetAppInfo(context.WithValue(ctx, global.LDAPConnName, conn), fields["id"].(string), "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newAppInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newAppInfo, nil
}

// PatchApps
//
//	@Description[en-US]: Incrementally update information of multiple applications.
//	@Description[zh-CN]: 增量更新多个应用的信息。
//	@param ctx     context.Context
//	@param patch   []map[string]interface{}
//	@return total  int64
//	@return err    error
func (s UserAndAppService) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		id, ok := patchInfo["id"].(string)
		if !ok {
			return total, errors.ParameterError("unknown id")
		}
		dn, err := s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id)
		if err != nil {
			return total, err
		}
		req := goldap.NewModifyRequest(dn, nil)
		for name, value := range patchInfo {
			switch name {
			case GroupStatusName:
				status := value.(float64)
				req.Replace(GroupStatusName, []string{strconv.Itoa(int(status))})
			default:
				return total, errors.ParameterError("unsupported field name: " + name)
			}
		}
		if len(req.Changes) > 0 {
			if err = conn.Modify(req); err != nil {
				return total, err
			}
		}
		total++
	}
	return
}

// PatchAppRole
//
//	@Description[en-US]: Update App Role.
//	@Description[zh-CN]: 更新应用角色。
//	@param ctx     context.Context
//	@param dn      string
//	@param patch   *models.AppRole
//	@return err    error
func (s UserAndAppService) PatchAppRole(ctx context.Context, dn string, role *models.AppRole) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	var member []string
	if len(role.Users) == 0 {
		member = []string{""}
	} else {
		member = make([]string, len(role.Users))
		for idx, user := range role.Users {
			if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
				return err
			}
		}
	}

	searchReq := goldap.NewSearchRequest(dn,
		goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		s.GetAppRoleSearchFilter(role.Name), nil, nil,
	)
	if _, err = conn.Search(searchReq); err != nil {
		if !ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) { // 如果返回错误且错误不为搜索结果为空，则返回异常
			return errors.NewServerError(http.StatusInternalServerError, "Internal Server Error.")
		}
	} else {
		updateReq := goldap.NewModifyRequest(dn, nil)
		updateReq.Replace("objectClass", s.GetAppRoleGroupClass().Insert("groupOfNames", "top").List())
		updateReq.Replace("member", member)
		updateReq.Replace("isDefault", []string{strings.ToUpper(strconv.FormatBool(role.IsDefault))})
		if err = conn.Modify(updateReq); err == nil {
			return nil
		} else if ldap.IsLdapError(err, goldap.LDAPResultObjectClassModsProhibited) {
			if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	createReq := goldap.NewAddRequest(dn, nil)
	createReq.Attribute("objectClass", append(s.GetAppRoleGroupClass().List(), "groupOfNames", "top"))
	createReq.Attribute("member", member)
	createReq.Attribute("isDefault", []string{strings.ToUpper(strconv.FormatBool(role.IsDefault))})
	return conn.Add(createReq)
}

// VerifyUserAuthorizationForApp
//
//	@Description[en-US]: Verify user authorization for the application.
//	@Description[zh-CN]: 验证应用程序的用户授权
//	@param ctx    context.Context
//	@param appId  string
//	@param userId string
//	@return role  string   :Role name, such as admin, viewer, editor ...
//	@return err   error
func (s UserAndAppService) VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (role string, err error) {
	info, err := s.GetAppInfo(ctx, appId, "")
	if err != nil {
		return "", err
	}
	for _, user := range info.Users {
		if user.Id == userId {
			return user.Role, nil
		}
	}
	return "", errors.StatusNotFound("authorization")
}
