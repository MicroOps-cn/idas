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

package transport

import (
	"context"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/global"
)

var apiServiceSet = []func(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService){
	UserService,
	AppService,
	FileService,
	SessionService,
	OAuthService,
	CurrentUserService,
	PermissionService,
	RoleService,
	PageService,
	ConfigService,
}

// UserService User Manager Service for restful Http container
func UserService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "users", Description: "Managing users"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetUsersRequest](ctx, endpoints.GetUsers, options)).
		Operation("getUsers").
		Doc("Get user list.").
		Params(StructToQueryParams(endpoint.GetUsersRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUsersResponse{}),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchUsersRequest](ctx, endpoints.PatchUsers, options)).
		Operation("patchUsers").
		Reads(endpoint.PatchUsersRequest{}).
		Doc("Batch update user information(Incremental).").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteUsersRequest](ctx, endpoints.DeleteUsers, options)).
		Operation("deleteUsers").
		Doc("Delete users in batch.").
		Reads(endpoint.DeleteUsersRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateUserRequest](ctx, endpoints.CreateUser, options)).
		Operation("createUser").
		Doc("Create a user.").
		Reads(endpoint.CreateUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetUserRequest](ctx, endpoints.GetUserInfo, options)).
		Operation("getUserInfo").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Get user information.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUserRequest{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateUserRequest](ctx, endpoints.UpdateUser, options)).
		Operation("updateUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Update user information(full).").
		Reads(endpoint.UpdateUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.UpdateUserRequest{}),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchUserRequest](ctx, endpoints.PatchUser, options)).
		Operation("patchUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Update user information(Incremental).").
		Reads(endpoint.PatchUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.PatchUserResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteUserRequest](ctx, endpoints.DeleteUser, options)).
		Operation("deleteUser").
		Param(v1ws.PathParameter("id", "identifier of the user").DataType("string")).
		Doc("Delete user.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/{userId}/key").
		To(NewKitHTTPServer[endpoint.CreateUserKeyRequest](ctx, endpoints.CreateUserKey, options)).
		Operation("createUserKey").
		Doc("Create a user key pair.").
		Param(v1ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Reads(endpoint.CreateUserKeyRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateUserKeyResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{userId}/key/{id}").
		To(NewKitHTTPServer[endpoint.DeleteUserKeyRequest](ctx, endpoints.DeleteUserKey, options)).
		Operation("deleteUserKey").
		Doc("Delete a user key pair.").
		Param(v1ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Param(v1ws.PathParameter("id", "identifier of the user key-pair").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/{userId}/key").
		To(NewKitHTTPServer[endpoint.GetUserKeysRequest](ctx, endpoints.GetUserKeys, options)).
		Operation("getUserKeys").
		Doc("Get a user key-pairs.").
		Param(v1ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetUserKeysResponse{}),
	)

	v1ws.Route(v1ws.POST("/key").
		To(NewKitHTTPServer[endpoint.CreateKeyRequest](ctx, endpoints.CreateKey, options)).
		Operation("createKey").
		Doc("Create your own key pair.").
		Reads(endpoint.CreateKeyRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateKeyResponse{}),
	)

	v1ws.Route(v1ws.POST("/sendActivateMail").
		To(NewKitHTTPServer[endpoint.SendActivationMailRequest](ctx, endpoints.SendActivateMail, options)).
		Operation("sendActivateMail").
		Reads(endpoint.SendActivationMailRequest{}).
		Doc("Send account activation email.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func AppService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "apps", Description: "Application manager"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetAppsRequest](ctx, endpoints.GetApps, options)).
		Operation("getApps").
		Doc("Get the application list.").
		Params(StructToQueryParams(endpoint.GetAppsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppsResponse{}),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchAppsRequest](ctx, endpoints.PatchApps, options)).
		Operation("patchApps").
		Doc("Batch update of application information (incremental).").
		Reads(endpoint.PatchAppsRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteAppsRequest](ctx, endpoints.DeleteApps, options)).
		Operation("deleteApps").
		Doc("Batch delete applications.").
		Reads(endpoint.DeleteAppsRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateAppRequest](ctx, endpoints.CreateApp, options)).
		Operation("createApp").
		Doc("Create an application.").
		Reads(endpoint.CreateAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetAppRequest](ctx, endpoints.GetAppInfo, options)).
		Operation("getAppInfo").
		Doc("Get Application info.").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppResponse{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateAppRequest](ctx, endpoints.UpdateApp, options)).
		Operation("updateApp").
		Doc("更新应用信息（全量）").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Reads(endpoint.UpdateAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.PATCH("/{id}").
		To(NewKitHTTPServer[endpoint.PatchAppRequest](ctx, endpoints.PatchApp, options)).
		Operation("patchApp").
		Doc("Update application information (full).").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Reads(endpoint.PatchAppRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteAppRequest](ctx, endpoints.DeleteApp, options)).
		Operation("deleteApp").
		Doc("Delete app.").
		Param(v1ws.PathParameter("id", "identifier of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/{appId}/key").
		To(NewKitHTTPServer[endpoint.CreateAppKeyRequest](ctx, endpoints.CreateAppKey, options)).
		Operation("createAppKey").
		Doc("Create a app key pair.").
		Param(v1ws.PathParameter("appId", "identifier of the app").DataType("string")).
		Reads(endpoint.CreateAppKeyRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.CreateAppKeyResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{appId}/key").
		To(NewKitHTTPServer[endpoint.DeleteAppKeysRequest](ctx, endpoints.DeleteAppKey, options)).
		Operation("deleteAppKeys").
		Doc("Delete a app key pairs.").
		Reads(endpoint.DeleteAppKeysRequest{}).
		Param(v1ws.PathParameter("appId", "identifier of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/{appId}/key").
		To(NewKitHTTPServer[endpoint.GetAppKeysRequest](ctx, endpoints.GetAppKeys, options)).
		Operation("getAppKeys").
		Doc("Get a app key-pairs.").
		Param(v1ws.PathParameter("appId", "identifier of the app").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetAppKeysResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func FileService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "files", Description: "Managing files"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[struct{}](ctx, endpoints.UploadFile, options)).
		Operation("uploadFile").
		Consumes("multipart/form-data").
		Doc("Upload file").
		Param(v1ws.MultiPartFormParameter("files", "files").AllowMultiple(true).DataType("file")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.FileUploadResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.FileDownloadRequest](ctx, endpoints.DownloadFile, options)).
		Operation("downloadFile").
		Param(v1ws.PathParameter("id", "identifier of the file").DataType("string").Required(true)).
		Doc("Download/View File").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func PageService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "pages", Description: "Managing pages"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetPagesRequest](ctx, endpoints.GetPages, options)).
		Operation("getPages").
		Doc("Get page list").
		Params(StructToQueryParams(endpoint.GetPagesRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(global.MetaForceOk, true).
		Returns(200, "OK", endpoint.GetPagesResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreatePageRequest](ctx, endpoints.CreatePage, options)).
		Operation("createPage").
		Doc("Create page.").
		Reads(endpoint.CreatePageRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchPagesRequest](ctx, endpoints.PatchPages, options)).
		Operation("patchPages").
		Reads(endpoint.PatchPagesRequest{}).
		Doc("Batch patch page config(Incremental).").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.GET("/{id}").
		To(NewKitHTTPServer[endpoint.GetPageRequest](ctx, endpoints.GetPage, options)).
		Operation("getPage").
		Doc("Get a page configs.").
		Param(v1ws.PathParameter("id", "identifier of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetPageResponse{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdatePageRequest](ctx, endpoints.UpdatePage, options)).
		Operation("updatePage").
		Doc("Update page (full).").
		Param(v1ws.PathParameter("id", "identifier of the page").DataType("string")).
		Reads(endpoint.UpdatePageRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeletePageRequest](ctx, endpoints.DeletePage, options)).
		Operation("deletePage").
		Doc("Delete a page.").
		Param(v1ws.PathParameter("id", "identifier of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/{pageId}/data").
		To(NewKitHTTPServer[endpoint.GetPageDatasRequest](ctx, endpoints.GetPageDatas, options)).
		Operation("getPageDatas").
		Doc("Get data list of page").
		Params(StructToQueryParams(endpoint.GetPageDatasRequest{})...).
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetPageDatasResponse{}),
	)
	v1ws.Route(v1ws.POST("/{pageId}/data").
		To(NewKitHTTPServer[endpoint.CreatePageDataRequest](ctx, endpoints.CreatePageData, options)).
		Operation("createPageData").
		Doc("Create a data of a page.").
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Reads(endpoint.CreatePageDataRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	v1ws.Route(v1ws.PATCH("/{pageId}/data").
		To(NewKitHTTPServer[endpoint.PatchPageDatasRequest](ctx, endpoints.PatchPageDatas, options)).
		Operation("patchPageDatas").
		Reads(endpoint.PatchPageDatasRequest{}).
		Doc("Batch patch data of a page(Incremental).").
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.GET("/{pageId}/data/{id}").
		To(NewKitHTTPServer[endpoint.GetPageDataRequest](ctx, endpoints.GetPageData, options)).
		Operation("getPageData").
		Doc("Get the specified data of a page.").
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Param(v1ws.PathParameter("id", "data id of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetPageDataResponse{}),
	)
	v1ws.Route(v1ws.PUT("/{pageId}/data/{id}").
		To(NewKitHTTPServer[endpoint.UpdatePageDataRequest](ctx, endpoints.UpdatePageData, options)).
		Operation("updatePageData").
		Doc("Update data of a page. (full).").
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Param(v1ws.PathParameter("id", "data id of the page").DataType("string")).
		Reads(endpoint.UpdatePageDataRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{pageId}/data/{id}").
		To(NewKitHTTPServer[endpoint.DeletePageDataRequest](ctx, endpoints.DeletePageData, options)).
		Operation("deletePageData").
		Doc("Delete data of a page.").
		Param(v1ws.PathParameter("pageId", "identifier of the page").DataType("string")).
		Param(v1ws.PathParameter("id", "data id of the page").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func SessionService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "sessions", Description: "Managing sessions"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetSessionsRequest](ctx, endpoints.GetSessions, options)).
		Operation("getSessions").
		Doc("Get session list.").
		Params(StructToQueryParams(endpoint.GetSessionsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetSessionsResponse{}),
	)
	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteSessionRequest](ctx, endpoints.DeleteSession, options)).
		Operation("deleteSession").
		Param(v1ws.PathParameter("id", "identifier of the session").DataType("string").Required(true)).
		Doc("Expire a session.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func OAuthService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "oauth", Description: "OAuth2.0 Support"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)

	// https://www.ruanyifeng.com/blog/2019/04/oauth-grant-types.html
	v1ws.Route(v1ws.POST("/token").
		To(NewSimpleKitHTTPServer[endpoint.OAuthTokenRequest](ctx, endpoints.OAuthTokens, decodeHTTPRequest[endpoint.OAuthTokenRequest], simpleEncodeHTTPResponse, options)).
		Operation("oAuthTokens").
		Doc("Get token.").
		Filter(HTTPApplicationAuthenticationFilter(endpoints)).
		Reads(endpoint.OAuthTokenRequest{}).
		Consumes("application/x-www-form-urlencoded", restful.MIME_JSON).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.OAuthTokenResponse{}),
	)

	v1ws.Route(v1ws.GET("/userinfo").
		To(NewSimpleKitHTTPServer[endpoint.OAuthTokenRequest](ctx, endpoints.CurrentUser, decodeHTTPRequest[endpoint.OAuthTokenRequest], simpleEncodeHTTPResponse, options)).
		Operation("oAuthUserInfo").
		Doc("Get user info.").
		Filter(HTTPAuthenticationFilter(endpoints)).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](ctx, endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("Application authorization.").
		Filter(HTTPAuthenticationFilter(endpoints)).
		Reads(endpoint.OAuthAuthorizeRequest{}).
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("/authorize").
		To(NewKitHTTPServer[endpoint.OAuthAuthorizeRequest](ctx, endpoints.OAuthAuthorize, options)).
		Operation("oAuthAuthorize").
		Doc("Application authorization.").
		Filter(HTTPAuthenticationFilter(endpoints)).
		Params(StructToQueryParams(endpoint.OAuthAuthorizeRequest{})...).
		Metadata(global.MetaAutoRedirectToLoginPage, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func CurrentUserService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "user", Description: "Current user service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.POST("/login").
		To(NewKitHTTPServer[endpoint.UserLoginRequest](ctx, endpoints.UserLogin, options)).
		Operation("userLogin").
		Doc("User login.").
		Reads(endpoint.UserLoginRequest{}).
		Metadata(global.MetaNeedLogin, false).
		Metadata(global.MetaSensitiveData, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.UserLoginResponse{}),
	)
	v1ws.Route(v1ws.POST("/logout").
		To(NewKitHTTPServer[struct{}](ctx, endpoints.UserLogout, options)).
		Operation("userLogout").
		Doc("User logout.").
		Consumes("*/*").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[struct{}](ctx, endpoints.CurrentUser, options)).
		Operation("currentUser").
		Doc("Get current login user information.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(global.MetaUpdateLastSeen, true).
		Returns(200, "OK", endpoint.GetUserResponse{}),
	)

	v1ws.Route(v1ws.PUT("").
		To(NewKitHTTPServer[endpoint.UpdateUserRequest](ctx, endpoints.UpdateCurrentUser, options)).
		Operation("updateCurrentUser").
		Doc("Update current login user information (full).").
		Reads(endpoint.UpdateUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.PATCH("").
		To(NewKitHTTPServer[endpoint.PatchCurrentUserRequest](ctx, endpoints.PatchCurrentUser, options)).
		Operation("patchCurrentUser").
		Doc("Update current login user information (increment).").
		Reads(endpoint.PatchCurrentUserRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.GET("totp/secret").
		To(NewKitHTTPServer[endpoint.CreateTOTPSecretRequest](ctx, endpoints.CreateTOTPSecret, options)).
		Operation("getTOTPSecret").
		Doc("get TOTP Secret").
		Params(StructToQueryParams(endpoint.CreateTOTPSecretRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(global.MetaSensitiveData, true).
		Metadata(global.MetaNeedLogin, false).
		Returns(200, "OK", endpoint.CreateTOTPSecretResponse{}),
	)

	v1ws.Route(v1ws.POST("totp").
		To(NewKitHTTPServer[endpoint.CreateTOTPRequest](ctx, endpoints.CreateTOTP, options)).
		Operation("bindingTOTP").
		Doc("binding TOTP Secret").
		Reads(endpoint.CreateTOTPRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/activateAccount").
		To(NewKitHTTPServer[endpoint.ActivateAccountRequest](ctx, endpoints.ActivateAccount, options)).
		Operation("activateAccount").
		Reads(endpoint.ActivateAccountRequest{}).
		Doc("Activate the user.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(global.MetaNeedLogin, false).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/forgotPassword").
		To(NewKitHTTPServer[endpoint.ForgotUserPasswordRequest](ctx, endpoints.ForgotPassword, options)).
		Operation("forgotPassword").
		Doc("Forgot the user password.").
		Reads(endpoint.ForgotUserPasswordRequest{}).
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/resetPassword").
		To(NewKitHTTPServer[endpoint.ResetUserPasswordRequest](ctx, endpoints.ResetPassword, options)).
		Operation("resetPassword").
		Reads(endpoint.ResetUserPasswordRequest{}).
		Doc("Reset the user password.").
		Metadata(global.MetaNeedLogin, false).
		Metadata(global.MetaSensitiveData, true).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.POST("/sendLoginCaptcha").
		To(NewKitHTTPServer[endpoint.SendLoginCaptchaRequest](ctx, endpoints.SendLoginCaptcha, options)).
		Operation("sendLoginCaptcha").
		Reads(endpoint.SendLoginCaptchaRequest{}).
		Doc("Send login code.").
		Metadata(global.MetaNeedLogin, false).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.SendLoginCaptchaResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func PermissionService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "permissions", Description: "permissions service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetPermissionsRequest](ctx, endpoints.GetPermissions, options)).
		Operation("getPermissions").
		Doc("Get permission list.").
		Params(StructToQueryParams(endpoint.GetPermissionsRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetPermissionsResponse{}),
	)
	return tag, []*restful.WebService{v1ws}
}

func RoleService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "roles", Description: "role service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("").
		To(NewKitHTTPServer[endpoint.GetRolesRequest](ctx, endpoints.GetRoles, options)).
		Operation("getRoles").
		Doc("Get role list.").
		Params(StructToQueryParams(endpoint.GetRolesRequest{})...).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetRolesResponse{}),
	)
	v1ws.Route(v1ws.DELETE("").
		To(NewKitHTTPServer[endpoint.DeleteRolesRequest](ctx, endpoints.DeleteRoles, options)).
		Operation("deleteRoles").
		Doc("Batch delete roles.").
		Reads(endpoint.DeleteRolesRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseTotalResponse{}),
	)
	v1ws.Route(v1ws.POST("").
		To(NewKitHTTPServer[endpoint.CreateRoleRequest](ctx, endpoints.CreateRole, options)).
		Operation("createRole").
		Doc("Create role.").
		Reads(endpoint.CreateRoleRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)
	v1ws.Route(v1ws.PUT("/{id}").
		To(NewKitHTTPServer[endpoint.UpdateRoleRequest](ctx, endpoints.UpdateRole, options)).
		Operation("updateRole").
		Doc("Update role information (full).").
		Param(v1ws.PathParameter("id", "identifier of the role").DataType("string")).
		Reads(endpoint.UpdateRoleRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	v1ws.Route(v1ws.DELETE("/{id}").
		To(NewKitHTTPServer[endpoint.DeleteRoleRequest](ctx, endpoints.DeleteRole, options)).
		Operation("deleteRole").
		Doc("删除角色").
		Param(v1ws.PathParameter("id", "identifier of the role").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}

func ConfigService(ctx context.Context, options []httptransport.ServerOption, endpoints endpoint.Set) (spec.Tag, []*restful.WebService) {
	tag := spec.Tag{TagProps: spec.TagProps{Name: "config", Description: "config service"}}
	tags := []string{tag.Name}
	v1ws := NewWebService(rootPath, schema.GroupVersion{Group: tag.Name, Version: "v1"}, tag.Description)
	v1ws.Filter(HTTPAuthenticationFilter(endpoints))

	v1ws.Route(v1ws.GET("security").
		To(NewKitHTTPServer[struct{}](ctx, endpoints.GetSecurityConfig, options)).
		Operation("getSecurityConfig").
		Doc("Obtain Security Configuration.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.GetSecurityConfigResponse{}),
	)
	v1ws.Route(v1ws.PATCH("security").
		To(NewKitHTTPServer[endpoint.PatchSecurityConfigRequest](ctx, endpoints.PatchSecurityConfig, options)).
		Operation("patchSecurityConfig").
		Doc("Update Security Configuration (Incremental).").
		Reads(endpoint.PatchSecurityConfigRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", endpoint.BaseResponse{}),
	)

	return tag, []*restful.WebService{v1ws}
}
