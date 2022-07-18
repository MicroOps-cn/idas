package endpoint

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"idas/pkg/service"
	"idas/pkg/service/models"
	"idas/pkg/utils/sets"
	"reflect"
	"strings"
)

type UserEndpoints struct {
	GetUsers       endpoint.Endpoint `description:"Get user list" role:"admin|viewer"`
	DeleteUsers    endpoint.Endpoint `description:"Batch delete users" role:"admin"`
	PatchUsers     endpoint.Endpoint `description:"Batch modify user information (incremental)" role:"admin"`
	UpdateUser     endpoint.Endpoint `description:"Modify user information" role:"admin"`
	GetUserInfo    endpoint.Endpoint `description:"Get user details" role:"admin|viewer"`
	CreateUser     endpoint.Endpoint `description:"Create a user" role:"admin"`
	PatchUser      endpoint.Endpoint `description:"Modify user information (incremental)" role:"admin"`
	DeleteUser     endpoint.Endpoint `description:"Delete a user" role:"admin"`
	GetUserSource  endpoint.Endpoint `description:"Get the data source that stores user information" role:"admin|viewer"`
	ForgotPassword endpoint.Endpoint `auth:"false"`
	ResetPassword  endpoint.Endpoint `auth:"false"`
	CurrentUser    endpoint.Endpoint `auth:"false"`
}

type AppEndpoints struct {
	PatchApps    endpoint.Endpoint `description:"Batch modify application information (incremental)" role:"admin"`
	DeleteApps   endpoint.Endpoint `description:"Batch delete applications" role:"admin"`
	GetAppSource endpoint.Endpoint `description:"Get the data source that stores applications information" role:"admin|viewer"`
	GetAppInfo   endpoint.Endpoint `description:"Get application details" role:"admin|viewer"`
	CreateApp    endpoint.Endpoint `description:"Create an application" role:"admin"`
	UpdateApp    endpoint.Endpoint `description:"Modify application information" role:"admin"`
	PatchApp     endpoint.Endpoint `description:"Modify application information (incremental)" role:"admin"`
	DeleteApp    endpoint.Endpoint `description:"Delete a application" role:"admin"`
	GetApps      endpoint.Endpoint `description:"Get application list" role:"admin|viewer"`
}

type SessionEndpoints struct {
	GetSessions     endpoint.Endpoint `description:"Get the user's session list" role:"admin"`
	DeleteSession   endpoint.Endpoint `description:"Delete the user's session" role:"admin"`
	UserLogin       endpoint.Endpoint `auth:"false"`
	UserLogout      endpoint.Endpoint `auth:"false"`
	GetLoginSession endpoint.Endpoint `auth:"false"`
	OAuthTokens     endpoint.Endpoint `auth:"false"`
	OAuthAuthorize  endpoint.Endpoint `auth:"false"`
	Authentication  endpoint.Endpoint `auth:"false"`
}

type RoleEndpoints struct {
	GetPermissions endpoint.Endpoint `description:"Get permission list" role:"admin|viewer"`
	GetRoles       endpoint.Endpoint `description:"Get role list" role:"admin|viewer"`
	DeleteRoles    endpoint.Endpoint `description:"Batch delete roles" role:"admin"`
	CreateRole     endpoint.Endpoint `description:"Create a role" role:"admin"`
	UpdateRole     endpoint.Endpoint `description:"Modify role information" role:"admin"`
	DeleteRole     endpoint.Endpoint `description:"Delete a role" role:"admin"`
}

type FileEndpoints struct {
	UploadFile   endpoint.Endpoint `name:"" description:"Upload files to the server" role:"admin"`
	DownloadFile endpoint.Endpoint `name:"" description:"Download/view files" role:"admin|viewer"`
}

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	UserEndpoints    `name:"User" description:"User management"`
	SessionEndpoints `name:"Session" description:"User session management"`
	AppEndpoints     `name:"App" description:"Application management"`
	RoleEndpoints    `name:"Role" description:"Role of current platform"`
	FileEndpoints    `name:"File" description:"File"`
}

func GetPermissionsDefine(typeOf reflect.Type) models.Permissions {
	var ret models.Permissions
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	for i := 0; i < typeOf.NumField(); i++ {
		var p models.Permission
		field := typeOf.Field(i)
		if p.Name = field.Tag.Get("name"); len(p.Name) == 0 {
			p.Name = field.Name
		}

		p.Description = field.Tag.Get("description")
		if auth := field.Tag.Get("auth"); len(auth) == 0 || auth == "true" {
			p.EnableAuth = true
		}
		p.Role = strings.Split(field.Tag.Get("role"), "|")
		if field.Type.Kind() == reflect.Struct {
			p.Children = GetPermissionsDefine(field.Type)
			if len(p.Children) > 0 {
				ret = append(ret, &p)
				continue
			}
		} else if field.Type.Kind() == reflect.Func {
			ret = append(ret, &p)
		}
	}
	return ret
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(ctx context.Context, svc service.Service, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	ps := Set{}.GetPermissionsDefine()
	var eps = sets.New[string]()
	var injectEndpoint = func(name string, ep endpoint.Endpoint) endpoint.Endpoint {
		if eps.Has(name) {
			panic("duplicate endpoint: " + name)
		}
		if count := len(ps.Get(name)); count == 0 {
			panic("endpoint not found: " + name)
		} else if count > 1 {
			panic("duplicate endpoint define: " + name)
		}
		eps.Insert(name)
		ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(ep)
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
		ep = opentracing.TraceServer(otTracer, name)(ep)
		if zipkinTracer != nil {
			ep = zipkin.TraceEndpoint(zipkinTracer, name)(ep)
		}
		ep = LoggingMiddleware(name)(ep)
		ep = InstrumentingMiddleware(duration.With("method", "Concat"))(ep)
		ep = AuthorizationMiddleware(svc, name)(ep)
		return ep
	}

	return Set{
		FileEndpoints: FileEndpoints{
			UploadFile:   injectEndpoint("UploadFile", MakeUploadFileEndpoint(svc)),
			DownloadFile: injectEndpoint("DownloadFile", MakeDownloadFileEndpoint(svc)),
		},
		UserEndpoints: UserEndpoints{
			CurrentUser:    injectEndpoint("CurrentUser", MakeCurrentUserEndpoint(svc)),
			ResetPassword:  injectEndpoint("ResetPassword", MakeResetUserPasswordEndpoint(svc)),
			ForgotPassword: injectEndpoint("ForgotPassword", MakeForgotPasswordEndpoint(svc)),
			GetUsers:       injectEndpoint("GetUsers", MakeGetUsersEndpoint(svc)),
			DeleteUsers:    injectEndpoint("DeleteUsers", MakeDeleteUsersEndpoint(svc)),
			PatchUsers:     injectEndpoint("PatchUsers", MakePatchUsersEndpoint(svc)),
			UpdateUser:     injectEndpoint("UpdateUser", MakeUpdateUserEndpoint(svc)),
			GetUserInfo:    injectEndpoint("GetUserInfo", MakeGetUserInfoEndpoint(svc)),
			CreateUser:     injectEndpoint("CreateUser", MakeCreateUserEndpoint(svc)),
			PatchUser:      injectEndpoint("PatchUser", MakePatchUserEndpoint(svc)),
			DeleteUser:     injectEndpoint("DeleteUser", MakeDeleteUserEndpoint(svc)),
			GetUserSource:  injectEndpoint("GetUserSource", MakeGetUserSourceRequestEndpoint(svc)),
		},
		SessionEndpoints: SessionEndpoints{
			GetSessions:     injectEndpoint("GetSessions", MakeGetSessionsEndpoint(svc)),
			DeleteSession:   injectEndpoint("DeleteSession", MakeDeleteSessionEndpoint(svc)),
			UserLogin:       injectEndpoint("UserLogin", MakeUserLoginEndpoint(svc)),
			Authentication:  injectEndpoint("Authentication", MakeAuthenticationEndpoint(svc)),
			UserLogout:      injectEndpoint("UserLogout", MakeUserLogoutEndpoint(svc)),
			GetLoginSession: injectEndpoint("GetLoginSession", MakeGetLoginSessionEndpoint(svc)),
			OAuthTokens:     injectEndpoint("OAuthTokens", MakeOAuthTokensEndpoint(svc)),
			OAuthAuthorize:  injectEndpoint("OAuthAuthorize", MakeOAuthAuthorizeEndpoint(svc)),
		},
		AppEndpoints: AppEndpoints{
			GetApps:      injectEndpoint("GetApps", MakeGetAppsEndpoint(svc)),
			DeleteApps:   injectEndpoint("DeleteApps", MakeDeleteAppsEndpoint(svc)),
			PatchApps:    injectEndpoint("PatchApps", MakePatchAppsEndpoint(svc)),
			UpdateApp:    injectEndpoint("UpdateApp", MakeUpdateAppEndpoint(svc)),
			GetAppInfo:   injectEndpoint("GetAppInfo", MakeGetAppInfoEndpoint(svc)),
			CreateApp:    injectEndpoint("CreateApp", MakeCreateAppEndpoint(svc)),
			PatchApp:     injectEndpoint("PatchApp", MakePatchAppEndpoint(svc)),
			DeleteApp:    injectEndpoint("DeleteApp", MakeDeleteAppEndpoint(svc)),
			GetAppSource: injectEndpoint("GetAppSource", MakeGetAppSourceRequestEndpoint(svc)),
		},
		RoleEndpoints: RoleEndpoints{
			GetRoles:       injectEndpoint("GetRoles", MakeGetRolesEndpoint(svc)),
			DeleteRoles:    injectEndpoint("DeleteRoles", MakeDeleteRolesEndpoint(svc)),
			CreateRole:     injectEndpoint("CreateRole", MakeCreateRoleEndpoint(svc)),
			UpdateRole:     injectEndpoint("UpdateRole", MakeUpdateRoleEndpoint(svc)),
			DeleteRole:     injectEndpoint("DeleteRole", MakeDeleteRoleEndpoint(svc)),
			GetPermissions: injectEndpoint("GetPermissions", MakeGetPermissionsEndpoint(svc)),
		},
	}
}

func (s Set) GetPermissionsDefine() models.Permissions {
	return GetPermissionsDefine(reflect.TypeOf(s))
}
