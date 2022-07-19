package ldapservice

import (
	"context"
	"fmt"
	goldap "github.com/go-ldap/ldap"
	uuid "github.com/satori/go.uuid"
	"idas/pkg/client/ldap"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service/models"
	"idas/pkg/utils/httputil"
	"idas/pkg/utils/sets"
	"idas/pkg/utils/wrapper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s UserAndAppService) GetApps(ctx context.Context, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{fmt.Sprintf(s.Options().ParseGroupSearchFilter())}
	if len(keywords) > 0 {
		filters = append(filters, fmt.Sprintf("(|(cn=*%s*)(description=*%s*))", keywords, keywords))
	}
	var filter string
	if len(filters) >= 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(filters, ""))
	}
	req := goldap.NewSearchRequest(
		s.Options().GroupSearchBase,
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
				CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Name:        entry.GetAttributeValue("cn"),
			Description: entry.GetAttributeValue("description"),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.AppMeta_Status(wrapper.Must[int](httputil.NewValue(entry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
			GrantMode:   models.AppMeta_GrantMode(wrapper.Must[int](httputil.NewValue(entry.GetAttributeValue("grantMode")).Default("0").Int())),
			GrantType:   models.AppMeta_GrantType(models.AppMeta_GrantType_value[entry.GetAttributeValue("grantType")]),
			Storage:     s.name,
		})
	}
	return
}

const GroupStatusName = "status"

func (s UserAndAppService) getAppDnByEntryUUID(ctx context.Context, id string) (dn string, err error) {
	u, err := uuid.FromString(id)
	if err != nil {
		return "", errors.ParameterError(fmt.Sprintf("id <%s> format error", id))
	}
	conn := s.Session(ctx)
	defer conn.Close()
	if result, err := conn.Search(goldap.NewSearchRequest(
		s.Options().GroupSearchBase,
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

func (s UserAndAppService) getAppInfoByDn(ctx context.Context, dn string) (*models.App, error) {
	searchReq := goldap.NewSearchRequest(
		dn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)
	return s.getAppInfoByReq(ctx, searchReq)
}

func (s UserAndAppService) getAppInfoByReq(ctx context.Context, searchReq *goldap.SearchRequest) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq.Attributes = []string{"entryUUID", "description", "cn", "avatar", "uniqueMember", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName}

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

	member := appEntry.GetAttributeValues("uniqueMember")
	if len(member) > 0 {
		userMap = make(map[string]*models.User, len(member))
		users = make([]*models.User, len(member))
		for idx, userDn := range member {
			userInfo, err := s.getUserInfoByDn(ctx, userDn)
			if err != nil {
				return nil, err
			}
			userMap[userDn] = userInfo
			users[idx] = userInfo
		}
		searchMemberGroupReq := goldap.NewSearchRequest(
			appEntry.DN,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
			"(objectClass=idasRoleGroup)",
			[]string{"entryUUID", "member", "cn", "isDefault"},
			nil,
		)
		searchMemberGroupRet, err := conn.Search(searchMemberGroupReq)
		if err != nil {
			if ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) {
				return nil, err
			}
		} else {
			roles = make([]*models.AppRole, len(searchMemberGroupRet.Entries))
			for idx, entry := range searchMemberGroupRet.Entries {
				var roleName = entry.GetAttributeValue("cn")
				roles[idx] = &models.AppRole{Model: models.Model{Id: entry.GetAttributeValue("entryUUID")}, Name: roleName}
				if strings.ToLower(entry.GetAttributeValue("isDefault")) == "true" {
					roles[idx].IsDefault = true
				}
				for _, m := range entry.GetAttributeValues("member") {
					for userDn, user := range userMap {
						if userDn == m {
							user.Role = models.UserRole(roleName)
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
			CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Name:        appEntry.GetAttributeValue("cn"),
		Description: appEntry.GetAttributeValue("description"),
		Avatar:      appEntry.GetAttributeValue("avatar"),
		Status:      models.AppMeta_Status(wrapper.Must[int](httputil.NewValue(appEntry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
		GrantMode:   models.AppMeta_GrantMode(wrapper.Must[int](httputil.NewValue(appEntry.GetAttributeValue("grantMode")).Default("0").Int())),
		GrantType:   models.AppMeta_GrantType(models.AppMeta_GrantType_value[appEntry.GetAttributeValue("grantType")]),
		Storage:     s.name,
		Role:        roles,
		User:        users,
	}, nil
}

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
			case "status":
				status := value.(float64)
				req.Replace(GroupStatusName, []string{strconv.Itoa(int(status))})
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

func (s UserAndAppService) DeleteApps(ctx context.Context, ids []string) (total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, id := range ids {
		var dn string
		if dn, err = s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id); err != nil {
			return total, err
		} else if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
			return
		}
		total++
	}
	return
}

type ldapUpdateColumn struct {
	columnName     string
	ldapColumnName string
	val            []string
}

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

	var member = make([]string, len(app.User))
	for idx, user := range app.User {
		if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
			return nil, err
		}
	}
	var replace = []ldapUpdateColumn{
		{columnName: "name", ldapColumnName: "cn", val: []string{app.Name}},
		{columnName: "description", ldapColumnName: "description", val: []string{app.Description}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{app.Avatar}},
		{columnName: "grant_type", ldapColumnName: "grantType", val: []string{string(app.GrantType)}},
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

	if len(app.Role) > 0 {
		for _, role := range app.Role {
			for _, user := range app.User {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.User = append(role.User, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && string(user.Role) == role.Name {
					role.User = append(role.User, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			roleDn := fmt.Sprintf("cn=%s,%s", role.Name, dn)

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

//getAppRoleByUserDnAndAppName Get the permission of the specified user under the application
// 获取应用下指定用户的权限
func (s UserAndAppService) getAppRoleByUserDnAndAppDn(ctx context.Context, appDn string, userDn string) (roleId, roleName string, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq := goldap.NewSearchRequest(
		appDn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", []string{"cn", "uniqueMember", "objectClass"}, nil,
	)
	ret, err := conn.Search(searchReq)
	if err != nil {
		return "", "", err
	}
	if len(ret.Entries) == 0 {
		return "", "", errors.StatusNotFound("app")
	}
	appEntry := ret.Entries[0]
	if !sets.New[string](appEntry.GetAttributeValues("uniqueMember")...).Has(userDn) {
		return "", "", fmt.Errorf("%s is not authorized to user %s", appDn, userDn)
	}
	searchMemberGroupReq := goldap.NewSearchRequest(
		appEntry.DN,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=idasRoleGroup)",
		[]string{"entryUUID", "member", "cn", "isDefault"},
		nil,
	)
	searchMemberGroupRet, err := conn.Search(searchMemberGroupReq)
	if err != nil {
		if ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) {
			return "", "", err
		}
	} else {
		for _, entry := range searchMemberGroupRet.Entries {
			if strings.ToLower(entry.GetAttributeValue("isDefault")) == "true" {
				roleName = entry.GetAttributeValue("cn")
				roleId = entry.GetAttributeValue("entryUUID")
			}
			for _, m := range entry.GetAttributeValues("member") {
				if m == userDn {
					return entry.GetAttributeValue("entryUUID"), entry.GetAttributeValue("cn"), nil
				}
			}
		}
	}
	return roleId, roleName, nil
}

func (s UserAndAppService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) == 0 && len(name) == 0 {
		return nil, errors.ParameterError("require id or appname")
	}

	var appEntry *goldap.Entry

	if len(id) != 0 {
		dn, err := s.getAppDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id)
		if err != nil {
			return nil, err
		}
		searchReq := goldap.NewSearchRequest(
			dn,
			goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
			"(objectClass=*)",
			[]string{"entryUUID", "description", "cn", "avatar", "uniqueMember", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
			nil,
		)
		ret, err := conn.Search(searchReq)
		if err != nil {
			return nil, err
		}
		if len(ret.Entries) > 0 {
			appEntry = ret.Entries[0]
		}
	}
	if appEntry == nil && len(name) != 0 {
		searchReq := goldap.NewSearchRequest(
			s.Options().GroupSearchBase,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
			s.Options().ParseGroupSearchFilter(name),
			[]string{"entryUUID", "description", "cn", "avatar", "uniqueMember", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
			nil,
		)
		ret, err := conn.Search(searchReq)
		if err != nil {
			return nil, err
		}
		if len(ret.Entries) > 0 {
			appEntry = ret.Entries[0]
		}
	}

	if appEntry == nil {
		return nil, errors.StatusNotFound("app")
	}
	var users []*models.User
	var userMap map[string]*models.User
	var roles []*models.AppRole

	member := appEntry.GetAttributeValues("uniqueMember")
	if len(member) > 0 {
		userMap = make(map[string]*models.User, len(member))
		users = make([]*models.User, len(member))
		for idx, userDn := range member {
			userInfo, err := s.getUserInfoByDn(ctx, userDn)
			if err != nil {
				return nil, err
			}
			userMap[userDn] = userInfo
			users[idx] = userInfo
		}
		searchReq := goldap.NewSearchRequest(
			appEntry.DN,
			goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 0, 0, false,
			"(objectClass=idasRoleGroup)",
			[]string{"entryUUID", "member", "cn", "isDefault"},
			nil,
		)
		ret, err := conn.Search(searchReq)
		if err != nil {
			if ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) {
				return nil, err
			}
		} else {
			roles = make([]*models.AppRole, len(ret.Entries))
			for idx, entry := range ret.Entries {
				var roleName = entry.GetAttributeValue("cn")
				roles[idx] = &models.AppRole{Model: models.Model{Id: entry.GetAttributeValue("entryUUID")}, Name: roleName}
				if strings.ToLower(entry.GetAttributeValue("isDefault")) == "true" {
					roles[idx].IsDefault = true
				}
				for _, m := range entry.GetAttributeValues("member") {
					for userDn, user := range userMap {
						if userDn == m {
							user.Role = models.UserRole(roleName)
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
			CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Name:        appEntry.GetAttributeValue("cn"),
		Description: appEntry.GetAttributeValue("description"),
		Avatar:      appEntry.GetAttributeValue("avatar"),
		Status:      models.AppMeta_Status(wrapper.Must[int](httputil.NewValue(appEntry.GetAttributeValue(GroupStatusName)).Default("0").Int())),
		GrantMode:   models.AppMeta_GrantMode(wrapper.Must[int](httputil.NewValue(appEntry.GetAttributeValue("grantMode")).Default("0").Int())),
		GrantType:   models.AppMeta_GrantType(models.AppMeta_GrantType_value[appEntry.GetAttributeValue("grantType")]),
		Storage:     s.name,
		User:        users,
		Role:        roles,
	}, nil
}

func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	dn := fmt.Sprintf("cn=%s,%s", app.Name, s.Options().GroupSearchBase)
	req := goldap.NewAddRequest(dn, nil)

	var err error
	member := make([]string, len(app.User))
	for idx, user := range app.User {
		if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
			return nil, err
		}
	}

	var attrs = map[string][]string{
		"description":  {app.Description},
		"avatar":       {app.Avatar},
		"grantType":    {string(app.GrantType)},
		"grantMode":    {strconv.Itoa(int(app.GrantMode))},
		"status":       {strconv.Itoa(int(app.Status))},
		"objectClass":  {"groupOfUniqueNames", "idasApp", "idasCore", "top"},
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

	if len(app.Role) > 0 {
		for _, role := range app.Role {
			for _, user := range app.User {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.User = append(role.User, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && string(user.Role) == role.Name {
					role.User = append(role.User, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			roleDn := fmt.Sprintf("cn=%s,%s", role.Name, dn)
			fmt.Println(role)
			fmt.Println(role.User)
			if err = s.PatchAppRole(context.WithValue(ctx, global.LDAPConnName, conn), roleDn, role); err != nil {
				return nil, err
			}
		}
	}
	newAppInfo, err := s.getAppInfoByDn(context.WithValue(ctx, global.LDAPConnName, conn), dn)
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed: "+err.Error())
	} else if newAppInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed. ")
	}
	return newAppInfo, nil
}

func (s UserAndAppService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
	id, ok := fields["id"].(string)
	if !ok {
		return nil, errors.ParameterError("unknown id")
	}
	if !strings.HasSuffix(id, s.Options().GroupSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	req := goldap.NewModifyRequest(id, nil)
	columns := sets.New[string]("appname", "email", "phone_number", "full_name", "avatar")
	for name, value := range fields {
		if !columns.Has(name) {
			continue
		}
		switch val := value.(type) {
		case float64:
			req.Replace(GroupStatusName, []string{strconv.Itoa(int(val))})
		case string:
			req.Replace(GroupStatusName, []string{val})
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

func (s UserAndAppService) DeleteApp(ctx context.Context, id string) (err error) {
	return wrapper.Error[int64](s.DeleteApps(ctx, []string{id}))
}

func (s UserAndAppService) PatchAppRole(ctx context.Context, dn string, role *models.AppRole) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	var member []string
	if len(role.User) == 0 {
		member = []string{""}
	} else {
		member = make([]string, len(role.User))
		for idx, user := range role.User {
			if member[idx], err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id); err != nil {
				return err
			}
		}
	}

	searchReq := goldap.NewSearchRequest(dn,
		goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&(objectClass=idasRoleGroup)(cn=%s))", role.Name),
		[]string{},
		nil,
	)
	if _, err = conn.Search(searchReq); err != nil {
		if !ldap.IsLdapError(err, goldap.LDAPResultNoSuchObject) { // 如果返回错误且错误不为搜索结果为空，则返回异常
			return errors.NewServerError(http.StatusInternalServerError, "Internal Server Error.")
		}
		createReq := goldap.NewAddRequest(dn, nil)
		createReq.Attribute("objectClass", []string{"idasRoleGroup", "groupOfNames", "top"})
		createReq.Attribute("member", member)
		createReq.Attribute("isDefault", []string{strings.ToUpper(strconv.FormatBool(role.IsDefault))})
		return conn.Add(createReq)
	} else {
		updateReq := goldap.NewModifyRequest(dn, nil)
		updateReq.Replace("objectClass", []string{"idasRoleGroup", "groupOfNames", "top"})
		updateReq.Replace("member", member)
		updateReq.Replace("isDefault", []string{strings.ToUpper(strconv.FormatBool(role.IsDefault))})
		return conn.Modify(updateReq)
	}
}

func (s UserAndAppService) VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (scope string, err error) {
	//TODO implement me
	panic("implement me")
}
