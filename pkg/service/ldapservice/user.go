package ldapservice

import (
	"context"
	"fmt"
	"idas/pkg/errors"
	"idas/pkg/global"
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

func NewUserService(name string, client *ldap.Client) *UserService {
	return &UserService{name: name, Client: client}
}

type UserService struct {
	*ldap.Client
	name string
}

func (s UserService) AutoMigrate(ctx context.Context) error {
	return nil
}

func (s UserService) GetUsers(ctx context.Context, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	filters := []string{fmt.Sprintf(s.Options().ParseUserSearchFilter())}
	if status != models.UserStatusUnknown {
		filters = append(filters, fmt.Sprintf("(%s=%d)", UserStatusName, status))
	}
	if len(keywords) > 0 {
		filters = append(filters, fmt.Sprintf("(|(uid=*%s*)(sn=*%s*)(telephoneNumber=*%s*)(email=*%s*))", keywords, keywords, keywords, keywords))
	}
	var filter string
	if len(filters) >= 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(filters, ""))
	}
	req := goldap.NewSearchRequest(
		s.Options().UserSearchBase,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"uid", "email", "cn", "avatar", "createTimestamp", "modifyTimestamp", "telephoneNumber", UserStatusName},
		nil,
	)
	ret, err := conn.Search(req)
	if err != nil {
		return nil, 0, err
	}
	total = int64(len(ret.Entries))
	for _, entry := range ret.Entries[(current-1)*pageSize : current*pageSize] {
		users = append(users, &models.User{
			Model: models.Model{
				Id:         entry.DN,
				CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Username:    entry.GetAttributeValue("uid"),
			Email:       entry.GetAttributeValue("email"),
			PhoneNumber: entry.GetAttributeValue("telephoneNumber"),
			FullName:    entry.GetAttributeValue("cn"),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.UserStatus(wrapper.Must[int](strconv.Atoi(entry.GetAttributeValue(UserStatusName)))),
			Storage:     s.name,
		})
	}
	return
}

func (s UserService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		dn, ok := patchInfo["id"].(string)
		if !ok {
			return count, errors.ParameterError("unknown id")
		}
		if !strings.HasSuffix(dn, s.Options().UserSearchBase) {
			return count, errors.ParameterError("Illegal parameter id")
		}
		req := goldap.NewModifyRequest(dn, nil)
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

func (s UserService) DeleteUsers(ctx context.Context, id []string) (count int64, err error) {
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

func (s UserService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	if !strings.HasSuffix(user.Id, s.Options().UserSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"username", "email", "phone_number", "full_name", "avatar", "status"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(user.Id, nil)
	var replace = []ldapUpdateColumn{
		{columnName: "username", ldapColumnName: "uid", val: []string{user.Username}},
		{columnName: "email", ldapColumnName: "email", val: []string{user.Avatar}},
		{columnName: "phone_number", ldapColumnName: "telephoneNumber", val: []string{user.PhoneNumber}},
		{columnName: "full_name", ldapColumnName: "cn", val: []string{user.FullName}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{user.Avatar}},
		{columnName: "status", ldapColumnName: "status", val: []string{strconv.Itoa(int(user.Status))}},
	}

	for _, value := range replace {
		if columns.Has(value.columnName) {
			req.Replace(value.ldapColumnName, value.val)
		}
	}
	if len(req.Changes) > 0 {
		if err := conn.Modify(req); err != nil {
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

func (s UserService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
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
			goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 1, 0, false,
			"(objectClass=*)",
			[]string{"uid", "email", "cn", "avatar", "createTimestamp", "modifyTimestamp", "telephoneNumber", UserStatusName},
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
			[]string{"uid", "email", "cn", "avatar", "createTimestamp", "modifyTimestamp", "telephoneNumber", UserStatusName},
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
		Username:    userEntry.GetAttributeValue("uid"),
		Email:       userEntry.GetAttributeValue("email"),
		PhoneNumber: userEntry.GetAttributeValue("telephoneNumber"),
		FullName:    userEntry.GetAttributeValue("cn"),
		Avatar:      userEntry.GetAttributeValue("avatar"),
		Status:      models.UserStatus(wrapper.Must[int](strconv.Atoi(userEntry.GetAttributeValue(UserStatusName)))),
		Storage:     s.name,
	}, nil
}

func (s UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	user.Id = fmt.Sprintf("uid=%s,%s", user.Username, s.Options().UserSearchBase)
	req := goldap.NewAddRequest(user.Id, nil)
	var attrs = map[string][]string{
		"uid":             {user.Username},
		"email":           {user.Email},
		"telephoneNumber": {user.PhoneNumber},
		"cn":              {user.FullName},
		"avatar":          {user.Avatar},
		"status":          {strconv.Itoa(int(user.Status))},
	}
	for name, value := range attrs {
		req.Attribute(name, value)
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

func (s UserService) PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error) {
	id, ok := user["id"].(string)
	if !ok {
		return nil, errors.ParameterError("unknown id")
	}
	if !strings.HasSuffix(id, s.Options().UserSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	req := goldap.NewModifyRequest(id, nil)
	columns := sets.New[string]("username", "email", "phone_number", "full_name", "avatar")
	for name, value := range user {
		if !columns.Has(name) {
			continue
		}
		switch val := value.(type) {
		case float64:
			req.Replace(UserStatusName, []string{strconv.Itoa(int(val))})
		case string:
			req.Replace(UserStatusName, []string{val})
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

func (s UserService) DeleteUser(ctx context.Context, id string) error {
	return wrapper.Error[int64](s.DeleteUsers(ctx, []string{id}))
}

func (s UserService) VerifyPassword(ctx context.Context, username string, password string) (*models.User, error) {
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

func (s UserService) Name() string {
	return s.name
}
