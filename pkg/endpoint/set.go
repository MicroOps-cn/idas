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
	CurrentUser endpoint.Endpoint
}

type AppEndpoints struct{}

type AuthEndpoints struct {
	UserLogin,
	UserLogout,
	GetLoginSession,
	OAuthTokens,
	OAuthAuthorize endpoint.Endpoint
}

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	UserEndpoints
	AuthEndpoints
	AppEndpoints
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	return Set{
		UserEndpoints: UserEndpoints{
			CurrentUser:   InjectEndpoint(logger, "CurrentUser", duration, otTracer, zipkinTracer, MakeCurrentUserEndpoint(svc)),
			GetUsers:      InjectEndpoint(logger, "GetUsers", duration, otTracer, zipkinTracer, MakeGetUsersEndpoint(svc)),
			DeleteUsers:   InjectEndpoint(logger, "DeleteUsers", duration, otTracer, zipkinTracer, MakeDeleteUsersEndpoint(svc)),
			PatchUsers:    InjectEndpoint(logger, "PatchUsers", duration, otTracer, zipkinTracer, MakePatchUsersEndpoint(svc)),
			UpdateUser:    InjectEndpoint(logger, "UpdateUser", duration, otTracer, zipkinTracer, MakeUpdateUserEndpoint(svc)),
			GetUserInfo:   InjectEndpoint(logger, "GetUserInfo", duration, otTracer, zipkinTracer, MakeGetUserInfoEndpoint(svc)),
			CreateUser:    InjectEndpoint(logger, "CreateUser", duration, otTracer, zipkinTracer, MakeCreateUserEndpoint(svc)),
			PatchUser:     InjectEndpoint(logger, "PatchUser", duration, otTracer, zipkinTracer, MakePatchUserEndpoint(svc)),
			DeleteUser:    InjectEndpoint(logger, "DeleteUser", duration, otTracer, zipkinTracer, MakeDeleteUserEndpoint(svc)),
			GetUserSource: InjectEndpoint(logger, "DeleteUser", duration, otTracer, zipkinTracer, MakeGetUserSourceRequestEndpoint(svc)),
		},
		AuthEndpoints: AuthEndpoints{
			UserLogin:       InjectEndpoint(logger, "UserLogin", duration, otTracer, zipkinTracer, MakeUserLoginEndpoint(svc)),
			UserLogout:      InjectEndpoint(logger, "UserLogout", duration, otTracer, zipkinTracer, MakeUserLogoutEndpoint(svc)),
			GetLoginSession: InjectEndpoint(logger, "GetLoginSession", duration, otTracer, zipkinTracer, MakeGetLoginSessionEndpoint(svc)),
			OAuthTokens:     InjectEndpoint(logger, "OAuthTokens", duration, otTracer, zipkinTracer, MakeOAuthTokensEndpoint(svc)),
			OAuthAuthorize:  InjectEndpoint(logger, "OAuthAuthorize", duration, otTracer, zipkinTracer, MakeOAuthAuthorizeEndpoint(svc)),
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
