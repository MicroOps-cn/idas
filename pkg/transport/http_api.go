package transport

import (
	"context"
	stdlog "log"

	"github.com/emicklei/go-restful/v3"
	kitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

func NewWebService(rootPath string, gv schema.GroupVersion, doc string) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(rootPath + "/" + gv.String()).
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

func InstallHTTPApi(logger log.Logger, container *restful.Container, options []httptransport.ServerOption, endpoints endpoint.Set) {
	container.Filter(HTTPLogging)
	restful.TraceLogger(stdlog.New(log.NewStdlibAdapter(level.Info(logger)), "[restful]", stdlog.LstdFlags|stdlog.Lshortfile))
	container.Filter(HTTPLoginAuthentication(endpoints))
	v1Ws := NewSimpleWebService("/api/v1", "基础接口")
	v1Ws.Route(v1Ws.POST("/login").Doc("用户登陆").To(NewKitHTTPServer[endpoint.UserLoginRequest](endpoints.UserLogin, options)).Metadata(global.MetaNeedLogin, false))
	v1Ws.Route(v1Ws.POST("/logout").Doc("用户退出登录").To(NewKitHTTPServer[endpoint.UserLogoutRequest](endpoints.UserLogout, options)))
	v1Ws.Route(v1Ws.GET("/user").Doc("获取当前登陆用户信息").To(NewKitHTTPServer[endpoint.CurrentUserRequest](endpoints.CurrentUser, options)))
	v1Ws.Route(v1Ws.GET("/resetPassword").Doc("重置用户密码").To(NewKitHTTPServer[endpoint.CurrentUserRequest](endpoints.CurrentUser, options)))

	container.Add(v1Ws)
	managerWs := NewWebService("/api", schema.GroupVersion{Group: "manager", Version: "v1"}, "管理接口")

	// 通用接口
	managerWs.Route(managerWs.POST("/file").Consumes("multipart/form-data").Doc("上传文件").To(NewKitHTTPServer[endpoint.FileUploadRequest](endpoints.UploadFile, options)))
	managerWs.Route(managerWs.GET("/file/{id}").Doc("下载/查看文件").To(NewKitHTTPServer[endpoint.FileDownloadRequest](endpoints.DownloadFile, options)))

	// 用户管理接口
	managerWs.Route(managerWs.GET("/users").Doc("获取用户列表").To(NewKitHTTPServer[endpoint.GetUsersRequest](endpoints.GetUsers, options)))
	managerWs.Route(managerWs.PATCH("/users").Doc("批量更新用户信息（增量）").To(NewKitHTTPServer[endpoint.PatchUsersRequest](endpoints.PatchUsers, options)))
	managerWs.Route(managerWs.DELETE("/users").Doc("批量删除用户").To(NewKitHTTPServer[endpoint.DeleteUsersRequest](endpoints.DeleteUsers, options)))
	managerWs.Route(managerWs.GET("/users/source").Doc("获取用户存储源").To(NewKitHTTPServer[endpoint.GetUserSourceRequest](endpoints.GetUserSource, options)))
	managerWs.Route(managerWs.GET("/user/{id}").Doc("获取用户信息").To(NewKitHTTPServer[endpoint.GetUserRequest](endpoints.GetUserInfo, options)))
	managerWs.Route(managerWs.POST("/user").Doc("创建用户").To(NewKitHTTPServer[endpoint.CreateUserRequest](endpoints.CreateUser, options)))
	managerWs.Route(managerWs.PUT("/user/{id}").Doc("更新用户信息（全量）").To(NewKitHTTPServer[endpoint.UpdateUserRequest](endpoints.UpdateUser, options)))
	managerWs.Route(managerWs.PATCH("/user/{id}").Doc("更新用户信息（增量）").To(NewKitHTTPServer[endpoint.PatchUserRequest](endpoints.PatchUser, options)))
	managerWs.Route(managerWs.DELETE("/user/{id}").Doc("删除用户").To(NewKitHTTPServer[endpoint.DeleteUserRequest](endpoints.DeleteUser, options)))

	// 会话管理接口
	managerWs.Route(managerWs.GET("/sessions").Doc("获取会话列表").To(NewKitHTTPServer[endpoint.GetSessionsRequest](endpoints.GetSessions, options)))
	managerWs.Route(managerWs.DELETE("/session/{id}").Doc("会话过期").To(NewKitHTTPServer[endpoint.DeleteSessionRequest](endpoints.DeleteSession, options)))

	// 应用管理接口
	managerWs.Route(managerWs.GET("/apps").Doc("获取应用列表").To(NewKitHTTPServer[endpoint.GetAppsRequest](endpoints.GetApps, options)))
	managerWs.Route(managerWs.PATCH("/apps").Doc("批量更新应用信息（增量）").To(NewKitHTTPServer[endpoint.PatchAppsRequest](endpoints.PatchApps, options)))
	managerWs.Route(managerWs.DELETE("/apps").Doc("批量删除应用").To(NewKitHTTPServer[endpoint.DeleteAppsRequest](endpoints.DeleteApps, options)))
	managerWs.Route(managerWs.GET("/apps/source").Doc("获取应用存储源").To(NewKitHTTPServer[endpoint.GetAppSourceRequest](endpoints.GetAppSource, options)))
	managerWs.Route(managerWs.GET("/app/{id}").Doc("获取应用信息").To(NewKitHTTPServer[endpoint.GetAppRequest](endpoints.GetAppInfo, options)))
	managerWs.Route(managerWs.POST("/app").Doc("创建应用").To(NewKitHTTPServer[endpoint.CreateAppRequest](endpoints.CreateApp, options)))
	managerWs.Route(managerWs.PUT("/app/{id}").Doc("更新应用信息（全量）").To(NewKitHTTPServer[endpoint.UpdateAppRequest](endpoints.UpdateApp, options)))
	managerWs.Route(managerWs.PATCH("/app/{id}").Doc("更新应用信息（增量）").To(NewKitHTTPServer[endpoint.PatchAppRequest](endpoints.PatchApp, options)))
	managerWs.Route(managerWs.DELETE("/app/{id}").Doc("删除应用").To(NewKitHTTPServer[endpoint.DeleteAppRequest](endpoints.DeleteApp, options)))
	container.Add(managerWs)

	oauthWs := NewWebService("/api", schema.GroupVersion{Group: "oauth", Version: "v1"}, "OAuth2")
	// https://www.ruanyifeng.com/blog/2019/04/oauth-grant-types.html
	oauthWs.Route(oauthWs.POST("/token").Doc("获取令牌").To(NewKitHTTPServer[endpoint.OAuthTokenRequest](endpoints.OAuthTokens, options)).Metadata(global.MetaNeedLogin, false).Consumes("application/x-www-form-urlencoded"))
	oauthWs.Route(oauthWs.POST("/authorize").Doc("应用授权").To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).Metadata(global.MetaAutoRedirectToLoginPage, true))
	oauthWs.Route(oauthWs.GET("/authorize").Doc("应用授权").To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).Metadata(global.MetaAutoRedirectToLoginPage, true))
	container.Add(oauthWs)
}
