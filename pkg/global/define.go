package global

import "time"

const (
	AppName                    = "idas"
	LoginSession               = "LOGIN_SESSION"
	LoggerName                 = "__logger__"
	TraceIdName                = "traceId"
	CallerName                 = "caller"
	RestfulRequestContextName  = "__restful_request__"
	RestfulResponseContextName = "__restful_response__"
	MetaNeedLogin              = "__need_login__"
	MySQLConnName              = "__mysql_conn__"
	AttrUser                   = "__user__"
	LoginSessionExpiration     = 7 * 24 * time.Hour
	LoginSessionExpiresFormat  = "Mon, 02-Jan-06 15:04:05 MST"
)
