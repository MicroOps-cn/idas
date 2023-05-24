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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	goldap "github.com/go-ldap/ldap"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/pkg/client/ldap"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
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
func (s UserAndAppService) getAppDetailByDn(ctx context.Context, dn string, options *opts.GetAppOptions) (*models.App, error) {
	searchReq := goldap.NewSearchRequest(
		dn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)
	return s.getAppDetailByReq(ctx, searchReq, options)
}

// getAppDetailByReq
//
//	@Description[en-US]: Use the <ldap.SearchRequest> to search for application information from the LDAP directory specified by "app_search_base". The directory level of the search is 1.
//	@Description[zh-CN]: 使用<ldap.SearchRequest>从 app_search_base 指定的LDAP目录内搜索应用信息, 搜索的目录层级为1
//	@param ctx         context.Context
//	@param searchReq   *ldap.SearchRequest
//	@return appDetail  *models.App
//	@return err        error
func (s UserAndAppService) getAppDetailByReq(ctx context.Context, searchReq *goldap.SearchRequest, opt *opts.GetAppOptions) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq.Attributes = []string{"entryUUID", "description", "cn", "avatar", "uniqueMember", "member", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", "displayName", GroupStatusName, "url"}

	ret, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(ret.Entries) == 0 {
		return nil, errors.StatusNotFound("App")
	}
	appEntry := ret.Entries[0]
	var users []*models.User

	if opt == nil || !opt.DisableGetUsers {
		member := append(appEntry.GetAttributeValues("uniqueMember"), appEntry.GetAttributeValues("member")...)
		if len(member) > 0 {
			for _, userDn := range member {
				userInfo, err := s.getUserByDn(ctx, userDn)
				if err != nil {
					if ldap.IsLdapError(err, goldap.LDAPResultControlNotFound) {
						return nil, err
					}
					continue
				}
				if len(opt.UserId) == 0 || w.Include(opt.UserId, userInfo.Id) {
					users = append(users, userInfo)
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
		DisplayName: appEntry.GetAttributeValue("displayName"),
		Status:      models.AppMeta_Status(w.M[int](httputil.NewValue(appEntry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
		GrantMode:   models.AppMeta_GrantMode(w.M[int](httputil.NewValue(appEntry.GetAttributeValue("grantMode")).Default("0").Int())),
		GrantType:   models.AppMeta_GrantType(w.M[int](httputil.NewValue(appEntry.GetAttributeValue("grantType")).Default("0").Int())),
		Users:       users,
		Url:         appEntry.GetAttributeValue("url"),
	}, nil
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
func (s UserAndAppService) GetApps(ctx context.Context, keywords string, filters map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	fts := []string{s.GetAppSearchFilter()}
	if len(keywords) > 0 {
		fts = append(fts, fmt.Sprintf("(|(cn=*%s*)(description=*%s*))", keywords, keywords))
	}
	for name, val := range filters {
		switch name {
		case "id":
			name = "entryUUID"
		case "user_id":
			searchReq := goldap.NewSearchRequest(
				s.Options().UserSearchBase,
				goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 10, 0, false,
				fmt.Sprintf("(entryUUID=%s)", val), nil, nil,
			)
			dns, err := s.getDNSByReq(ldap.WithConnContext(ctx, conn), searchReq)
			if err != nil {
				return 0, nil, err
			}
			val = dns[0]
			fts = append(fts, fmt.Sprintf("(|(uniqueMember=%v)(member=%v))", val, val))
			continue
		}
		if ok, _ := regexp.MatchString("^[-_a-zA-Z0-9]+$", name); !ok {
			return 0, nil, errors.ParameterError(name)
		}
		value := fmt.Sprintf("%v", val)
		if ok, _ := regexp.MatchString("^[-_a-zA-Z0-9*]+$", value); !ok {
			return 0, nil, errors.ParameterError(name)
		}
		fts = append(fts, fmt.Sprintf("(%s=%v)", name, value))
	}
	var filter string
	if len(fts) >= 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(fts, ""))
	}
	req := goldap.NewSearchRequest(
		s.Options().AppSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"entryUUID", "description", "cn", "avatar", "createTimestamp", "modifyTimestamp", "displayName", "grantMode", "grantType", GroupStatusName, "url"},
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
			DisplayName: entry.GetAttributeValue("displayName"),
			Status:      models.AppMeta_Status(w.M[int](httputil.NewValue(entry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
			GrantMode:   models.AppMeta_GrantMode(w.M[int](httputil.NewValue(entry.GetAttributeValue("grantMode")).Default("0").Int())),
			GrantType:   models.AppMeta_GrantType(w.M[int](httputil.NewValue(entry.GetAttributeValue("grantType")).Default("0").Int())),
			Url:         entry.GetAttributeValue("url"),
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
		if dn, err = s.getAppDnByEntryUUID(ldap.WithConnContext(ctx, conn), id); err != nil {
			return total, err
		} else if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
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
	return w.E[int64](s.DeleteApps(ctx, []string{id}))
}

// UpdateApp
//
//	@Description[en-US]: Update applies the value of the specified column. If no column is specified, all column information is updated.
//	@Description[zh-CN]: 更新应用指定列的值，如果未指定列，则表示更新所有列信息。
//	@param ctx           context.Context
//	@param app           *models.App
//	@param updateColumns ...string
//	@return err          error
func (s UserAndAppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	ret, err := conn.Search(goldap.NewSearchRequest(
		s.Options().AppSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(entryUUID=%s)", app.Id), []string{"objectClass"}, nil))
	if err != nil {
		return err
	} else if len(ret.Entries) == 0 {
		return errors.StatusNotFound("App")
	}
	var memberAttr string
	entry := ret.Entries[0]
	objectClass := sets.New[string](entry.GetAttributeValues("objectClass")...)
	if objectClass.Has("groupOfNames") {
		memberAttr = "member"
	} else if objectClass.Has("groupOfUniqueNames") {
		memberAttr = "uniqueMember"
	} else {
		return errors.NewServerError(500, "This application does not support modification.")
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"name", "description", "avatar", "display_name", "grant_type", "grant_mode", "status", "user", "url"}
	}

	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(entry.DN, nil)
	if !objectClass.HasAll("idasApp", "idasCore") {
		objectClass.Insert("idasApp", "idasCore")
		req.Replace("objectClass", objectClass.List())
	}
	member := make([]string, len(app.Users))
	for idx, user := range app.Users {
		if member[idx], err = s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), user.Id); err != nil {
			return err
		}
	}

	replace := []ldapUpdateColumn{
		{columnName: "name", ldapColumnName: "cn", val: []string{app.Name}},
		{columnName: "description", ldapColumnName: "description", val: []string{app.Description}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{app.Avatar}},
		{columnName: "grant_type", ldapColumnName: "grantType", val: []string{strconv.Itoa(int(app.GrantType))}},
		{columnName: "display_name", ldapColumnName: "displayName", val: []string{app.DisplayName}},
		{columnName: "grant_mode", ldapColumnName: "grantMode", val: []string{strconv.Itoa(int(app.GrantMode))}},
		{columnName: "status", ldapColumnName: "status", val: []string{strconv.Itoa(int(app.Status))}},
		{columnName: "user", ldapColumnName: memberAttr, val: member},
		{columnName: "url", ldapColumnName: "url", val: []string{app.Url}},
	}
	for _, value := range replace {
		if columns.Has(value.columnName) && len(value.val) > 0 && len(value.val[0]) > 0 {
			req.Replace(value.ldapColumnName, value.val)
		}
	}

	if len(req.Changes) > 0 {
		if err = conn.Modify(req); err != nil {
			return err
		}
	}

	return nil
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
func (s UserAndAppService) GetAppInfo(ctx context.Context, o ...opts.WithGetAppOptions) (appDetail *models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	opt := opts.NewAppOptions(o...)
	if len(opt.Id) == 0 && len(opt.Name) == 0 {
		return nil, errors.ParameterError("require id or name")
	}
	if len(opt.Id) > 0 {
		u, err := uuid.FromString(opt.Id)
		if err != nil {
			return nil, errors.ParameterError(fmt.Sprintf("id <%s> format error", opt.Id))
		}
		searchReq := goldap.NewSearchRequest(
			s.Options().AppSearchBase,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
			fmt.Sprintf("(entryUUID=%s)", u.String()),
			[]string{},
			nil,
		)
		return s.getAppDetailByReq(ldap.WithConnContext(ctx, conn), searchReq, opt)
	}

	if len(opt.Name) > 0 {
		req := goldap.NewSearchRequest(
			s.Options().AppSearchBase,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
			s.GetAppSearchFilter(opt.Name),
			nil,
			nil,
		)
		return s.getAppDetailByReq(ctx, req, opt)
	}
	return nil, errors.NotFoundError()
}

// CreateApp
//
//	@Description[en-US]: Create an app.
//	@Description[zh-CN]: 创建应用
//	@param ctx        context.Context
//	@param app        *models.App
//	@return error
func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	dn := fmt.Sprintf("cn=%s,%s", app.Name, s.Options().AppSearchBase)
	req := goldap.NewAddRequest(dn, nil)

	member := make([]string, len(app.Users))
	for idx, user := range app.Users {
		if member[idx], err = s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), user.Id); err != nil {
			return err
		}
	}

	attrs := map[string][]string{
		"description":     {app.Description},
		"avatar":          {app.Avatar},
		"grantType":       {strconv.Itoa(int(app.GrantType))},
		"displayName":     {app.DisplayName},
		"grantMode":       {strconv.Itoa(int(app.GrantMode))},
		"status":          {strconv.Itoa(int(app.Status))},
		"objectClass":     s.GetAppClass(),
		s.GetMemberAttr(): member,
		"url":             {app.Url},
	}

	if len(app.Id) > 0 {
		attrs["entryUUID"] = []string{app.Id}
	}

	for name, value := range attrs {
		if len(value) > 0 && len(value[0]) > 0 {
			req.Attribute(name, value)
		}
	}

	err = conn.Add(req)
	if err != nil {
		return err
	}
	newApp, err := s.getAppDetailByDn(ldap.WithConnContext(ctx, conn), dn, opts.NewAppOptions(opts.WithoutUsers))
	if err != nil {
		return err
	}
	app.Model = newApp.Model
	return nil
}

// PatchApp
//
//	@Description[en-US]: Incremental update application.
//	@Description[zh-CN]: 增量更新应用。
//	@param ctx        context.Context
//	@param fields     map[string]interface{}
//	@return err       error
func (s UserAndAppService) PatchApp(ctx context.Context, fields map[string]interface{}) (err error) {
	id, ok := fields["id"].(string)
	if !ok {
		return errors.ParameterError("unknown id")
	}
	delete(fields, "id")

	conn := s.Session(ctx)
	defer conn.Close()

	ret, err := conn.Search(goldap.NewSearchRequest(
		s.Options().AppSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(entryUUID=%s)", id), []string{"objectClass"}, nil))
	if err != nil {
		return err
	} else if len(ret.Entries) == 0 {
		return errors.StatusNotFound("App")
	}
	var memberAttr string
	entry := ret.Entries[0]
	objectClass := sets.New[string](entry.GetAttributeValues("objectClass")...)
	if objectClass.Has("groupOfNames") {
		memberAttr = "member"
	} else if objectClass.Has("groupOfUniqueNames") {
		memberAttr = "uniqueMember"
	} else {
		return errors.NewServerError(500, "This application does not support modification.")
	}

	req := goldap.NewModifyRequest(entry.DN, []goldap.Control{goldap.NewControlManageDsaIT(false)})

	if !objectClass.HasAll("idasApp", "idasCore") {
		objectClass.Insert("idasApp", "idasCore")
		req.Replace("objectClass", objectClass.List())
	}

	ldapColumnMap := map[string]string{
		"appname":      "name",
		"description":  "description",
		"avatar":       "avatar",
		"grant_type":   "grantType",
		"display_name": "displayName",
		"grant_mode":   "grantMode",
		"user":         memberAttr,
		"status":       "status",
		"url":          "url",
	}
	for name, value := range fields {
		if value == nil {
			continue
		}
		ldapColumnName, ok := ldapColumnMap[name]
		if !ok {
			return errors.ParameterError("unsupported field name: " + name)
		}
		switch val := value.(type) {
		case float64:
			req.Replace(ldapColumnName, []string{strconv.Itoa(int(val))})
		case string:
			req.Replace(ldapColumnName, []string{val})
		case *models.AppMeta_Status:
			req.Replace(ldapColumnName, []string{strconv.Itoa(int(*val))})
		default:
			return errors.ParameterError(fmt.Sprintf("unsupported field value type: name=%s,type=%T", name, value))
		}
	}
	if len(req.Changes) == 0 {
		return nil
	}
	req.Replace("objectClass", s.GetAppClass())
	return conn.Modify(req)
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

		dn, err := s.getAppDnByEntryUUID(ldap.WithConnContext(ctx, conn), id)
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
