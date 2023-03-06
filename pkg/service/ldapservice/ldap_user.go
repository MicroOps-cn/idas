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
	"strconv"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log/level"
	goldap "github.com/go-ldap/ldap"
	uuid "github.com/satori/go.uuid"
	"github.com/tredoe/osutil/user/crypt/sha256_crypt"

	"github.com/MicroOps-cn/idas/pkg/client/ldap"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
)

const UserStatusName = "status"

// hash
//
//	@Description[en-US]: Use sha256 to generate the hash value of the password
//	@Description[zh-CN]: 使用sha256生成密码的hash值
//	@param password string
//	@return sha256 string
//	@return error
func hash(password []byte) (string, error) {
	c := sha256_crypt.New()
	salt := string(uuid.NewV4().Bytes())
	return c.Generate(password, []byte("$5$"+salt))
}

// getAppDnByEntryUUID
//
//	@Description[en-US]: Search and obtain the user dn (LDAP distinguished name) through UUID under "user_search_base".
//	@Description[zh-CN]: 从“user_search_base”下通过UUID搜索并获取用户dn(LDAP distinguished name)。
//	@param ctx  context.Context
//	@param id   string
//	@return dn  string
//	@return err error
func (s UserAndAppService) getUserDnByEntryUUID(ctx context.Context, id string) (dn string, err error) {
	u, err := uuid.FromString(id)
	if err != nil {
		return "", errors.ParameterError(fmt.Sprintf("id <%s> format error", id))
	}
	conn := s.Session(ctx)
	if result, err := conn.Search(goldap.NewSearchRequest(
		s.Options().UserSearchBase,
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

// getUserObjectClass
//
//	@Description[en-US]: Get the current objectClass of the user
//	@Description[zh-CN]: 获取用户当前的objectClass
//	@param ctx           context.Context
//	@param dn            string
//	@return objectClass  sets.Set[string]
//	@return err          error
func (s UserAndAppService) getUserObjectClass(ctx context.Context, dn string) (objectClass sets.Set[string], err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	searchReq := goldap.NewSearchRequest(
		dn,
		goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)",
		[]string{"objectClass"},
		nil,
	)
	ret, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(ret.Entries) == 0 {
		return nil, errors.StatusNotFound(fmt.Sprintf("user %s", dn))
	}
	objectClass = sets.New[string](ret.Entries[0].GetAttributeValues("objectClass")...)
	return objectClass, nil
}

// getDNSByReq
//
//	@Description[en-US]: Get DNs (LDAP distinguished name) by <*ldap.SearchRequest>.
//	@Description[zh-CN]: 通过<*ldap.SearchRequest>获取DNs(LDAP distinguished name)
//	@param ctx        context.Context
//	@param searchReq  *ldap.SearchRequest
//	@return dns       []string
//	@return err       error
func (s UserAndAppService) getDNSByReq(ctx context.Context, searchReq *goldap.SearchRequest) (dns []string, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	search, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	for _, entry := range search.Entries {
		dns = append(dns, entry.DN)
	}
	return dns, err
}

// getUserByReq
//
//	@Description[en-US]: Use the <ldap.SearchRequest> to search for application information from the LDAP directory specified by "user_search_base". The directory level of the search is 1.
//	@Description[zh-CN]: 使用<ldap.SearchRequest>从 user_search_base 指定的LDAP目录内搜索应用信息, 搜索的目录层级为1
//	@param ctx         context.Context
//	@param searchReq   *ldap.SearchRequest
//	@return userDetail *models.User
//	@return err        error
func (s UserAndAppService) getUserByReq(ctx context.Context, searchReq *goldap.SearchRequest) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq.Attributes = []string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "entryUUID", "avatar", "createTimestamp", "modifyTimestamp", UserStatusName, "memberOf"}
	ret, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(ret.Entries) == 0 {
		return nil, errors.StatusNotFound("user")
	}

	userEntry := ret.Entries[0]
	userInfo := &models.User{
		Model: models.Model{
			Id:         userEntry.GetAttributeValue("entryUUID"),
			CreateTime: w.M[time.Time](time.Parse("20060102150405Z", userEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: w.M[time.Time](time.Parse("20060102150405Z", userEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Username:    userEntry.GetAttributeValue(s.Options().GetAttrUsername()),
		Email:       userEntry.GetAttributeValue(s.Options().GetAttrEmail()),
		PhoneNumber: userEntry.GetAttributeValue(s.Options().GetAttrUserPhoneNo()),
		FullName:    userEntry.GetAttributeValue(s.Options().GetAttrUserDisplayName()),
		Avatar:      userEntry.GetAttributeValue("avatar"),
		Status:      models.UserMeta_UserStatus(w.M[int](httputil.NewValue(userEntry.GetAttributeValue(UserStatusName)).Default("0").Int())),
		Storage:     s.name,
	}
	return userInfo, nil
}

// getUserByDn
//
//	@Description[en-US]: Use DN to obtain user information.
//	@Description[zh-CN]: 使用DN获取用户信息。
//	@param ctx           context.Context
//	@param dn            string
//	@return userDetail   *models.User
//	@return err          error
func (s UserAndAppService) getUserByDn(ctx context.Context, dn string) (*models.User, error) {
	searchReq := goldap.NewSearchRequest(
		dn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)
	return s.getUserByReq(ctx, searchReq)
}

// getUserByUsername
//
//	@Description[en-US]: Use username to obtain user information.
//	@Description[zh-CN]: 使用用户名获取用户信息。
//	@param ctx           context.Context
//	@param username      string
//	@return userDetail   *models.User
//	@return err          error
func (s UserAndAppService) getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 2, 0, false,
		s.Options().ParseUserSearchFilter(username), nil, nil,
	)
	return s.getUserByReq(ctx, searchReq)
}

// GetUserInfoByUsernameAndEmail
//
//	@Description[en-US]: Use username or email to obtain user information.
//	@Description[zh-CN]: 使用用户名或email获取用户信息。
//	@param ctx           context.Context
//	@param username      string
//	@param email         string
//	@return userDetail   *models.User
//	@return err          error
func (s UserAndAppService) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (*models.User, error) {
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&(%s=%s)(%s=%s))",
			s.Options().GetAttrUsername(), username,
			s.Options().GetAttrEmail(), email),
		nil, nil,
	)
	return s.getUserByReq(ctx, searchReq)
}

// ResetPassword
//
//	@Description[en-US]: Reset User Password.
//	@Description[zh-CN]: 重置用户密码。
//	@param ctx       context.Context
//	@param id        string
//	@param password  string           : New password.
//	@return err      error
func (s UserAndAppService) ResetPassword(ctx context.Context, id string, password string) error {
	conn := s.Session(ctx)
	defer conn.Close()
	phash, err := hash([]byte(password))
	if err != nil {
		logger := log.GetContextLogger(ctx)
		level.Error(logger).Log("failed to general password hash")
		return err
	}
	dn, err := s.getUserDnByEntryUUID(ctx, id)
	if err != nil {
		return err
	}
	info, err := s.getUserByDn(ctx, dn)
	if err != nil {
		return err
	}
	if info.Status.Is(models.UserMeta_disabled) {
		return errors.NewServerError(500, "unknown user status: "+info.Status.String())
	}
	req := goldap.NewModifyRequest(dn, nil)
	req.Replace("userPassword", []string{"{CRYPT}" + phash})
	req.Replace(UserStatusName, []string{strconv.Itoa(int(models.UserMeta_normal))})
	return conn.Modify(req)
}

// UpdateLoginTime [Not Supported]
//
//	@Description[en-US]: Update the user's last login time.
//	@Description[zh-CN]: 更新用户最后一次登陆时间。
//	@param _
//	@param _
//	@return error
func (s UserAndAppService) UpdateLoginTime(_ context.Context, _ string) error {
	return nil
}

// GetUsers
//
//	@Description[en-US]: Get user list.
//	@Description[zh-CN]: 获取用户列表。
//	@param ctx       context.Context
//	@param keywords  string
//	@param status    models.UserMeta_UserStatus
//	@param appId     string
//	@param current   int64
//	@param pageSize  int64
//	@return total    int64
//	@return users    []*models.User
//	@return err      error
func (s UserAndAppService) GetUsers(ctx context.Context, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{s.Options().ParseUserSearchFilter()}
	if status != models.UserMetaStatusAll {
		filters = append(filters, fmt.Sprintf("(%s=%d)", UserStatusName, status))
	}
	if len(keywords) > 0 {
		filters = append(filters, fmt.Sprintf("(|(%s=*%s*)(%s=*%s*)(%s=*%s*)(%s=*%s*))",
			s.Options().GetAttrUsername(), keywords,
			s.Options().GetAttrUserDisplayName(), keywords,
			s.Options().GetAttrEmail(), keywords,
			s.Options().GetAttrUserPhoneNo(), keywords))
	}
	var filter string
	if len(filters) >= 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(filters, ""))
	}
	req := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "entryUUID", "avatar", "createTimestamp", "modifyTimestamp", UserStatusName},
		nil,
	)
	ret, err := conn.Search(req)
	if ldap.IsLdapError(err, 32) {
		return 0, nil, nil
	} else if err != nil {
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
		users = append(users, &models.User{
			Model: models.Model{
				Id:         entry.GetAttributeValue("entryUUID"),
				CreateTime: w.M[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: w.M[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Username:    entry.GetAttributeValue(s.Options().GetAttrUsername()),
			Email:       entry.GetAttributeValue(s.Options().GetAttrEmail()),
			PhoneNumber: entry.GetAttributeValue(s.Options().GetAttrUserPhoneNo()),
			FullName:    entry.GetAttributeValue(s.Options().GetAttrUserDisplayName()),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.UserMeta_UserStatus(w.M[int](httputil.NewValue(entry.GetAttributeValue(UserStatusName)).Default("0").Int())),
			Storage:     s.name,
		})
	}
	return
}

// PatchUsers
//
//	@Description[en-US]: Incrementally update information of multiple users.
//	@Description[zh-CN]: 增量更新多个用户的信息。
//	@param ctx 		context.Context
//	@param patch 	[]map[string]interface{}
//	@return count	int64
//	@return err		error
func (s UserAndAppService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		id, ok := patchInfo["id"].(string)
		if !ok {
			return count, errors.ParameterError("unknown id")
		}
		delete(patchInfo, "id")
		var dn string
		dn, err = s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), id)
		if err != nil {
			return count, err
		}
		objectClass, err := s.getUserObjectClass(ldap.WithConnContext(ctx, conn), dn)
		if err != nil {
			return count, err
		}
		objectClass.Insert(s.GetUserClass().List()...)
		req := goldap.NewModifyRequest(dn, nil)
		req.Replace("objectClass", objectClass.List())
		for name, value := range patchInfo {
			switch name {
			case "status":
				status := value.(int32)
				req.Replace(UserStatusName, []string{strconv.Itoa(int(status))})
			default:
				return count, errors.ParameterError("unsupported field name: " + name)
			}
		}
		if len(req.Changes) > 0 {
			if err = conn.Modify(req); err != nil {
				return count, err
			}
		}
		count++
	}
	return
}

// DeleteUsers
//
//	@Description[en-US]: Delete users in batch.
//	@Description[zh-CN]: 批量删除用户。
//	@param ctx 		context.Context
//	@param ids 		[]string
//	@return count	int64
//	@return err		error
func (s UserAndAppService) DeleteUsers(ctx context.Context, ids []string) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, id := range ids {
		var dn string
		if dn, err = s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), id); err != nil {
			return count, err
		} else if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
			return
		}
		count++
	}
	return
}

// UpdateUser
//
//	@Description[en-US]: Update user information.
//	@Description[zh-CN]: 更新用户信息.
//	@param ctx	context.Context
//	@param user	*models.User
//	@param updateColumns	...string
//	@return err	error
func (s UserAndAppService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	dn, err := s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), user.Id)
	if err != nil {
		return err
	}
	objectClass, err := s.getUserObjectClass(ldap.WithConnContext(ctx, conn), dn)
	if err != nil {
		return err
	}
	objectClass.Insert(s.GetUserClass().List()...)

	if len(updateColumns) == 0 {
		updateColumns = []string{"object_class", "username", "email", "phone_number", "full_name", "avatar", "status"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(dn, nil)
	replace := []ldapUpdateColumn{
		{columnName: "object_class", ldapColumnName: "objectClass", val: objectClass.List()},
		{columnName: "email", ldapColumnName: s.Options().GetAttrEmail(), val: []string{user.Email}},
		{columnName: "phone_number", ldapColumnName: s.Options().GetAttrUserPhoneNo(), val: []string{user.PhoneNumber}},
		{columnName: "full_name", ldapColumnName: s.Options().GetAttrUserDisplayName(), val: []string{user.FullName}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{user.Avatar}},
		{columnName: "status", ldapColumnName: "status", val: []string{strconv.Itoa(int(user.Status))}},
	}

	for _, value := range replace {
		if columns.Has(value.columnName) && len(value.val[0]) > 0 {
			req.Replace(value.ldapColumnName, value.val)
		}
	}

	if len(req.Changes) > 0 {
		return conn.Modify(req)
	}
	return nil
}

// GetUserInfoById
//
//	@Description[en-US]: Obtain user information through ID.
//	@Description[zh-CN]: 通过ID获取用户信息。
//	@param ctx 	context.Context
//	@param id 	string
//	@return userDetail	*models.User
//	@return err	error
func (s UserAndAppService) GetUserInfoById(ctx context.Context, id string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(entryUUID=%s)", id), nil, nil,
	)
	return s.getUserByReq(ldap.WithConnContext(ctx, conn), searchReq)
}

// GetUserInfo
//
//	@Description[en-US]: Obtain user information through ID or username.
//	@Description[zh-CN]: 通过ID或用户名获取用户信息。
//	@param ctx 	context.Context
//	@param id 	string
//	@param username 	string
//	@return userDetail	*models.User
//	@return err	error
func (s UserAndAppService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) == 0 && len(username) == 0 {
		return nil, errors.ParameterError("require id or username")
	}

	var userEntry *goldap.Entry

	if len(id) != 0 {
		return s.GetUserInfoById(ldap.WithConnContext(ctx, conn), id)
	}
	if userEntry == nil && len(username) != 0 {
		return s.getUserByUsername(ldap.WithConnContext(ctx, conn), username)
	}
	return nil, errors.StatusNotFound("user")
}

// CreateUser
//
//	@Description[en-US]: Create a user.
//	@Description[zh-CN]: 创建用户。
//	@param ctx 	context.Context
//	@param user 	*models.User
//	@return err	error
func (s UserAndAppService) CreateUser(ctx context.Context, user *models.User) (err error) {
	logger := log.GetContextLogger(ctx)
	conn := s.Session(ctx)
	defer conn.Close()
	dn := fmt.Sprintf("%s=%s,%s", s.Options().GetAttrUsername(), user.Username, s.Options().UserSearchBase)
	req := goldap.NewAddRequest(dn, nil)

	attrs := map[string][]string{
		"mail":                               {user.Email},
		s.Options().GetAttrUserPhoneNo():     {user.PhoneNumber},
		s.Options().GetAttrUsername():        {user.Username},
		s.Options().GetAttrUserDisplayName(): {user.FullName},
		"sn":                                 {" "},
		"avatar":                             {user.Avatar},
		"status":                             {strconv.Itoa(int(user.Status))},
		"objectClass":                        append(s.GetUserClass().List(), "inetOrgPerson", "organizationalPerson", "person", "top"),
	}
	if len(attrs["cn"]) == 0 || len(attrs["cn"][0]) == 0 {
		attrs["cn"] = []string{user.Username}
	}

	if len(user.Password) > 0 {
		passwordHash, err := hash(user.Password)
		if err != nil {
			level.Error(logger).Log("failed to general password hash")
		} else {
			attrs["userPassword"] = []string{"{CRYPT}" + passwordHash} // RFC4519/2307: password of user
		}
	}

	for name, value := range attrs {
		if value[0] != "" {
			req.Attribute(name, value)
		}
	}

	err = conn.Add(req)
	if err != nil {
		return err
	}
	newUser, err := s.getUserByDn(ldap.WithConnContext(ctx, conn), dn)
	if err != nil {
		return err
	}
	user.Model = newUser.Model
	return nil
}

// PatchUser
//
//	@Description[en-US]: Incremental update user.
//	@Description[zh-CN]: 增量更新用户。
//	@param ctx 	context.Context
//	@param user 	map[string]interface{}
//	@return err	error
func (s UserAndAppService) PatchUser(ctx context.Context, user map[string]interface{}) (err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	id, ok := user["id"].(string)
	if !ok {
		return errors.ParameterError("unknown id")
	}
	dn, err := s.getUserDnByEntryUUID(ldap.WithConnContext(ctx, conn), id)
	if err != nil {
		return err
	}
	req := goldap.NewModifyRequest(dn, nil)
	columns := []ldapUpdateColumn{
		{columnName: "email", ldapColumnName: s.Options().GetAttrEmail()},
		{columnName: "phone_number", ldapColumnName: s.Options().GetAttrUserPhoneNo()},
		{columnName: "full_name", ldapColumnName: s.Options().GetAttrUserDisplayName()},
		{columnName: "avatar", ldapColumnName: "avatar"},
		{columnName: "full_name", ldapColumnName: s.Options().GetAttrUserDisplayName()},
		{columnName: "status", ldapColumnName: "status"},
	}
	for _, column := range columns {
		if value, ok := user[column.columnName]; ok {
			switch val := value.(type) {
			case float64:
				req.Replace(column.ldapColumnName, []string{strconv.Itoa(int(val))})
			case models.UserMeta_UserStatus:
				req.Replace(column.ldapColumnName, []string{strconv.Itoa(int(val))})
			case string:
				req.Replace(column.ldapColumnName, []string{val})
			default:
				req.Replace(column.ldapColumnName, []string{fmt.Sprintf("%v", val)})
			}
		}
	}

	return conn.Modify(req)
}

// DeleteUser
//
//	@Description[en-US]: Delete a user.
//	@Description[zh-CN]: 删除用户。
//	@param ctx 	context.Context
//	@param id 	string
//	@return error
func (s UserAndAppService) DeleteUser(ctx context.Context, id string) error {
	return w.E[int64](s.DeleteUsers(ctx, []string{id}))
}

// VerifyPasswordById
//
//	@Description[en-US]: Verify the user's password through ID.
//	@Description[zh-CN]: 通过ID验证用户密码。
//	@param ctx 	context.Context
//	@param id 	string
//	@param password 	string
//	@return users	[]*models.User
func (s UserAndAppService) VerifyPasswordById(ctx context.Context, id, password string) (users []*models.User) {
	conn := s.Session(ctx)
	defer conn.Close()
	logger := log.GetContextLogger(ctx)
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 10, 0, false,
		fmt.Sprintf("(entryUUID=%s)", id), nil, nil,
	)
	dns, err := s.getDNSByReq(ldap.WithConnContext(ctx, conn), searchReq)
	if err != nil {
		level.Error(logger).Log("msg", "unknown error", "id", id, "err", err)
		return nil
	} else if len(dns) == 0 {
		level.Debug(logger).Log("msg", "no such user", "id", id)
	}

	for _, dn := range dns {
		if err = conn.Bind(dn, password); err != nil {
			level.Debug(logger).Log("msg", "incorrect password", "id", id, "err", err)
			continue
		}
		userInfo, err := s.getUserByDn(ctx, dn)
		if err != nil {
			level.Error(logger).Log("msg", "failed to get user info", "id", id, "err", err)
			return nil
		} else if userInfo == nil {
			level.Warn(logger).Log("msg", "failed to get user info", "id", id, "err", err)
			continue
		}
		users = append(users, userInfo)
	}
	return users
}

// VerifyPassword
//
//	@Description[en-US]: Verify password for user.
//	@Description[zh-CN]: 验证用户密码。
//	@param ctx 	context.Context
//	@param username 	string
//	@param password 	string
//	@return users	[]*models.User
func (s UserAndAppService) VerifyPassword(ctx context.Context, username string, password string) (users []*models.User) {
	conn := s.Session(ctx)
	defer conn.Close()
	logger := log.GetContextLogger(ctx)
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 10, 0, false,
		s.Options().ParseUserSearchFilter(username), nil, nil,
	)
	dns, err := s.getDNSByReq(ldap.WithConnContext(ctx, conn), searchReq)
	if err != nil {
		level.Error(logger).Log("msg", "unknown error", "username", username, "err", err)
		return nil
	} else if len(dns) == 0 {
		level.Debug(logger).Log("msg", "no such user", "username", username)
	}

	for _, dn := range dns {
		if err = conn.Bind(dn, password); err != nil {
			level.Debug(logger).Log("msg", "incorrect password", "username", username, "err", err)
			continue
		}
		userInfo, err := s.getUserByDn(ctx, dn)
		if err != nil {
			level.Error(logger).Log("msg", "failed to get user info", "username", username, "err", err)
			return nil
		} else if userInfo == nil {
			level.Warn(logger).Log("msg", "failed to get user info", "username", username, "err", err)
			continue
		}
		users = append(users, userInfo)
	}
	return users
}

func (s UserAndAppService) GetUsersById(ctx context.Context, ids []string) (users models.Users, err error) {
	for _, id := range ids {
		userInfo, err := s.GetUserInfoById(ctx, id)
		if err != nil && !strings.Contains(err.Error(), "Not Found") {
			return nil, err
		} else if userInfo == nil {
			continue
		}
		users = append(users, userInfo)
	}
	return users, nil
}
