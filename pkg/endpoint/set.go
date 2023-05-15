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

package endpoint

import (
	"context"
	"reflect"
	"strings"

	"github.com/MicroOps-cn/fuck/sets"
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

	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type UserEndpoints struct {
	GetUsers          endpoint.Endpoint `description:"Get user list" role:"admin|viewer" audit:"false"`
	DeleteUsers       endpoint.Endpoint `description:"Batch delete users" role:"admin" audit:"false"`
	PatchUsers        endpoint.Endpoint `description:"Batch modify user information (incremental)" role:"admin" audit:"true"`
	UpdateUser        endpoint.Endpoint `description:"Modify user information" role:"admin" audit:"true"`
	GetUserInfo       endpoint.Endpoint `description:"Get user details" role:"admin|viewer" audit:"false"`
	CreateUser        endpoint.Endpoint `description:"Create a user" role:"admin" audit:"true"`
	PatchUser         endpoint.Endpoint `description:"Modify user information (incremental)" role:"admin" audit:"true"`
	DeleteUser        endpoint.Endpoint `description:"Delete a user" role:"admin" audit:"true"`
	ForgotPassword    endpoint.Endpoint `auth:"false" audit:"true"`
	ResetPassword     endpoint.Endpoint `auth:"false" audit:"true"`
	CurrentUser       endpoint.Endpoint `auth:"false" audit:"false"`
	CreateUserKey     endpoint.Endpoint `description:"Create a user key-pair" role:"admin" audit:"true"`
	DeleteUserKey     endpoint.Endpoint `description:"Delete a user key-pair" role:"admin" audit:"true"`
	GetUserKeys       endpoint.Endpoint `description:"Get a user key-pairs" role:"admin|viewer" audit:"false"`
	CreateKey         endpoint.Endpoint `auth:"false" audit:"true"`
	CreateTOTPSecret  endpoint.Endpoint `auth:"false" audit:"false"`
	CreateTOTP        endpoint.Endpoint `auth:"false" audit:"true"`
	UnbindTOTP        endpoint.Endpoint `auth:"false" audit:"true"`
	SendLoginCaptcha  endpoint.Endpoint `auth:"false" audit:"true"`
	UpdateCurrentUser endpoint.Endpoint `auth:"false" audit:"true"`
	PatchCurrentUser  endpoint.Endpoint `auth:"false" audit:"true"`
	SendActivateMail  endpoint.Endpoint `description:"Send activation link to user mail" role:"admin" audit:"true"`
	ActivateAccount   endpoint.Endpoint `auth:"false" audit:"true"`
}

type AppEndpoints struct {
	PatchApps          endpoint.Endpoint `description:"Batch modify application information (incremental)" role:"admin" audit:"true"`
	DeleteApps         endpoint.Endpoint `description:"Batch delete applications" role:"admin" audit:"true"`
	GetAppInfo         endpoint.Endpoint `description:"Get application details" role:"admin|viewer" audit:"false"`
	CreateApp          endpoint.Endpoint `description:"Create an application" role:"admin" audit:"true"`
	UpdateApp          endpoint.Endpoint `description:"Modify application information" role:"admin" audit:"true"`
	PatchApp           endpoint.Endpoint `description:"Modify application information (incremental)" role:"admin" audit:"true"`
	DeleteApp          endpoint.Endpoint `description:"Delete a application" role:"admin" audit:"true"`
	GetApps            endpoint.Endpoint `description:"Get application list" role:"admin|viewer" audit:"false"`
	AppAuthentication  endpoint.Endpoint `auth:"false" audit:"false"`
	CreateAppKey       endpoint.Endpoint `description:"Create a app key-pair" role:"admin" audit:"true"`
	DeleteAppKey       endpoint.Endpoint `description:"Delete a app key-pair" role:"admin" audit:"true"`
	GetAppKeys         endpoint.Endpoint `description:"Get a app key-pairs" role:"admin" audit:"false"`
	GetCurrentUserApps endpoint.Endpoint `auth:"false" audit:"false"`
}

type SessionEndpoints struct {
	GetSessions              endpoint.Endpoint `description:"Get the user's session list" role:"admin" audit:"false"`
	DeleteSession            endpoint.Endpoint `description:"Delete the user's session" role:"admin" audit:"true"`
	GetCurrentUserSessions   endpoint.Endpoint `description:"Get current user session list" auth:"false" audit:"false"`
	DeleteCurrentUserSession endpoint.Endpoint `description:"Get current user session list" auth:"false" audit:"false"`
	UserLogin                endpoint.Endpoint `auth:"false" audit:"true"`
	UserLogout               endpoint.Endpoint `auth:"false" audit:"true"`
	GetSessionByToken        endpoint.Endpoint `auth:"false" audit:"false"`
	GetProxySessionByToken   endpoint.Endpoint `auth:"false" audit:"false"`
	OAuthTokens              endpoint.Endpoint `auth:"false" audit:"true"`
	OAuthAuthorize           endpoint.Endpoint `auth:"false" audit:"true"`
	Authentication           endpoint.Endpoint `auth:"false" audit:"false"`
	SessionRenewal           endpoint.Endpoint `auth:"false" audit:"false"`
}

type RoleEndpoints struct {
	GetPermissions endpoint.Endpoint `description:"Get permission list" role:"admin|viewer" audit:"false"`
	GetRoles       endpoint.Endpoint `description:"Get role list" role:"admin|viewer" audit:"false"`
	DeleteRoles    endpoint.Endpoint `description:"Batch delete roles" role:"admin"`
	CreateRole     endpoint.Endpoint `description:"Create a role" role:"admin"`
	UpdateRole     endpoint.Endpoint `description:"Modify role information" role:"admin"`
	DeleteRole     endpoint.Endpoint `description:"Delete a role" role:"admin"`
}

type PageEndpoints struct {
	GetPages   endpoint.Endpoint `description:"Get page list" role:"admin|viewer" audit:"false"`
	GetPage    endpoint.Endpoint `description:"Get page" role:"admin|viewer" audit:"false"`
	CreatePage endpoint.Endpoint `description:"Create a page" role:"admin"`
	UpdatePage endpoint.Endpoint `description:"Modify page information" role:"admin"`
	DeletePage endpoint.Endpoint `description:"Delete a page" role:"admin"`
	PatchPages endpoint.Endpoint `description:"Patch pages" role:"admin"`

	GetPageDatas   endpoint.Endpoint `description:"Get data list of page" role:"admin|viewer" audit:"false"`
	GetPageData    endpoint.Endpoint `description:"Get a data of page" role:"admin|viewer" audit:"false"`
	CreatePageData endpoint.Endpoint `description:"Create a data of page" role:"admin"`
	UpdatePageData endpoint.Endpoint `description:"Modify a data of page" role:"admin"`
	DeletePageData endpoint.Endpoint `description:"Delete a data of page" role:"admin"`
	PatchPageDatas endpoint.Endpoint `description:"Patch page data" role:"admin"`
}

type FileEndpoints struct {
	UploadFile   endpoint.Endpoint `name:"" description:"Upload files to the server" auth:"false" audit:"false"`
	DownloadFile endpoint.Endpoint `auth:"false" audit:"false"`
}

type ProxyEndpoints struct {
	ProxyRequest   endpoint.Endpoint `auth:"false" audit:"false"`
	GetProxyConfig endpoint.Endpoint `auth:"false" audit:"false"`
}

type ConfigEndpoints struct {
	GetSecurityConfig   endpoint.Endpoint `description:"Get security config." role:"admin" audit:"false"`
	PatchSecurityConfig endpoint.Endpoint `description:"Patch security config." role:"admin"`
}
type EventEndpoints struct {
	GetEvents               endpoint.Endpoint `description:"Get events." role:"admin" audit:"false"`
	GetEventLogs            endpoint.Endpoint `description:"Get event logs." role:"admin" audit:"false"`
	GetCurrentUserEvents    endpoint.Endpoint `description:"Get current user events." auth:"false" audit:"false"`
	GetCurrentUserEventLogs endpoint.Endpoint `description:"Get current user event logs." auth:"false" audit:"false"`
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
	ProxyEndpoints   `name:"Proxy" description:"Proxy"`
	PageEndpoints    `name:"Page" description:"Page"`
	ConfigEndpoints  `name:"Config" description:"System Config Manage"`
	EventEndpoints   `name:"Event" description:"Event"`
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
		if field.Type.Kind() == reflect.Struct {
			if auth := field.Tag.Get("auth"); len(auth) == 0 || auth == "true" {
				p.EnableAuth = true
			}
			p.Children = GetPermissionsDefine(field.Type)
			if len(p.Children) > 0 {
				ret = append(ret, &p)
				continue
			}
		} else if field.Type.Kind() == reflect.Func {
			if audit := field.Tag.Get("audit"); len(audit) == 0 || audit == "true" {
				p.EnableAudit = true
			}

			if auth := field.Tag.Get("auth"); len(auth) == 0 || auth == "true" {
				p.EnableAuth = true
			}
			if p.EnableAuth {
				if role := field.Tag.Get("role"); len(role) > 0 {
					p.Role = strings.Split(role, "|")
				}
			}
			ret = append(ret, &p)
		}
	}
	return ret
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(_ context.Context, svc service.Service, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	ps := Set{}.GetPermissionsDefine()
	eps := sets.New[string]()
	injectEndpoint := func(name string, ep endpoint.Endpoint) endpoint.Endpoint {
		if eps.Has(name) {
			panic("duplicate endpoint: " + name)
		}
		psd := ps.Get(name)
		if len(psd) == 0 {
			panic("endpoint not found: " + name)
		} else if len(psd) > 1 {
			panic("duplicate endpoint define: " + name)
		}
		eps.Insert(name)
		ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(ep)
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
		if otTracer != nil {
			ep = opentracing.TraceServer(otTracer, name)(ep)
		}
		if zipkinTracer != nil {
			ep = zipkin.TraceEndpoint(zipkinTracer, name)(ep)
		}
		ep = LoggingMiddleware(svc, name, ps)(ep)
		if duration != nil {
			ep = InstrumentingMiddleware(duration.With("method", "Concat"))(ep)
		}
		if psd[0].EnableAuth {
			ep = AuthorizationMiddleware(svc, name)(ep)
		}
		return ep
	}

	return Set{
		FileEndpoints: FileEndpoints{
			UploadFile:   injectEndpoint("UploadFile", MakeUploadFileEndpoint(svc)),
			DownloadFile: injectEndpoint("DownloadFile", MakeDownloadFileEndpoint(svc)),
		},
		UserEndpoints: UserEndpoints{
			CurrentUser:       injectEndpoint("CurrentUser", MakeCurrentUserEndpoint(svc)),
			ResetPassword:     injectEndpoint("ResetPassword", MakeResetUserPasswordEndpoint(svc)),
			ForgotPassword:    injectEndpoint("ForgotPassword", MakeForgotPasswordEndpoint(svc)),
			GetUsers:          injectEndpoint("GetUsers", MakeGetUsersEndpoint(svc)),
			DeleteUsers:       injectEndpoint("DeleteUsers", MakeDeleteUsersEndpoint(svc)),
			PatchUsers:        injectEndpoint("PatchUsers", MakePatchUsersEndpoint(svc)),
			UpdateUser:        injectEndpoint("UpdateUser", MakeUpdateUserEndpoint(svc)),
			GetUserInfo:       injectEndpoint("GetUserInfo", MakeGetUserInfoEndpoint(svc)),
			CreateUser:        injectEndpoint("CreateUser", MakeCreateUserEndpoint(svc)),
			PatchUser:         injectEndpoint("PatchUser", MakePatchUserEndpoint(svc)),
			DeleteUser:        injectEndpoint("DeleteUser", MakeDeleteUserEndpoint(svc)),
			CreateUserKey:     injectEndpoint("CreateUserKey", MakeCreateUserKeyEndpoint(svc)),
			DeleteUserKey:     injectEndpoint("DeleteUserKey", MakeDeleteUserKeyEndpoint(svc)),
			GetUserKeys:       injectEndpoint("GetUserKeys", MakeGetUserKeysEndpoint(svc)),
			CreateKey:         injectEndpoint("CreateKey", MakeCreateKeyEndpoint(svc)),
			CreateTOTPSecret:  injectEndpoint("CreateTOTPSecret", MakeCreateTOTPSecretEndpoint(svc)),
			CreateTOTP:        injectEndpoint("CreateTOTP", MakeCreateTOTPEndpoint(svc)),
			SendLoginCaptcha:  injectEndpoint("SendLoginCaptcha", MakeSendLoginCaptchaEndpoint(svc)),
			UpdateCurrentUser: injectEndpoint("UpdateCurrentUser", MakeUpdateCurrentUserEndpoint(svc)),
			PatchCurrentUser:  injectEndpoint("PatchCurrentUser", MakePatchCurrentUserEndpoint(svc)),
			ActivateAccount:   injectEndpoint("ActivateAccount", MakeActivateAccountEndpoint(svc)),
			SendActivateMail:  injectEndpoint("SendActivateMail", MakeSendActivationMailEndpoint(svc)),
		},
		SessionEndpoints: SessionEndpoints{
			GetSessions:              injectEndpoint("GetSessions", MakeGetSessionsEndpoint(svc)),
			GetCurrentUserSessions:   injectEndpoint("GetCurrentUserSessions", MakeGetCurrentUserSessionsEndpoint(svc)),
			DeleteCurrentUserSession: injectEndpoint("DeleteCurrentUserSession", MakeDeleteCurrentUserSessionEndpoint(svc)),
			DeleteSession:            injectEndpoint("DeleteSession", MakeDeleteSessionEndpoint(svc)),
			UserLogin:                injectEndpoint("UserLogin", MakeUserLoginEndpoint(svc)),
			Authentication:           injectEndpoint("Authentication", MakeAuthenticationEndpoint(svc)),
			UserLogout:               injectEndpoint("UserLogout", MakeUserLogoutEndpoint(svc)),
			GetSessionByToken:        injectEndpoint("GetSessionByToken", MakeGetSessionByTokenEndpoint(svc)),
			GetProxySessionByToken:   injectEndpoint("GetProxySessionByToken", MakeGetProxySessionByTokenEndpoint(svc)),
			OAuthTokens:              injectEndpoint("OAuthTokens", MakeOAuthTokensEndpoint(svc)),
			OAuthAuthorize:           injectEndpoint("OAuthAuthorize", MakeOAuthAuthorizeEndpoint(svc)),
		},
		AppEndpoints: AppEndpoints{
			GetApps:            injectEndpoint("GetApps", MakeGetAppsEndpoint(svc)),
			DeleteApps:         injectEndpoint("DeleteApps", MakeDeleteAppsEndpoint(svc)),
			PatchApps:          injectEndpoint("PatchApps", MakePatchAppsEndpoint(svc)),
			UpdateApp:          injectEndpoint("UpdateApp", MakeUpdateAppEndpoint(svc)),
			GetAppInfo:         injectEndpoint("GetAppInfo", MakeGetAppInfoEndpoint(svc)),
			CreateApp:          injectEndpoint("CreateApp", MakeCreateAppEndpoint(svc)),
			PatchApp:           injectEndpoint("PatchApp", MakePatchAppEndpoint(svc)),
			DeleteApp:          injectEndpoint("DeleteApp", MakeDeleteAppEndpoint(svc)),
			AppAuthentication:  injectEndpoint("AppAuthentication", MakeAppAuthenticationEndpoint(svc)),
			CreateAppKey:       injectEndpoint("CreateAppKey", MakeCreateAppKeyEndpoint(svc)),
			DeleteAppKey:       injectEndpoint("DeleteAppKey", MakeDeleteAppKeyEndpoint(svc)),
			GetAppKeys:         injectEndpoint("GetAppKeys", MakeGetAppKeysEndpoint(svc)),
			GetCurrentUserApps: injectEndpoint("GetCurrentUserApps", MakeGetCurrentUserAppsEndpoint(svc)),
		},
		RoleEndpoints: RoleEndpoints{
			GetRoles:       injectEndpoint("GetRoles", MakeGetRolesEndpoint(svc)),
			DeleteRoles:    injectEndpoint("DeleteRoles", MakeDeleteRolesEndpoint(svc)),
			CreateRole:     injectEndpoint("CreateRole", MakeCreateRoleEndpoint(svc)),
			UpdateRole:     injectEndpoint("UpdateRole", MakeUpdateRoleEndpoint(svc)),
			DeleteRole:     injectEndpoint("DeleteRole", MakeDeleteRoleEndpoint(svc)),
			GetPermissions: injectEndpoint("GetPermissions", MakeGetPermissionsEndpoint(svc)),
		},
		ProxyEndpoints: ProxyEndpoints{
			ProxyRequest:   injectEndpoint("ProxyRequest", MakeProxyRequestEndpoint(svc)),
			GetProxyConfig: injectEndpoint("GetProxyConfig", MakeGetProxyConfigEndpoint(svc)),
		},
		PageEndpoints: PageEndpoints{
			GetPages:   injectEndpoint("GetPages", MakeGetPagesEndpoint(svc)),
			GetPage:    injectEndpoint("GetPage", MakeGetPageEndpoint(svc)),
			CreatePage: injectEndpoint("CreatePage", MakeCreatePageEndpoint(svc)),
			UpdatePage: injectEndpoint("UpdatePage", MakeUpdatePageEndpoint(svc)),
			DeletePage: injectEndpoint("DeletePage", MakeDeletePageEndpoint(svc)),
			PatchPages: injectEndpoint("PatchPages", MakePatchPagesEndpoint(svc)),

			GetPageDatas:   injectEndpoint("GetPageDatas", MakeGetPageDatasEndpoint(svc)),
			GetPageData:    injectEndpoint("GetPageData", MakeGetPageDataEndpoint(svc)),
			CreatePageData: injectEndpoint("CreatePageData", MakeCreatePageDataEndpoint(svc)),
			UpdatePageData: injectEndpoint("UpdatePageData", MakeUpdatePageDataEndpoint(svc)),
			DeletePageData: injectEndpoint("DeletePageData", MakeDeletePageDataEndpoint(svc)),
			PatchPageDatas: injectEndpoint("PatchPageDatas", MakePatchPageDatasEndpoint(svc)),
		},
		ConfigEndpoints: ConfigEndpoints{
			GetSecurityConfig:   injectEndpoint("GetSecurityConfig", MakeGetSecurityConfigEndpoint(svc)),
			PatchSecurityConfig: injectEndpoint("PatchSecurityConfig", MakePatchSecurityConfigEndpoint(svc)),
		},
		EventEndpoints: EventEndpoints{
			GetEvents:               injectEndpoint("GetEvents", MakeGetEventsEndpoint(svc)),
			GetEventLogs:            injectEndpoint("GetEventLogs", MakeGetEventLogsEndpoint(svc)),
			GetCurrentUserEvents:    injectEndpoint("GetCurrentUserEvents", MakeGetCurrentUserEventsEndpoint(svc)),
			GetCurrentUserEventLogs: injectEndpoint("GetCurrentUserEventLogs", MakeGetCurrentUserEventLogsEndpoint(svc)),
		},
	}
}

func (s Set) GetPermissionsDefine() models.Permissions {
	return GetPermissionsDefine(reflect.TypeOf(s))
}
