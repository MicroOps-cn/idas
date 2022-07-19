package global

import "time"

const (
	AppName                     = "idas"
	IdasAppName                 = "IDAS"
	LoginSession                = "LOGIN_SESSION"
	AuthCode                    = "AUTH_CODE"
	LoggerName                  = "__logger__"
	TraceIdName                 = "traceId"
	CallerName                  = "caller"
	RestfulRequestContextName   = "__restful_request__"
	RestfulResponseContextName  = "__restful_response__"
	MetaUser                    = "__user__"
	MetaNeedLogin               = "__need_login__"
	MetaAutoRedirectToLoginPage = "__auto_redirect_to_login_page__"
	GormConnName                = "__mysql_conn__"
	LDAPConnName                = "__ldap_conn__"
	LoginSessionExpiration      = 7 * 24 * time.Hour
	AuthCodeExpiration          = 5 * time.Minute
	TokenExpiration             = 1 * time.Hour
	RefreshTokenExpiration      = 30 * 24 * time.Hour
	ResetPasswordExpiration     = 30 * time.Minute
	LoginSessionExpiresFormat   = "Mon, 02-Jan-06 15:04:05 MST"
)
