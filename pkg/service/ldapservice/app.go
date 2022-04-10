package ldapservice

import (
	"context"
	"fmt"
	goldap "github.com/go-ldap/ldap"
	"idas/pkg/client/ldap"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service/models"
	"idas/pkg/utils/sets"
	"idas/pkg/utils/wrapper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AppService struct {
	name string
	*ldap.Client
}

func (s AppService) AutoMigrate(ctx context.Context) error {
	return nil
}

func (s AppService) Name() string {
	return s.name
}

func (s AppService) GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error) {
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
		[]string{"description", "cn", "avatar", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
		nil,
	)
	ret, err := conn.Search(req)
	if err != nil {
		return nil, 0, err
	}
	total = int64(len(ret.Entries))
	for _, entry := range ret.Entries[(current-1)*pageSize : current*pageSize] {
		apps = append(apps, &models.App{
			Model: models.Model{
				Id:         entry.DN,
				CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("createTimestamp"))),
				UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", entry.GetAttributeValue("modifyTimestamp"))),
			},
			Name:        entry.GetAttributeValue("cn"),
			Description: entry.GetAttributeValue("description"),
			Avatar:      entry.GetAttributeValue("avatar"),
			Status:      models.GroupStatus(wrapper.Must[int](strconv.Atoi(entry.GetAttributeValue(GroupStatusName)))),
			GrantMode:   models.GrantMode(wrapper.Must[int](strconv.Atoi(entry.GetAttributeValue("grantMode")))),
			GrantType:   models.GrantType(entry.GetAttributeValue("grantType")),
			Storage:     s.name,
		})
	}
	return
}

const GroupStatusName = "status"

func (s AppService) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, patchInfo := range patch {
		dn, ok := patchInfo["id"].(string)
		if !ok {
			return total, errors.ParameterError("unknown id")
		}
		if !strings.HasSuffix(dn, s.Options().GroupSearchBase) {
			return total, errors.ParameterError("Illegal parameter id")
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

func (s AppService) DeleteApps(ctx context.Context, id []string) (total int64, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	for _, dn := range id {
		if !strings.HasSuffix(dn, s.Options().GroupSearchBase) {
			return total, errors.ParameterError("Illegal parameter id")
		}
		if err = conn.Del(goldap.NewDelRequest(dn, nil)); err != nil {
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

func (s AppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()

	if !strings.HasSuffix(app.Id, s.Options().GroupSearchBase) {
		return nil, errors.ParameterError("Illegal parameter id")
	}
	if len(updateColumns) == 0 {
		updateColumns = []string{"name", "description", "avatar", "grant_type", "grant_mode", "status"}
	}
	columns := sets.New[string](updateColumns...)
	req := goldap.NewModifyRequest(app.Id, nil)
	var replace = []ldapUpdateColumn{
		{columnName: "name", ldapColumnName: "cn", val: []string{app.Name}},
		{columnName: "description", ldapColumnName: "description", val: []string{app.Description}},
		{columnName: "avatar", ldapColumnName: "avatar", val: []string{app.Avatar}},
		{columnName: "grant_type", ldapColumnName: "grantType", val: []string{string(app.GrantType)}},
		{columnName: "grant_mode", ldapColumnName: "grantMode", val: []string{strconv.Itoa(int(app.GrantMode))}},
		{columnName: "status", ldapColumnName: "status", val: []string{strconv.Itoa(int(app.Status))}},
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
	newAppInfo, err := s.GetAppInfo(context.WithValue(ctx, global.LDAPConnName, conn), app.Id, "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed: "+err.Error())
	} else if newAppInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been modified successfully, but the query result failed. ")
	}
	return newAppInfo, nil
}

func (s AppService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	conn := s.Session(ctx)
	defer conn.Close()
	if len(id) == 0 && len(name) == 0 {
		return nil, errors.ParameterError("require id or appname")
	}

	var appEntry *goldap.Entry

	if len(id) != 0 {
		if !strings.HasSuffix(id, s.Options().GroupSearchBase) {
			return nil, errors.ParameterError("Illegal parameter id")
		}
		searchReq := goldap.NewSearchRequest(
			id,
			goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 1, 0, false,
			"(objectClass=*)",
			[]string{"description", "cn", "avatar", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
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
			s.Options().GroupSearchFilter,
			goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 1, 0, false,
			s.Options().ParseGroupSearchFilter(name),
			[]string{"description", "cn", "avatar", "createTimestamp", "modifyTimestamp", "grantMode", "grantType", GroupStatusName},
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
	return &models.App{
		Model: models.Model{
			Id:         appEntry.DN,
			CreateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("createTimestamp"))),
			UpdateTime: wrapper.Must[time.Time](time.Parse("20060102150405Z", appEntry.GetAttributeValue("modifyTimestamp"))),
		},
		Name:        appEntry.GetAttributeValue("cn"),
		Description: appEntry.GetAttributeValue("description"),
		Avatar:      appEntry.GetAttributeValue("avatar"),
		Status:      models.GroupStatus(wrapper.Must[int](strconv.Atoi(appEntry.GetAttributeValue(GroupStatusName)))),
		GrantMode:   models.GrantMode(wrapper.Must[int](strconv.Atoi(appEntry.GetAttributeValue("grantMode")))),
		GrantType:   models.GrantType(appEntry.GetAttributeValue("grantType")),
		Storage:     s.name,
	}, nil
}

func (s AppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	conn := s.Session(ctx)
	defer conn.Close()
	app.Id = fmt.Sprintf("cn=%s,%s", app.Name, s.Options().GroupSearchBase)
	req := goldap.NewAddRequest(app.Id, nil)
	var attrs = map[string][]string{
		"description": {app.Description},
		"avatar":      {app.Avatar},
		"grantType":   {string(app.GrantType)},
		"grantMode":   {strconv.Itoa(int(app.GrantMode))},
		"status":      {strconv.Itoa(int(app.Status))},
	}
	for name, value := range attrs {
		req.Attribute(name, value)
	}
	if err := conn.Add(req); err != nil {
		return nil, err
	}

	newAppInfo, err := s.GetAppInfo(context.WithValue(ctx, global.LDAPConnName, conn), app.Id, "")
	if err != nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed: "+err.Error())
	} else if newAppInfo == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "Internal Server Error. It may have been created successfully, but the query result failed. ")
	}
	return newAppInfo, nil
}

func (s AppService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
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

func (s AppService) DeleteApp(ctx context.Context, id string) (err error) {
	return wrapper.Error[int64](s.DeleteApps(ctx, []string{id}))
}

func NewAppService(name string, client *ldap.Client) *AppService {
	return &AppService{name: name, Client: client}
}
