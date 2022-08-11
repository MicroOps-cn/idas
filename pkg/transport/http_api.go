package transport

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"idas/pkg/endpoint"
	"idas/pkg/global"
)

var apiServiceSet = []func(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService){
	UserService,
	AppService,
	FileService,
	SessionService,
	OAuthService,
	UserAuthService,
	PermissionService,
	RoleService,
}

// UserService User Manager Service for restful Http container
func UserService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "users", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetUsersRequest](endpoints.GetUsers, options)).
		Operation("getUsers").
		Doc("Get user list.").
		Params(StructToQueryParams(endpoint.GetUsersRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUsersResponse{}),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchUsersRequest](endpoints.PatchUsers, options)).
		Operation("patchUsers").
		Reads(endpoint.PatchUsersRequest{}).
		Doc("Batch update user information(Incremental).").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.PatchUsersResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteUsersRequest](endpoints.DeleteUsers, options)).
		Operation("deleteUsers").
		Doc("Delete users in batch.").
		Reads(endpoint.DeleteUsersRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateUserRequest](endpoints.CreateUser, options)).
		Operation("createUser").
		Doc("Create a user.").
		Reads(endpoint.CreateUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateUserResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetUserRequest](endpoints.GetUserInfo, options)).
		Operation("getUserInfo").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Param(v1ws.QueryParameter("storage", "storage of the user").DataType("string")).
		Doc("Get user information.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUserRequest{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateUserRequest](endpoints.UpdateUser, options)).
		Operation("updateUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Update user information(full).").
		Reads(endpoint.UpdateUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.UpdateUserRequest{}),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchUserRequest](endpoints.PatchUser, options)).
		Operation("patchUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Update user information(Incremental).").
		Reads(endpoint.PatchUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.PatchUserResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteUserRequest](endpoints.DeleteUser, options)).
		Operation("deleteUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Param(v1ws.QueryParameter("storage", "storage source of the user").DataType("string")).
		Doc("Delete user.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.POST("/forgotPassword").
		To(NewKitHTTPServer[endpoint.ForgotUserPasswordRequest](endpoints.ForgotPassword, options)).
		Operation("forgotPassword").
		Doc("Forgot the user password.").
		Reads(endpoint.ForgotUserPasswordRequest{}).
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/resetPassword").
		To(NewKitHTTPServer[endpoint.ResetUserPasswordRequest](endpoints.ResetPassword, options)).
		Operation("resetPassword").
		Reads(endpoint.ResetUserPasswordRequest{}).
		Doc("Reset the user password.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.GET("/source").
		To(NewKitHTTPServer[struct{}](endpoints.GetUserSource, options)).
		Operation("getUserSource").
		Doc("Get the user storage source.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUserSourceResponse{}),
	)

	v1ws.Route(v1ws.POST("/{userId}/key").
		To(NewKitHTTPServer[endpoint.CreateUserKeyRequest](endpoints.CreateUserKey, options)).
		Operation("createUserKey").
		Doc("Create a user key pair.").
		Param(v1ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Reads(endpoint.CreateUserKeyRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateUserKeyResponse{}),
	)
	v1ws.Route(v1ws.POST("/key").
		To(NewKitHTTPServer[endpoint.CreateKeyRequest](endpoints.CreateKey, options)).
		Operation("createKey").
		Doc("Create your own key pair.").
		Reads(endpoint.CreateKeyRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateKeyResponse{}),
	)

	v1ws.Route(v1ws.POST("/sendActivateMail").
		To(NewKitHTTPServer[endpoint.SendActivationMailRequest](endpoints.SendActivateMail, options)).
		Operation("sendActivateMail").
		Reads(endpoint.SendActivationMailRequest{}).
		Doc("Send account activation email.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/activateAccount").
		To(NewKitHTTPServer[endpoint.ActivateAccountRequest](endpoints.ActivateAccount, options)).
		Operation("activateAccount").
		Reads(endpoint.ActivateAccountRequest{}).
		Doc("Activate the user.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(global.MetaNeedLogin, false).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func AppService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "apps", Description: "Application manager"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetAppsRequest](endpoints.GetApps, options)).
		Operation("getApps").
		Doc("Get the application list.").
		Params(StructToQueryParams(endpoint.GetAppsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppsResponse{}),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchAppsRequest](endpoints.PatchApps, options)).
		Operation("patchApps").
		Doc("批量更新应用信息（增量）").
		Reads(endpoint.PatchAppsRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.PatchAppsResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteAppsRequest](endpoints.DeleteApps, options)).
		Operation("deleteApps").
		Doc("批量删除应用").
		Reads(endpoint.DeleteAppsRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateAppRequest](endpoints.CreateApp, options)).
		Operation("createApp").
		Doc("创建应用").
		Reads(endpoint.CreateAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/source").
		To(NewKitHTTPServer[struct{}](endpoints.GetAppSource, options)).
		Operation("getAppSource").
		Doc("获取应用存储源").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppSourceResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetAppRequest](endpoints.GetAppInfo, options)).
		Operation("getAppInfo").
		Doc("获取应用信息").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Param(v1ws.QueryParameter("storage", "storage of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppResponse{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateAppRequest](endpoints.UpdateApp, options)).
		Operation("updateApp").
		Doc("更新应用信息（全量）").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Reads(endpoint.UpdateAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchAppRequest](endpoints.PatchApp, options)).
		Operation("patchApp").
		Doc("更新应用信息（增量）").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Reads(endpoint.PatchAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteAppRequest](endpoints.DeleteApp, options)).
		Operation("deleteApp").
		Doc("删除应用").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Param(v1ws.QueryParameter("storage", "storage source of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func FileService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "files", Description: "Managing files"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[struct{}](endpoints.UploadFile, options)).
		Operation("uploadFile").
		Consumes("multipart/form-data").
		Doc("Upload file").
		Param(v1ws.MultiPartFormParameter("files", "files").AllowMultiple(true).DataType("file")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.FileUploadResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.FileDownloadRequest](endpoints.DownloadFile, options)).
		Operation("downloadFile").
		Param(v1ws.PathParameter("id", "identifier of the file").DataType("string").Required(true)).
		Doc("Download/View File").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func SessionService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "sessions", Description: "Managing sessions"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetSessionsRequest](endpoints.GetSessions, options)).
		Operation("getSessions").
		Doc("获取会话列表").
		Params(StructToQueryParams(endpoint.GetSessionsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetSessionsResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteSessionRequest](endpoints.DeleteSession, options)).
		Operation("deleteSession").
		Param(v1ws.PathParameter("id", "identifier of the session").DataType("string").Required(true)).
		Doc("会话过期").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func OAuthService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "oauth", Description: "OAuth2.0 Support"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	// https://www.ruanyifeng.com/blog/2019/04/oauth-grant-types.html
	v1ws.Route(v1ws.POST("/token").
		To(NewSimpleKitHTTPServer[endpoint.OAuthTokenRequest](endpoints.OAuthTokens, decodeHTTPRequest[endpoint.OAuthTokenRequest], simpleEncodeHTTPResponse, options)).
		Operation("oAuthTokens").
		Doc("获取令牌").
		Metadata(global.MetaNeedLogin, false).
		Reads(endpoint.OAuthTokenRequest{}).
		Consumes("application/x-www-form-urlencoded", restful.MIME_JSON).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.OAuthTokenResponse{}),
	)
	v1ws.Route(v1ws.POST("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("应用授权").
		Reads(endpoint.OAuthAuthorizeRequest{}).
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("应用授权").
		Params(StructToQueryParams(endpoint.OAuthAuthorizeRequest{})...).
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/userinfo").
		To(NewSimpleKitHTTPServer[endpoint.OAuthTokenRequest](endpoints.CurrentUser, decodeHTTPRequest[endpoint.OAuthTokenRequest], simpleEncodeHTTPResponse, options)).
		Operation("oAuthUserInfo").
		Doc("获取用户信息").
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func UserAuthService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "user", Description: "user login service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.POST("/login").
		To(NewKitHTTPServer[endpoint.UserLoginRequest](endpoints.UserLogin, options)).
		Operation("userLogin").
		Doc("用户登陆").
		Reads(endpoint.UserLoginRequest{}).
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.POST("/logout").
		To(NewKitHTTPServer[struct{}](endpoints.UserLogout, options)).
		Operation("userLogout").
		Doc("用户退出登录").
		Consumes("*/*").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[struct{}](endpoints.CurrentUser, options)).
		Operation("currentUser").
		Doc("获取当前登陆用户信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUserResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func PermissionService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "permissions", Description: "permissions service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetPermissionsRequest](endpoints.GetPermissions, options)).
		Operation("getPermissions").
		Doc("获取权限列表").
		Params(StructToQueryParams(endpoint.GetPermissionsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetPermissionsResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func RoleService(options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "roles", Description: "role service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetRolesRequest](endpoints.GetRoles, options)).
		Operation("getRoles").
		Doc("获取角色列表").
		Params(StructToQueryParams(endpoint.GetRolesRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetRolesResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteRolesRequest](endpoints.DeleteRoles, options)).
		Operation("deleteRoles").
		Doc("批量删除角色").
		Reads(endpoint.DeleteRolesRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateRoleRequest](endpoints.CreateRole, options)).
		Operation("createRole").
		Doc("创建角色").
		Reads(endpoint.CreateRoleRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateRoleRequest](endpoints.UpdateRole, options)).
		Operation("updateRole").
		Doc("更新角色信息（全量）").
		Param(v1ws.PathParameter("id", "identifier of the role").DataType("string")).
		Reads(endpoint.UpdateRoleRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteRoleRequest](endpoints.DeleteRole, options)).
		Operation("deleteRole").
		Doc("删除角色").
		Param(v1ws.PathParameter("id", "identifier of the role").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}
