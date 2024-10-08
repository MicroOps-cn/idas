declare namespace API {
  
        type ActivateAccountRequest =
          {
                'newPassword': string;
                'token'?: string;
                'userId': string;
          }
        
  
        type App =
          {
                'avatar': string;
                'createTime'?: string;
                'description': string;
                'displayName'?: string;
                'grantMode': AppMetaGrantMode;
                'grantType': AppMetaGrantType;
                'i18n': AppI18NOptions;
                'id': string;
                'isDelete'?: boolean;
                'name': string;
                'oAuth2': AppOAuth2;
                'proxy'?: AppProxy;
                'role'?: string;
                'roleId'?: string;
                'roles'?: AppRole[];
                'status': AppMetaStatus;
                'updateTime'?: string;
                'url': string;
                'users'?: User[];
          }
        
  
        type AppI18NOptions =
          {
                'description'?: Record<string, any>;
                'displayName'?: Record<string, any>;
          }
        
  
        type AppInfo =
          {
                'avatar'?: string;
                'createTime': string;
                'description'?: string;
                'displayName'?: string;
                'grantMode': AppMetaGrantMode;
                'grantType': AppMetaGrantType[];
                'i18n'?: AppI18NOptions;
                'id': string;
                'isDelete': boolean;
                'name': string;
                'oAuth2'?: AppOAuth2Info;
                'proxy'?: AppProxyInfo;
                'roles'?: AppRoleInfo[];
                'status': AppMetaStatus;
                'updateTime': string;
                'url': string;
                'users'?: UserInfo[];
          }
        
  
        type AppKeyInfo =
          {
                'appId': string;
                'createTime': string;
                'id': string;
                'key': string;
                'name': string;
                'privateKey': string;
                'secret': string;
                'updateTime': string;
          }
        
  
        type AppMetaGrantMode = "manual" | 0 | "full" | 1;
  
        type AppMetaGrantType = "none" | 0 | "authorization_code" | 1 | "implicit" | 2 | "password" | 4 | "client_credentials" | 8 | "proxy" | 16 | "oidc" | 32 | "radius" | 64;
  
        type AppMetaJWTSignatureMethod = "default" | 0 | "HS256" | 1 | "HS384" | 2 | "HS512" | 3 | "RS256" | 4 | "RS384" | 5 | "RS512" | 6 | "ES256" | 7 | "ES384" | 8 | "ES512" | 9;
  
        type AppMetaStatus = "unknown" | 0 | "normal" | 1 | "disable" | 2;
  
        type AppOAuth2 =
          {
                'authorizedRedirectUrl': string[];
                'jwtSignatureKey'?: string;
                'jwtSignatureMethod': AppMetaJWTSignatureMethod;
          }
        
  
        type AppOAuth2Info =
          {
                'authorizedRedirectUrl': string[];
                'jwtSignatureKey'?: string;
                'jwtSignatureMethod': AppMetaJWTSignatureMethod;
          }
        
  
        type AppProxy =
          {
                'appId': string;
                'createTime'?: string;
                'domain': string;
                'hstsOffload': boolean;
                'id': string;
                'insecureSkipVerify': boolean;
                'isDelete'?: boolean;
                'jwtCookieName': string;
                'jwtProvider': boolean;
                'jwt_secret': string;
                'jwt_secret_salt': string;
                'transparentServerName': boolean;
                'updateTime'?: string;
                'upstream': string;
                'urls': AppProxyUrl[];
          }
        
  
        type AppProxyInfo =
          {
                'domain': string;
                'hstsOffload': boolean;
                'insecureSkipVerify': boolean;
                'jwtCookieName': string;
                'jwtProvider': boolean;
                'jwtSecret': string;
                'transparentServerName': boolean;
                'upstream': string;
                'urls': AppProxyUrlInfo[];
          }
        
  
        type AppProxyUrl =
          {
                'createTime'?: string;
                'id': string;
                'index'?: number;
                'isDelete'?: boolean;
                'method': string;
                'name': string;
                'updateTime'?: string;
                'upstream'?: string;
                'url': string;
          }
        
  
        type AppProxyUrlInfo =
          {
                'id': string;
                'method': string;
                'name': string;
                'upstream'?: string;
                'url': string;
          }
        
  
        type AppRole =
          {
                'appId': string;
                'createTime'?: string;
                'id': string;
                'isDefault'?: boolean;
                'isDelete'?: boolean;
                'name': string;
                'updateTime'?: string;
                'urls': string[];
                'users': User[];
          }
        
  
        type AppRoleInfo =
          {
                'id': string;
                'isDefault'?: boolean;
                'name': string;
                'urls'?: string[];
          }
        
  
        type AppUser =
          {
                'id': string;
                'role'?: string;
                'roleId'?: string;
          }
        
  
        type BaseListResponse =
          {
                'current': number;
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type BaseResponse =
          {
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type BaseTotalResponse =
          {
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type createAppKeyParams =
          {
              /** identifier of the app */
                'appId': string;
          }
        
  
        type CreateAppKeyRequest =
          {
                'appId': string;
                'name': string;
          }
        
  
        type CreateAppKeyResponse =
          {
                'data'?: AppKeyInfo;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type CreateAppRequest =
          {
                'avatar'?: string;
                'description'?: string;
                'displayName'?: string;
                'grantMode'?: AppMetaGrantMode;
                'grantType'?: AppMetaGrantType[];
                'i18n'?: AppI18NOptions;
                'name': string;
                'oAuth2': AppOAuth2Info;
                'proxy'?: AppProxyInfo;
                'roles'?: AppRoleInfo[];
                'url': string;
                'users'?: AppUser[];
          }
        
  
        type createPageDataParams =
          {
              /** identifier of the page */
                'pageId': string;
          }
        
  
        type CreatePageDataRequest =
          {
                'data'?: Record<string, any>;
                'pageId': string;
          }
        
  
        type CreatePageRequest =
          {
                'description'?: string;
                'fields'?: FieldConfig[];
                'icon'?: string;
                'name': string;
          }
        
  
        type CreateRoleRequest =
          {
                'description'?: string;
                'name': string;
                'permission'?: string[];
          }
        
  
        type CreateTOTPRequest =
          {
                'firstCode': string;
                'secondCode': string;
                'token': string;
          }
        
  
        type CreateTOTPSecretResponse =
          {
                'data'?: CreateTOTPSecretResponseData;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type CreateTOTPSecretResponseData =
          {
                'secret': string;
                'token': string;
          }
        
  
        type CreateUserRequest =
          {
                'apps'?: UserApp[];
                'avatar'?: string;
                'email'?: string;
                'fullName'?: string;
                'isDelete'?: boolean;
                'phoneNumber'?: string;
                'status'?: UserMetaUserStatus;
                'username': string;
          }
        
  
        type currentUserAppsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
          }
        
  
        type currentUserEventLogsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'eventId': string;
          }
        
  
        type currentUserEventsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'action'?: string;
                'startTime': string;
                'endTime': string;
          }
        
  
        type deleteAppKeysParams =
          {
              /** identifier of the app */
                'appId': string;
          }
        
  
        type DeleteAppKeysRequest =
          {
                'appId': string;
                'id': string;
          }
        
  
        type deleteAppParams =
          {
              /** identifier of the app */
                'id': string;
          }
        
  
        type DeleteAppRequest =
          {
                'id': string;
          }
        
  
        type deleteCurrentUserSessionParams =
          {
              /** identifier of the session */
                'id': string;
          }
        
  
        type deletePageDataParams =
          {
              /** identifier of the page */
                'pageId': string;
              /** data id of the page */
                'id': string;
          }
        
  
        type deletePageParams =
          {
              /** identifier of the page */
                'id': string;
          }
        
  
        type deleteRoleParams =
          {
              /** identifier of the role */
                'id': string;
          }
        
  
        type DeleteRoleRequest =
          {
                'id': string;
          }
        
  
        type deleteSessionParams =
          {
              /** identifier of the session */
                'id': string;
          }
        
  
        type deleteUserParams =
          {
              /** identifier of the user */
                'id': string;
          }
        
  
        type DeleteUserRequest =
          {
                'id': string;
          }
        
  
        type downloadFileParams =
          {
              /** identifier of the file */
                'id': string;
          }
        
  
        type Event =
          {
                'action': string;
                'client_ip': string;
                'createTime': string;
                'id': string;
                'location': string;
                'message': string;
                'status': string;
                'took': number;
                'updateTime': string;
                'userId': string;
                'username': string;
          }
        
  
        type EventLog =
          {
                'createTime': string;
                'id': string;
                'log': string;
                'updateTime': string;
                'userId': string;
          }
        
  
        type FieldConfig =
          {
                'defaultValue'?: string;
                'displayName'?: string;
                'max'?: number;
                'maxWidth'?: number;
                'min'?: number;
                'minWidth'?: number;
                'name': string;
                'tooltip'?: string;
                'valueEnum'?: Record<string, any>;
                'valueType': FieldType;
          }
        
  
        type FieldType = "text" | 0 | "textarea" | 2 | "digit" | 3 | "digitRange" | 4 | "checkbox" | 5 | "radio" | 6 | "switch" | 7 | "select" | 8 | "multiSelect" | 9 | "timeRange" | 10 | "date" | 11 | "dateRange" | 12 | "dateTime" | 13 | "dateTimeRange" | 14;
  
        type FileUploadResponse =
          {
                'data': Record<string, any>;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type ForgotUserPasswordRequest =
          {
                'email': string;
                'username': string;
          }
        
  
        type getAppIconsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
          }
        
  
        type GetAppIconsResponse =
          {
                'current': number;
                'data'?: Model[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getAppInfoParams =
          {
              /** identifier of the app */
                'id': string;
          }
        
  
        type getAppKeysParams =
          {
              /** identifier of the app */
                'appId': string;
          }
        
  
        type GetAppKeysResponse =
          {
                'data'?: SimpleAppKeyInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type GetAppResponse =
          {
                'data'?: AppInfo;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type getAppsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
          }
        
  
        type GetAppsResponse =
          {
                'current': number;
                'data'?: AppInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type GetCurrentUserEventLogsResponse =
          {
                'current': number;
                'data'?: EventLog[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type GetCurrentUserEventsResponse =
          {
                'current': number;
                'data'?: Event[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getCurrentUserSessionsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'userId'?: string;
          }
        
  
        type getEventLogsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'eventId': string;
          }
        
  
        type GetEventLogsResponse =
          {
                'current': number;
                'data'?: EventLog[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getEventsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'username'?: string;
                'action'?: string;
                'startTime': string;
                'endTime': string;
          }
        
  
        type GetEventsResponse =
          {
                'current': number;
                'data'?: Event[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getPageDataParams =
          {
              /** identifier of the page */
                'pageId': string;
              /** data id of the page */
                'id': string;
          }
        
  
        type GetPageDataResponse =
          {
                'data'?: PageData;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type getPageDatasParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'pageId': string;
                'filters'?: any;
              /** identifier of the page */
                'pageId': string;
          }
        
  
        type GetPageDatasResponse =
          {
                'current': number;
                'data'?: PageData[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getPageParams =
          {
              /** identifier of the page */
                'id': string;
          }
        
  
        type GetPageResponse =
          {
                'data'?: PageConfig;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type getPagesParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'status'?: PageStatus;
          }
        
  
        type GetPagesResponse =
          {
                'current': number;
                'data'?: PageConfig[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getPermissionsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
          }
        
  
        type GetPermissionsResponse =
          {
                'current': number;
                'data'?: PermissionInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getRolesParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
          }
        
  
        type GetRolesResponse =
          {
                'current': number;
                'data'?: RoleInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type GetSecurityConfigResponse =
          {
                'data'?: RuntimeSecurityConfig;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type getSessionsParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'userId'?: string;
          }
        
  
        type GetSessionsResponse =
          {
                'current': number;
                'data'?: SessionInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type getTOTPSecretParams =
          {
                'token'?: any;
          }
        
  
        type getUserInfoParams =
          {
              /** identifier of the user */
                'id': string;
          }
        
  
        type GetUserResponse =
          {
                'data'?: UserInfo;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type getUsersParams =
          {
                'pageSize'?: number;
                'current'?: number;
                'keywords'?: string;
                'app'?: string;
                'status'?: UserMetaUserStatus;
          }
        
  
        type GetUsersResponse =
          {
                'current': number;
                'data'?: UserInfo[];
                'errorCode'?: string;
                'errorMessage'?: string;
                'pageSize': number;
                'success': boolean;
                'total': number;
                'traceId': string;
          }
        
  
        type GlobalConfig =
          {
                'copyright'?: string;
                'defaultLoginType': LoginType;
                'loginType': GlobalLoginType[];
                'logo'?: string;
                'subTitle'?: string;
                'title'?: string;
                'version'?: string;
          }
        
  
        type GlobalConfigResponse =
          {
                'data'?: GlobalConfig;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type GlobalLoginType =
          {
                'autoLogin'?: boolean;
                'autoRedirect'?: boolean;
                'icon'?: string;
                'id'?: string;
                'name'?: string;
                'type': LoginType;
          }
        
  
        type LoginType = "normal" | 0 | "mfa_totp" | 1 | "mfa_email" | 2 | "mfa_sms" | 3 | "email" | 4 | "sms" | 5 | "oauth2" | 6 | "enable_mfa_totp" | 10 | "enable_mfa_email" | 11 | "enable_mfa_sms" | 12;
  
        type Model =
          {
                'createTime'?: string;
                'id': string;
                'isDelete'?: boolean;
                'updateTime'?: string;
          }
        
  
        type oAuthAuthorizeParams =
          {
                'response_type'?: any;
                'client_id'?: string;
                'redirect_uri'?: string;
                'state'?: string;
                'scope'?: string;
                'access_type'?: string;
                'code_challenge'?: string;
                'code_challenge_method'?: OAuthAuthorizeRequestCodeChallengeMethod;
          }
        
  
        type OAuthAuthorizeRequest =
          {
                'access_type'?: string;
                'client_id'?: string;
                'code_challenge'?: string;
                'code_challenge_method'?: OAuthAuthorizeRequestCodeChallengeMethod;
                'redirect_uri'?: string;
                'response_type'?: OAuthAuthorizeRequestResponseTypes;
                'scope'?: string;
                'state'?: string;
          }
        
  
        type OAuthAuthorizeRequestCodeChallengeMethod = "plain" | 0 | "S256" | 1;
  
        type OAuthAuthorizeRequestResponseType = "none" | 0 | "code" | 1 | "token" | 2 | "id_token" | 3;
  
        type OAuthAuthorizeRequestResponseTypes =
          {
                'types'?: OAuthAuthorizeRequestResponseType[];
          }
        
  
        type OAuthGrantType = "refresh_token" | 0 | "authorization_code" | 1 | "password" | 2 | "client_credentials" | 3;
  
        type OAuthJWKSResponse =
          {
                'keys'?: OAuthJWKSResponseKey[];
          }
        
  
        type OAuthJWKSResponseKey =
          {
                'alg'?: string;
                'crv'?: string;
                'e'?: string;
                'kid'?: string;
                'kty'?: string;
                'n'?: string;
                'use'?: string;
                'x'?: string;
                'y'?: string;
          }
        
  
        type OAuthTokenRequest =
          {
                'client_id'?: string;
                'client_secret': string;
                'code'?: string;
                'disable_refresh_token'?: boolean;
                'grant_type'?: OAuthGrantType;
                'password': string;
                'redirect_uri'?: string;
                'refresh_token'?: string;
                'state'?: string;
                'token_type'?: OAuthTokenType;
                'username'?: string;
          }
        
  
        type OAuthTokenResponse =
          {
                'access_token'?: string;
                'cookies'?: string[];
                'error'?: string;
                'expires_in': number;
                'headers'?: Record<string, any>;
                'id_token'?: string;
                'next_method'?: LoginType[];
                'refresh_token'?: string;
                'token_type': OAuthTokenType;
          }
        
  
        type OAuthTokenType = "Bearer" | 0 | "Mac" | 1 | "Cookie" | 2;
  
        type OIDCWellKnownResponse =
          {
                'authorization_endpoint'?: string;
                'claims_supported'?: string[];
                'code_challenge_methods_supported'?: string[];
                'device_authorization_endpoint'?: string;
                'grant_types_supported'?: string[];
                'id_token_signing_alg_values_supported'?: string[];
                'issuer'?: string;
                'jwks_uri'?: string;
                'response_types_supported'?: string[];
                'revocation_endpoint'?: string;
                'scopes_supported'?: string[];
                'subject_types_supported'?: string[];
                'token_endpoint'?: string;
                'token_endpoint_auth_methods_supported'?: string[];
                'userinfo_endpoint'?: string;
          }
        
  
        type openidConfigurationParams =
          {
                'client_id'?: string;
          }
        
  
        type openidConfigurationParams =
          {
                'client_id'?: string;
          }
        
  
        type PageConfig =
          {
                'createTime': string;
                'description'?: string;
                'fields'?: FieldConfig[];
                'icon'?: string;
                'id': string;
                'isDisable': boolean;
                'name': string;
                'updateTime': string;
          }
        
  
        type PageData =
          {
                'createTime': string;
                'data'?: Record<string, any>;
                'id': string;
                'pageId': string;
                'updateTime': string;
          }
        
  
        type PageStatus = "all" | "disabled" | "enabled";
  
        type PasswordComplexity = "unsafe" | 0 | "general" | 1 | "safe" | 2 | "very_safe" | 3;
  
        type patchAppParams =
          {
              /** identifier of the app */
                'id': string;
          }
        
  
        type PatchAppRequest =
          {
                'avatar'?: string;
                'description'?: string;
                'displayName'?: string;
                'grantMode'?: AppMetaGrantMode;
                'grantType'?: AppMetaGrantType[];
                'i18n'?: AppI18NOptions;
                'id': string;
                'isDelete'?: boolean;
                'name'?: string;
                'status'?: AppMetaStatus;
                'url'?: string;
          }
        
  
        type PatchCurrentUserRequest =
          {
                'email_as_mfa'?: boolean;
                'sms_as_mfa'?: boolean;
                'totp_as_mfa'?: boolean;
          }
        
  
        type PatchPageDataRequest =
          {
                'data'?: Record<string, any>;
                'id': string;
                'isDelete'?: boolean;
                'pageId': string;
          }
        
  
        type patchPageDatasParams =
          {
              /** identifier of the page */
                'pageId': string;
          }
        
  
        type PatchPageRequest =
          {
                'description'?: string;
                'fields'?: FieldConfig[];
                'icon'?: string;
                'id': string;
                'isDelete'?: boolean;
                'isDisable'?: boolean;
                'name'?: string;
          }
        
  
        type PatchSecurityConfigRequest =
          {
                'accountInactiveLock'?: number;
                'forceEnableMfa'?: boolean;
                'loginSessionInactivityTime'?: number;
                'loginSessionMaxTime'?: number;
                'passwordComplexity'?: PasswordComplexity;
                'passwordExpireTime'?: number;
                'passwordFailedLockDuration'?: number;
                'passwordFailedLockThreshold'?: number;
                'passwordHistory'?: number;
                'passwordMinLength'?: number;
          }
        
  
        type patchUserParams =
          {
              /** identifier of the user */
                'id': string;
          }
        
  
        type PatchUserRequest =
          {
                'id'?: string;
                'isDelete'?: boolean;
                'status'?: UserMetaUserStatus;
          }
        
  
        type PatchUserResponse =
          {
                'User': string;
          }
        
  
        type PermissionInfo =
          {
                'createTime': string;
                'description'?: string;
                'enableAuth'?: boolean;
                'id': string;
                'name'?: string;
                'parentId'?: string;
                'updateTime': string;
          }
        
  
        type ResetUserPasswordRequest =
          {
                'newPassword': string;
                'oldPassword'?: string;
                'token'?: string;
                'userId': string;
                'username'?: string;
          }
        
  
        type RoleInfo =
          {
                'createTime': string;
                'description'?: string;
                'id': string;
                'name': string;
                'permission'?: PermissionInfo[];
                'updateTime': string;
          }
        
  
        type RuntimeSecurityConfig =
          {
                'accountInactiveLock': number;
                'forceEnableMfa': boolean;
                'loginSessionInactivityTime': number;
                'loginSessionMaxTime': number;
                'passwordComplexity': PasswordComplexity;
                'passwordExpireTime': number;
                'passwordFailedLockDuration': number;
                'passwordFailedLockThreshold': number;
                'passwordHistory': number;
                'passwordMinLength': number;
          }
        
  
        type SendActivationMailRequest =
          {
                'userId': string;
          }
        
  
        type SendLoginCaptchaRequest =
          {
                'email'?: string;
                'phone'?: string;
                'type': LoginType;
                'username'?: string;
          }
        
  
        type SendLoginCaptchaResponse =
          {
                'data'?: SendLoginCaptchaResponseData;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type SendLoginCaptchaResponseData =
          {
                'token'?: string;
          }
        
  
        type SessionInfo =
          {
                'createTime': string;
                'expiry': string;
                'id': string;
                'lastSeen'?: string;
          }
        
  
        type SimpleAppKeyInfo =
          {
                'appId': string;
                'createTime': string;
                'id': string;
                'key': string;
                'name': string;
                'updateTime': string;
          }
        
  
        type updateAppParams =
          {
              /** identifier of the app */
                'id': string;
          }
        
  
        type UpdateAppRequest =
          {
                'avatar'?: string;
                'description'?: string;
                'displayName'?: string;
                'grantMode'?: AppMetaGrantMode;
                'grantType'?: AppMetaGrantType[];
                'i18n'?: AppI18NOptions;
                'id': string;
                'isDelete'?: boolean;
                'name': string;
                'oAuth2': AppOAuth2Info;
                'proxy'?: AppProxyInfo;
                'roles'?: AppRoleInfo[];
                'status': AppMetaStatus;
                'url': string;
                'users'?: AppUser[];
          }
        
  
        type updatePageDataParams =
          {
              /** identifier of the page */
                'pageId': string;
              /** data id of the page */
                'id': string;
          }
        
  
        type UpdatePageDataRequest =
          {
                'data'?: Record<string, any>;
                'id': string;
                'pageId': string;
          }
        
  
        type updatePageParams =
          {
              /** identifier of the page */
                'id': string;
          }
        
  
        type UpdatePageRequest =
          {
                'description'?: string;
                'fields'?: FieldConfig[];
                'icon'?: string;
                'id': string;
                'isDisable': boolean;
                'name': string;
          }
        
  
        type updateRoleParams =
          {
              /** identifier of the role */
                'id': string;
          }
        
  
        type UpdateRoleRequest =
          {
                'description'?: string;
                'id': string;
                'name': string;
                'permission'?: string[];
          }
        
  
        type updateUserParams =
          {
              /** identifier of the user */
                'id': string;
          }
        
  
        type UpdateUserRequest =
          {
                'apps'?: UserApp[];
                'avatar'?: string;
                'email'?: string;
                'fullName'?: string;
                'id': string;
                'isDelete'?: boolean;
                'phoneNumber'?: string;
                'status'?: UserMetaUserStatus;
                'username': string;
          }
        
  
        type uploadFileParams =
          {
              /** files */
                'files'?: string[];
          }
        
  
        type User =
          {
                'apps'?: App[];
                'avatar': string;
                'createTime'?: string;
                'email': string;
                'extendedData'?: UserExt;
                'fullName': string;
                'id': string;
                'isDelete'?: boolean;
                'loginTime': string;
                'password'?: string;
                'phoneNumber': string;
                'role'?: string;
                'roleId'?: string;
                'status': UserMetaUserStatus;
                'updateTime'?: string;
                'username': string;
          }
        
  
        type UserApp =
          {
                'avatar'?: string;
                'description'?: string;
                'displayName'?: string;
                'id': string;
                'name'?: string;
                'role'?: string;
                'roleId'?: string;
                'roles'?: AppRole[];
          }
        
  
        type UserExt =
          {
                'ForceMFA': boolean;
                'activationTime': string;
                'emailAsMFA': boolean;
                'loginTime': string;
                'passwordModifyTime': string;
                'smsAsMFA': boolean;
                'totpAsMFA': boolean;
                'userId': string;
          }
        
  
        type UserInfo =
          {
                'apps'?: UserApp[];
                'avatar'?: string;
                'createTime': string;
                'email'?: string;
                'extendedData'?: UserExt;
                'fullName'?: string;
                'id': string;
                'isDelete': boolean;
                'loginTime'?: string;
                'phoneNumber'?: string;
                'role'?: string;
                'roleId'?: string;
                'status': UserMetaUserStatus;
                'updateTime': string;
                'username': string;
          }
        
  
        type UserLoginRequest =
          {
                'autoLogin'?: boolean;
                'bindingToken'?: string;
                'code'?: string;
                'email'?: string;
                'firstCode'?: string;
                'password'?: string;
                'phone'?: string;
                'secondCode'?: string;
                'token'?: string;
                'type'?: LoginType;
                'username'?: string;
          }
        
  
        type UserLoginResponse =
          {
                'data'?: UserLoginResponseData;
                'errorCode'?: string;
                'errorMessage'?: string;
                'success': boolean;
                'traceId': string;
          }
        
  
        type UserLoginResponseData =
          {
                'email'?: string;
                'nextMethod': LoginType[];
                'phone_number'?: string;
                'token'?: string;
          }
        
  
        type UserMetaUserStatus = "normal" | 0 | "disabled" | 1 | "user_inactive" | 2 | "password_expired" | 4;
  
        type userOAuthLoginParams =
          {
              /** identifier of the oauth */
                'id': string;
          }
        
  
}
