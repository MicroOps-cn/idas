package transport

import (
	"context"
	"reflect"
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	kitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"idas/pkg/endpoint"
	"idas/pkg/global"
)

func WrapHTTPHandler(h *httptransport.Server) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		request := req.Request.WithContext(context.WithValue(context.WithValue(ctx, global.RestfulResponseContextName, resp), global.RestfulRequestContextName, req))
		h.ServeHTTP(resp, request)
	}
}

func NewKitHTTPServer[RequestType any](dp kitendpoint.Endpoint, options []httptransport.ServerOption) restful.RouteFunction {
	return WrapHTTPHandler(httptransport.NewServer(
		dp,
		decodeHTTPRequest[RequestType],
		encodeHTTPResponse,
		options...,
	))
}

const QueryTypeKey = "__query_type__"

func NewWebService(rootPath string, gv schema.GroupVersion, doc string) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(rootPath + "/" + gv.Version + "/" + gv.Group).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).Doc(doc)
	return &webservice
}

func NewSimpleWebService(rootPath string, doc string) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(rootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).Doc(doc)
	return &webservice
}

const rootPath = "/api"

func StructToQueryParams(obj interface{}, nameFilter ...string) []*restful.Parameter {
	var params []*restful.Parameter
	typeOfObj := reflect.TypeOf(obj)
	valueOfObj := reflect.ValueOf(obj)
	// 通过 #NumField 获取结构体字段的数量
loopObjFields:
	for i := 0; i < typeOfObj.NumField(); i++ {
		field := typeOfObj.Field(i)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			params = append(params, StructToQueryParams(valueOfObj.Field(i).Interface(), nameFilter...)...)
		} else {
			if len(nameFilter) > 0 {
				for _, name := range nameFilter {
					if name == field.Name {
						goto handleField
					}
				}
				continue loopObjFields
			}
		handleField:
			jsonTag := strings.Split(field.Tag.Get("json"), ",")
			if len(jsonTag) > 0 {
				params = append(params, restful.QueryParameter(
					jsonTag[0],
					field.Tag.Get("description")).DataType(field.Type.String()),
				)
			}
		}
	}
	return params
}

//UserService User Manager Service for restful Http container
func UserService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "users", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "User Manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetUsersRequest](endpoints.GetUsers, options)).
		Operation("getUsers").
		Doc("获取用户列表").
		Params(StructToQueryParams(endpoint.GetUsersRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchUsersRequest](endpoints.PatchUsers, options)).
		Operation("patchUsers").
		Reads([]map[string]endpoint.PatchUsersRequest{}).
		Doc("批量更新用户信息（增量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteUsersRequest](endpoints.DeleteUsers, options)).
		Operation("deleteUsers").
		Doc("批量删除用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateUserRequest](endpoints.CreateUser, options)).
		Operation("createUser").
		Doc("创建用户").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetUserRequest](endpoints.GetUserInfo, options)).
		Operation("getUserInfo").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("获取用户信息").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateUserRequest](endpoints.UpdateUser, options)).
		Operation("updateUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("更新用户信息（全量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchUserRequest](endpoints.PatchUser, options)).
		Operation("patchUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("更新用户信息（增量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteUserRequest](endpoints.DeleteUser, options)).
		Operation("deleteUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("删除用户").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.POST("/forgotPassword").
		To(NewKitHTTPServer[endpoint.ForgotUserPasswordRequest](endpoints.ForgotPassword, options)).
		Operation("forgotPassword").
		Doc("忘记用户密码").
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	v1ws.Route(v1ws.POST("/resetPassword").
		To(NewKitHTTPServer[endpoint.ResetUserPasswordRequest](endpoints.ResetPassword, options)).
		Operation("resetPassword").
		Doc("重置用户密码").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/source").
		To(NewKitHTTPServer[endpoint.GetUserSourceRequest](endpoints.GetUserSource, options)).
		Operation("getUserSource").
		Doc("获取用户存储源").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	return tag, []*restful.WebService{v1ws}
}

func AppService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "apps", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "Application manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetAppsRequest](endpoints.GetApps, options)).
		Operation("getApps").
		Doc("获取应用列表").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchAppsRequest](endpoints.PatchApps, options)).
		Operation("patchApps").
		Doc("批量更新应用信息（增量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteAppsRequest](endpoints.DeleteApps, options)).
		Operation("deleteApps").
		Doc("批量删除应用").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateAppRequest](endpoints.CreateApp, options)).
		Operation("createApp").
		Doc("创建应用").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/source").
		To(NewKitHTTPServer[endpoint.GetAppSourceRequest](endpoints.GetAppSource, options)).
		Operation("getAppSource").
		Doc("获取应用存储源").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetAppRequest](endpoints.GetAppInfo, options)).
		Operation("getAppInfo").
		Doc("获取应用信息").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateAppRequest](endpoints.UpdateApp, options)).
		Operation("updateApp").
		Doc("更新应用信息（全量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchAppRequest](endpoints.PatchApp, options)).
		Operation("patchApp").
		Doc("更新应用信息（增量）").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteAppRequest](endpoints.DeleteApp, options)).
		Operation("deleteApp").
		Doc("删除应用").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	return tag, []*restful.WebService{v1ws}
}

func FileService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "files", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "File manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.FileUploadRequest](endpoints.UploadFile, options)).
		Operation("uploadFile").
		Consumes("multipart/form-data").Doc("上传文件").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.FileDownloadRequest](endpoints.DownloadFile, options)).
		Operation("downloadFile").
		Doc("下载/查看文件").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	return tag, []*restful.WebService{v1ws}
}

func SessionService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "sessions", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "File manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetSessionsRequest](endpoints.GetSessions, options)).
		Operation("getSessions").
		Doc("获取会话列表").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteSessionRequest](endpoints.DeleteSession, options)).
		Operation("deleteSession").
		Doc("会话过期").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	return tag, []*restful.WebService{v1ws}
}

func OAuthService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "oauth", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "File manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	// https://www.ruanyifeng.com/blog/2019/04/oauth-grant-types.html
	v1ws.Route(v1ws.POST("/token").
		To(NewKitHTTPServer[endpoint.OAuthTokenRequest](endpoints.OAuthTokens, options)).
		Operation("oAuthTokens").
		Doc("获取令牌").
		Metadata(global.MetaNeedLogin, false).
		Consumes("application/x-www-form-urlencoded").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.POST("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("应用授权").
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("应用授权").
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	return tag, []*restful.WebService{v1ws}
}

func UserAuthService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "user", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, "File manager")
	v1ws.Filter(HTTPLoginAuthentication(endpoints))

	v1ws.Route(v1ws.POST("/login").
		To(NewKitHTTPServer[endpoint.UserLoginRequest](endpoints.UserLogin, options)).
		Operation("userLogin").
		Doc("用户登陆").
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.POST("/logout").
		To(NewKitHTTPServer[endpoint.UserLogoutRequest](endpoints.UserLogout, options)).
		Operation("userLogout").
		Doc("用户退出登录").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.CurrentUserRequest](endpoints.CurrentUser, options)).
		Operation("currentUser").
		Doc("获取当前登陆用户信息").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	return tag, []*restful.WebService{v1ws}
}
