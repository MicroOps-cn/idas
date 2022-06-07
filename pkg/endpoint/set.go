package endpoint

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"idas/pkg/service"
)

type UserEndpoints struct {
	GetUsers,
	PatchUsers,
	DeleteUsers,
	UpdateUser,
	GetUserInfo,
	CreateUser,
	PatchUser,
	DeleteUser,
	GetUserSource,
	ForgotPassword,
	ResetPassword,
	CurrentUser endpoint.Endpoint
}

type AppEndpoints struct {
	PatchApps,
	DeleteApps,
	GetAppSource,
	GetAppInfo,
	CreateApp,
	UpdateApp,
	PatchApp,
	DeleteApp,
	GetApps endpoint.Endpoint
}

type SessionEndpoints struct {
	GetSessions,
	DeleteSession,
	UserLogin,
	UserLogout,
	GetLoginSession,
	OAuthTokens,
	OAuthAuthorize endpoint.Endpoint
}

type CommonEndpoints struct {
	UploadFile   endpoint.Endpoint
	DownloadFile endpoint.Endpoint
}

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	UserEndpoints
	SessionEndpoints
	AppEndpoints
	CommonEndpoints
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	return Set{
		CommonEndpoints: CommonEndpoints{
			UploadFile:   InjectEndpoint(logger, "UploadFile", duration, otTracer, zipkinTracer, MakeUploadFileEndpoint(svc)),
			DownloadFile: InjectEndpoint(logger, "UploadFile", duration, otTracer, zipkinTracer, MakeDownloadFileEndpoint(svc)),
		},
		UserEndpoints: UserEndpoints{
			CurrentUser:    InjectEndpoint(logger, "CurrentUser", duration, otTracer, zipkinTracer, MakeCurrentUserEndpoint(svc)),
			ResetPassword:  InjectEndpoint(logger, "ResetPassword", duration, otTracer, zipkinTracer, MakeResetUserPasswordEndpoint(svc)),
			ForgotPassword: InjectEndpoint(logger, "ForgotPassword", duration, otTracer, zipkinTracer, MakeForgotPasswordEndpoint(svc)),
			GetUsers:       InjectEndpoint(logger, "GetUsers", duration, otTracer, zipkinTracer, MakeGetUsersEndpoint(svc)),
			DeleteUsers:    InjectEndpoint(logger, "DeleteUsers", duration, otTracer, zipkinTracer, MakeDeleteUsersEndpoint(svc)),
			PatchUsers:     InjectEndpoint(logger, "PatchUsers", duration, otTracer, zipkinTracer, MakePatchUsersEndpoint(svc)),
			UpdateUser:     InjectEndpoint(logger, "UpdateUser", duration, otTracer, zipkinTracer, MakeUpdateUserEndpoint(svc)),
			GetUserInfo:    InjectEndpoint(logger, "GetUserInfo", duration, otTracer, zipkinTracer, MakeGetUserInfoEndpoint(svc)),
			CreateUser:     InjectEndpoint(logger, "CreateUser", duration, otTracer, zipkinTracer, MakeCreateUserEndpoint(svc)),
			PatchUser:      InjectEndpoint(logger, "PatchUser", duration, otTracer, zipkinTracer, MakePatchUserEndpoint(svc)),
			DeleteUser:     InjectEndpoint(logger, "DeleteUser", duration, otTracer, zipkinTracer, MakeDeleteUserEndpoint(svc)),
			GetUserSource:  InjectEndpoint(logger, "GetUserSource", duration, otTracer, zipkinTracer, MakeGetUserSourceRequestEndpoint(svc)),
		},
		SessionEndpoints: SessionEndpoints{
			GetSessions:     InjectEndpoint(logger, "GetSessions", duration, otTracer, zipkinTracer, MakeGetSessionsEndpoint(svc)),
			DeleteSession:   InjectEndpoint(logger, "DeleteSession", duration, otTracer, zipkinTracer, MakeDeleteSessionEndpoint(svc)),
			UserLogin:       InjectEndpoint(logger, "UserLogin", duration, otTracer, zipkinTracer, MakeUserLoginEndpoint(svc)),
			UserLogout:      InjectEndpoint(logger, "UserLogout", duration, otTracer, zipkinTracer, MakeUserLogoutEndpoint(svc)),
			GetLoginSession: InjectEndpoint(logger, "GetLoginSession", duration, otTracer, zipkinTracer, MakeGetLoginSessionEndpoint(svc)),
			OAuthTokens:     InjectEndpoint(logger, "OAuthTokens", duration, otTracer, zipkinTracer, MakeOAuthTokensEndpoint(svc)),
			OAuthAuthorize:  InjectEndpoint(logger, "OAuthAuthorize", duration, otTracer, zipkinTracer, MakeOAuthAuthorizeEndpoint(svc)),
		},
		AppEndpoints: AppEndpoints{
			GetApps:      InjectEndpoint(logger, "GetApps", duration, otTracer, zipkinTracer, MakeGetAppsEndpoint(svc)),
			DeleteApps:   InjectEndpoint(logger, "DeleteApps", duration, otTracer, zipkinTracer, MakeDeleteAppsEndpoint(svc)),
			PatchApps:    InjectEndpoint(logger, "PatchApps", duration, otTracer, zipkinTracer, MakePatchAppsEndpoint(svc)),
			UpdateApp:    InjectEndpoint(logger, "UpdateApp", duration, otTracer, zipkinTracer, MakeUpdateAppEndpoint(svc)),
			GetAppInfo:   InjectEndpoint(logger, "GetAppInfo", duration, otTracer, zipkinTracer, MakeGetAppInfoEndpoint(svc)),
			CreateApp:    InjectEndpoint(logger, "CreateApp", duration, otTracer, zipkinTracer, MakeCreateAppEndpoint(svc)),
			PatchApp:     InjectEndpoint(logger, "PatchApp", duration, otTracer, zipkinTracer, MakePatchAppEndpoint(svc)),
			DeleteApp:    InjectEndpoint(logger, "DeleteApp", duration, otTracer, zipkinTracer, MakeDeleteAppEndpoint(svc)),
			GetAppSource: InjectEndpoint(logger, "GetAppSource", duration, otTracer, zipkinTracer, MakeGetAppSourceRequestEndpoint(svc)),
		},
	}
}

func InjectEndpoint(logger log.Logger, name string, duration metrics.Histogram, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, ep endpoint.Endpoint) endpoint.Endpoint {
	ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(ep)
	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	ep = opentracing.TraceServer(tracer, name)(ep)
	if zipkinTracer != nil {
		ep = zipkin.TraceEndpoint(zipkinTracer, name)(ep)
	}
	ep = LoggingMiddleware(name)(ep)
	ep = InstrumentingMiddleware(duration.With("method", "Concat"))(ep)
	return ep
}
