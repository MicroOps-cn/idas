package ldapservice

import (
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"
	"github.com/tredoe/osutil/user/crypt/sha256_crypt"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/utils/httputil"
	"idas/pkg/utils/sets"
	"idas/pkg/utils/wrapper"
	"net/http"
	"strconv"
	"strings"
	"time"

	goldap "github.com/go-ldap/ldap"

	"idas/pkg/client/ldap"
	"idas/pkg/service/models"
)

const UserStatusName = "status"

func NewUserAndAppService(ctx context.Context, name string, client *ldap.Client) *UserAndAppService {
	return &UserAndAppService{name: name, Client: client}
}

type UserAndAppService struct {
	*ldap.Client
	name string
}

func (s UserAndAppService) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (*models.User, error) {
	conn := s.Session(ctx)

	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&(%s=*%s*)(%s=*%s*))",
			s.Options().GetAttrUsername(), username,
			s.Options().GetAttrEmail(), email),
		nil, nil,
	)
	return s.getUserInfoByReq(context.WithValue(ctx, global.LDAPConnName, conn), searchReq)
}

func (s UserAndAppService) ResetPassword(ctx context.Context, id string, password string) error {
	conn := s.Session(ctx)
	defer conn.Close()
	phash, err := hash([]byte(password))
	if err != nil {
		logger := logs.GetContextLogger(ctx)
		level.Error(logger).Log("failed to general password hash")
		return err
	}
	dn, err := s.getUserDnByEntryUUID(ctx, id)
	if err != nil {
		return err
	}
	req := goldap.NewModifyRequest(dn, nil)
	req.Replace("userPassword", []string{"{CRYPT}" + phash})
	return conn.Modify(req)
}

func (s UserAndAppService) UpdateLoginTime(_ context.Context, _ string) error {
	return nil
}

func (s UserAndAppService) AutoMigrate(ctx context.Context) error {
	return nil
}

func (s UserAndAppService) GetUsers(ctx context.Context, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{s.Options().ParseUserSearchFilter()}
	if status != models.UserStatusUnknown {
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
	if err != nil {
		return nil, 0, err
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
				CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Username:    entry.GetAttributeValue(s.Options().GetAttrUsername()),
			Email:       entry.GetAttributeValue(s.Options().GetAttrEmail()),
			PhoneNumber: entry.GetAttributeValue(s.Options().GetAttrUserPhoneNo()),
			FullName:    entry.GetAttributeValue(s.Options().GetAttrUserDisplayName()),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.UserStatus(wrapper.Must[int](httputil.NewValue(entry.GetAttributeValue(UserStatusName)).Default("0").Int())),
			Storage:     s.name,
		})
	}
	return
}

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

func (s UserAndAppService) getUserObjectClass(ctx context.Context, dn string) (objectClass []string, err error) {
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
	objectClass = ret.Entries[0].GetAttributeValues("objectClass")
	return objectClass, nil
}
func (s UserAndAppService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		id, ok := patchInfo["id"].(string)
		if !ok {
			return count, errors.ParameterError("unknown id")
		}
		var dn string
		dn, err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id)
		if err != nil {
			return count, err
		}
		objectClass, err := s.getUserObjectClass(context.WithValue(ctx, global.LDAPConnName, conn), dn)
		if !sets.New[string](objectClass...).Has("idasCore") {
			objectClass = append(objectClass, "idasCore")
		}

		req := goldap.NewModifyRequest(dn, nil)
		req.Replace("objectClass", objectClass)
		for name, value := range patchInfo {
			switch name {
			case "status":
				status := value.(float64)
				req.Replace(UserStatusName, []string{strconv.Itoa(int(status))})
			case "isDelete":
				//
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

func (s UserAndAppService) DeleteUsers(ctx context.Context, ids []string) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, id := range ids {
		var dn string
		if dn, err = s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id); err != nil {
			return count, err
		} else if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
			return
		}
		count++
	}
	return
}

func (s UserAndAppService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	dn, err := s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), user.Id)
	if err != nil {
		return nil, err
	}
	objectClass, err := s.getUserObjectClass(context.WithValue(ctx, global.LDAPConnName, conn), dn)
	if !sets.New[string](objectClass...).Has("idasCore") {
		objectClass = append(objectClass, "idasCore")
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"object_class", "username", "email", "phone_number", "full_name", "avatar", "status"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(dn, nil)
	var replace = []ldapUpdateColumn{
		{columnName: "object_class", ldapColumnName: "objectClass", val: objectClass},
		{columnName: "username", ldapColumnName: s.Options().GetAttrUsername(), val: []string{user.Username}},
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
		if err = conn.Modify(req); err != nil {
			return nil, err
		}
	}
	newUserInfo, err := s.getUserInfoByDn(context.WithValue(ctx, global.LDAPConnName, conn), dn)
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newUserInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newUserInfo, nil
}

func (s UserAndAppService) getUserInfoByReq(ctx context.Context, searchReq *goldap.SearchRequest) (*models.User, error) {
	conn := s.Session(ctx)
	searchReq.Attributes = []string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "entryUUID", "avatar", "createTimestamp", "modifyTimestamp", UserStatusName}
	ret, err := conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(ret.Entries) == 0 {
		return nil, errors.StatusNotFound("user")
	}
	userEntry := ret.Entries[0]
	return &models.User{
		Model: models.Model{
			Id:         userEntry.GetAttributeValue("entryUUID"),
			CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", userEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", userEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Username:    userEntry.GetAttributeValue(s.Options().GetAttrUsername()),
		Email:       userEntry.GetAttributeValue(s.Options().GetAttrEmail()),
		PhoneNumber: userEntry.GetAttributeValue(s.Options().GetAttrUserPhoneNo()),
		FullName:    userEntry.GetAttributeValue(s.Options().GetAttrUserDisplayName()),
		Avatar:      userEntry.GetAttributeValue("avatar"),
		Status:      models.UserStatus(wrapper.Must[int](httputil.NewValue(userEntry.GetAttributeValue(UserStatusName)).Default("0").Int())),
		Storage:     s.name,
	}, nil
}

func (s UserAndAppService) getUserInfoByDn(ctx context.Context, dn string) (*models.User, error) {
	conn := s.Session(ctx)
	searchReq := goldap.NewSearchRequest(
		dn, goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)
	return s.getUserInfoByReq(context.WithValue(ctx, global.LDAPConnName, conn), searchReq)
}

func (s UserAndAppService) GetUserInfoById(ctx context.Context, id string) (*models.User, error) {
	conn := s.Session(ctx)
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeSingleLevel, goldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(entryUUID=%s)", id), nil, nil,
	)
	return s.getUserInfoByReq(context.WithValue(ctx, global.LDAPConnName, conn), searchReq)
}

func (s UserAndAppService) getUserInfoByUsername(ctx context.Context, username string) (*models.User, error) {
	conn := s.Session(ctx)
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 2, 0, false,
		s.Options().ParseUserSearchFilter(username), nil, nil,
	)
	return s.getUserInfoByReq(context.WithValue(ctx, global.LDAPConnName, conn), searchReq)
}

func (s UserAndAppService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) == 0 && len(username) == 0 {
		return nil, errors.ParameterError("require id or username")
	}

	var userEntry *goldap.Entry

	if len(id) != 0 {
		return s.GetUserInfoById(context.WithValue(ctx, global.LDAPConnName, conn), id)
	}
	if userEntry == nil && len(username) != 0 {
		return s.getUserInfoByUsername(context.WithValue(ctx, global.LDAPConnName, conn), username)
	}
	return nil, errors.StatusNotFound("user")
}

func hash(password []byte) (string, error) {
	c := sha256_crypt.New()
	salt := string(uuid.NewV4().Bytes())
	return c.Generate(password, []byte("$5$"+salt))
}

func (s UserAndAppService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	logger := logs.GetContextLogger(ctx)
	conn := s.Session(ctx)
	defer conn.Close()
	dn := fmt.Sprintf("%s=%s,%s", s.Options().GetAttrUsername(), user.Username, s.Options().UserSearchBase)
	req := goldap.NewAddRequest(dn, nil)

	var attrs = map[string][]string{
		"mail":                               {user.Email},
		s.Options().GetAttrUserPhoneNo():     {user.PhoneNumber},
		s.Options().GetAttrUsername():        {user.Username},
		s.Options().GetAttrUserDisplayName(): {user.FullName},
		"sn":                                 {" "},
		"avatar":                             {user.Avatar},
		"status":                             {strconv.Itoa(int(user.Status))},
		"objectClass":                        {"idasCore", "inetOrgPerson", "organizationalPerson", "person", "top"},
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

	if err := conn.Add(req); err != nil {
		return nil, err
	}

	newUserInfo, err := s.getUserInfoByDn(context.WithValue(ctx, global.LDAPConnName, conn), dn)
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed: "+err.Error())
	} else if newUserInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed. ")
	}
	return newUserInfo, nil
}

func (s UserAndAppService) PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	id, ok := user["id"].(string)
	if !ok {
		return nil, errors.ParameterError("unknown id")
	}
	dn, err := s.getUserDnByEntryUUID(context.WithValue(ctx, global.LDAPConnName, conn), id)
	if err != nil {
		return nil, err
	}
	req := goldap.NewModifyRequest(dn, nil)
	var columns = []ldapUpdateColumn{
		{columnName: "username", ldapColumnName: s.Options().GetAttrUsername()},
		{columnName: "email", ldapColumnName: s.Options().GetAttrEmail()},
		{columnName: "phone_number", ldapColumnName: s.Options().GetAttrUserPhoneNo()},
		{columnName: "full_name", ldapColumnName: s.Options().GetAttrUserDisplayName()},
		{columnName: "avatar", ldapColumnName: "avatar"},
	}
	for _, column := range columns {
		if value, ok := user[column.columnName]; ok {
			switch val := value.(type) {
			case float64:
				req.Replace(column.ldapColumnName, []string{strconv.Itoa(int(val))})
			case string:
				req.Replace(column.ldapColumnName, []string{val})
			}
		}
	}

	if err := conn.Modify(req); err != nil {
		return nil, err
	}
	newUserInfo, err := s.GetUserInfoById(context.WithValue(ctx, global.LDAPConnName, conn), user["id"].(string))
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newUserInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newUserInfo, nil
}

func (s UserAndAppService) DeleteUser(ctx context.Context, id string) error {
	return wrapper.Error[int64](s.DeleteUsers(ctx, []string{id}))
}

func (s UserAndAppService) getDnsByReq(ctx context.Context, searchReq *goldap.SearchRequest) (dns []string, err error) {
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

func (s UserAndAppService) VerifyPassword(ctx context.Context, username string, password string) (users []*models.User) {
	conn := s.Session(ctx)
	defer conn.Close()
	logger := logs.GetContextLogger(ctx)
	searchReq := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 10, 0, false,
		s.Options().ParseUserSearchFilter(username), nil, nil,
	)
	dns, err := s.getDnsByReq(context.WithValue(ctx, global.LDAPConnName, conn), searchReq)
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
		userInfo, err := s.getUserInfoByDn(context.WithValue(ctx, global.LDAPConnName, conn), dn)
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
func (s UserAndAppService) GetAuthCodeByClientId() string {
	return s.name
}

func (s UserAndAppService) Name() string {
	return s.name
}
