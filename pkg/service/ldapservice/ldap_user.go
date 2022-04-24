package ldapservice

import (
	"context"
	"fmt"
	"idas/pkg/errors"
	"idas/pkg/global"
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

func NewUserAndAppService(name string, client *ldap.Client) *UserAndAppService {
	return &UserAndAppService{name: name, Client: client}
}

type UserAndAppService struct {
	*ldap.Client
	name string
}

func (s UserAndAppService) AutoMigrate(ctx context.Context) error {
	return nil
}

func (s UserAndAppService) GetUsers(ctx context.Context, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{fmt.Sprintf(s.Options().ParseUserSearchFilter())}
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
		[]string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "avatar", "createTimestamp", "modifyTimestamp", UserStatusName},
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
				Id:         entry.DN,
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

func (s UserAndAppService) GetUserObjectClass(ctx context.Context, id string) (objectClass []string, err error) {
	conn := s.Session(ctx)
	defer conn.Close()

	searchReq := goldap.NewSearchRequest(
		id,
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
		return nil, errors.StatusNotFound(fmt.Sprintf("user %s", id))
	}
	objectClass = ret.Entries[0].GetAttributeValues("objectClass")
	return objectClass, nil
}
func (s UserAndAppService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		dn, ok := patchInfo["id"].(string)
		if !ok {
			return count, errors.ParameterError("unknown id")
		}
		objectClass, err := s.GetUserObjectClass(context.WithValue(ctx, global.LDAPConnName, conn), dn)
		if !sets.New[string](objectClass...).Has("idasCore") {
			objectClass = append(objectClass, "idasCore")
		}
		if !strings.HasSuffix(dn, s.Options().UserSearchBase) {
			return count, errors.ParameterError("Illegal parameter id")
		}
		req := goldap.NewModifyRequest(dn, nil)
		req.Replace("objectClass", objectClass)
		for name, value := range patchInfo {
			switch name {
			case "status":
				status := value.(float64)
				req.Replace(UserStatusName, []string{strconv.Itoa(int(status))})
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

func (s UserAndAppService) DeleteUsers(ctx context.Context, id []string) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, dn := range id {
		if !strings.HasSuffix(dn, s.Options().UserSearchBase) {
			return count, errors.ParameterError("Illegal parameter id")
		}
		if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
			return
		}
		count++
	}
	return
}

func (s UserAndAppService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	if !strings.HasSuffix(user.Id, s.Options().UserSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	objectClass, err := s.GetUserObjectClass(context.WithValue(ctx, global.LDAPConnName, conn), user.Id)
	if !sets.New[string](objectClass...).Has("idasCore") {
		objectClass = append(objectClass, "idasCore")
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"object_class", "username", "email", "phone_number", "full_name", "avatar", "status"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(user.Id, nil)
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
	newUserInfo, err := s.GetUserInfo(context.WithValue(ctx, global.LDAPConnName, conn), user.Id, "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newUserInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newUserInfo, nil
}

func (s UserAndAppService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) == 0 && len(username) == 0 {
		return nil, errors.ParameterError("require id or username")
	}

	var userEntry *goldap.Entry

	if len(id) != 0 {
		if !strings.HasSuffix(id, s.Options().UserSearchBase) {
			return nil, errors.ParameterError("Illegal parameter id")
		}
		searchReq := goldap.NewSearchRequest(
			id,
			goldap.ScopeBaseObject, goldap.NeverDerefAliases, 1, 0, false,
			"(objectClass=*)",
			[]string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "avatar", "createTimestamp", "modifyTimestamp", UserStatusName},
			nil,
		)
		ret, err := conn.Search(searchReq)
		if err != nil {
			return nil, err
		}
		if len(ret.Entries) > 0 {
			userEntry = ret.Entries[0]
		}
	}
	if userEntry == nil && len(username) != 0 {
		searchReq := goldap.NewSearchRequest(
			s.Options().UserSearchBase,
			goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 1, 0, false,
			s.Options().ParseUserSearchFilter(username),
			[]string{s.Options().GetAttrUsername(), s.Options().GetAttrUserPhoneNo(), s.Options().GetAttrEmail(), s.Options().GetAttrUserDisplayName(), "avatar", "createTimestamp", "modifyTimestamp", UserStatusName},
			nil,
		)
		ret, err := conn.Search(searchReq)
		if err != nil {
			return nil, err
		}
		if len(ret.Entries) > 0 {
			userEntry = ret.Entries[0]
		}
	}

	if userEntry == nil {
		return nil, errors.StatusNotFound("user")
	}
	return &models.User{
		Model: models.Model{
			Id:         userEntry.DN,
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

func (s UserAndAppService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	user.Id = fmt.Sprintf("%s=%s,%s", s.Options().GetAttrUsername(), user.Username, s.Options().UserSearchBase)
	req := goldap.NewAddRequest(user.Id, nil)
	var attrs = map[string][]string{
		"mail":                               {user.Email},
		s.Options().GetAttrUserPhoneNo():     {user.PhoneNumber},
		s.Options().GetAttrUsername():        {user.Username},
		s.Options().GetAttrUserDisplayName(): {user.FullName},
		"avatar":                             {user.Avatar},
		"status":                             {strconv.Itoa(int(user.Status))},
		"objectClass":                        {"idasCore", "inetOrgPerson", "organizationalPerson", "person", "top"},
	}
	for name, value := range attrs {
		if value[0] != "" {
			req.Attribute(name, value)
		}
	}
	if err := conn.Add(req); err != nil {
		return nil, err
	}

	newUserInfo, err := s.GetUserInfo(context.WithValue(ctx, global.LDAPConnName, conn), user.Id, "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed: "+err.Error())
	} else if newUserInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed. ")
	}
	return newUserInfo, nil
}

func (s UserAndAppService) PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error) {
	id, ok := user["id"].(string)
	if !ok {
		return nil, errors.ParameterError("unknown id")
	}
	if !strings.HasSuffix(id, s.Options().UserSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	req := goldap.NewModifyRequest(id, nil)
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

	conn := s.Session(ctx)
	defer conn.Close()
	if err := conn.Modify(req); err != nil {
		return nil, err
	}
	newUserInfo, err := s.GetUserInfo(context.WithValue(ctx, global.LDAPConnName, conn), user["id"].(string), "")
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

func (s UserAndAppService) VerifyPassword(ctx context.Context, username string, password string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	userInfo, err := s.GetUserInfo(context.WithValue(ctx, global.LDAPConnName, conn), "", username)
	if err != nil {
		return nil, err
	} else if userInfo == nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if err = conn.Bind(userInfo.Id, password); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}
	return userInfo, nil
}

func (s UserAndAppService) Name() string {
	return s.name
}
