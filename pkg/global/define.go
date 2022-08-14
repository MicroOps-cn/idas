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

package global

import "time"

const (
	AppName                     = "idas"
	IdasAppName                 = "IDAS"
	LoginSession                = "LOGIN_SESSION"
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
	ActiveExpiration            = 7 * 24 * time.Hour
	AuthCodeExpiration          = 5 * time.Minute
	TokenExpiration             = 1 * time.Hour
	RefreshTokenExpiration      = 30 * 24 * time.Hour
	ResetPasswordExpiration     = 30 * time.Minute
	LoginSessionExpiresFormat   = "Mon, 02-Jan-06 15:04:05 MST"
)
