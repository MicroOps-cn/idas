package endpoint

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"

	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UserLoginRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		if loginCookie, err := s.CreateLoginSession(ctx, req.Username, req.Password, req.AutoLogin); err == nil {
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", loginCookie)
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Wrong user name or password")
		}
		return &resp, nil
	}
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		cookie, err := request.(RestfulRequester).GetRestfulRequest().Request.Cookie(global.LoginSession)
		if err != nil {
			resp.Error = errors.BadRequestError
		} else if len(cookie.Value) > 0 {
			for _, id := range strings.Split(cookie.Value, ",") {
				if err = s.DeleteLoginSession(ctx, id); err != nil {
					resp.Error = errors.InternalServerError
					return resp, nil
				}
			}
			loginCookie := fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, cookie.Value, time.Now().UTC().Format(global.LoginSessionExpiresFormat))
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", fmt.Sprintf(loginCookie))
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Invalid identity information")
		}
		return &resp, nil
	}
}

func MakeAuthenticationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*AuthenticationRequest)
		return s.Authentication(ctx, req.AuthMethod, req.AuthAlgorithm, req.AuthKey, req.AuthSecret, req.Payload, req.AuthSign)
	}
}

type GetSessionParams struct {
	Token     string
	TokenType models.TokenType
}

func MakeGetSessionByTokenEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		params := request.(GetSessionParams)
		var resp []*models.User
		if len(params.Token) > 0 {
			if resp, err = s.GetSessionByToken(ctx, params.Token, params.TokenType); err != nil {
				if err != errors.NotLoginError {
					level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get session")
					err = errors.NotLoginError
				}
			}
		} else {
			err = errors.NotLoginError
		}
		return resp, err
	}
}

func MakeGetSessionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetSessionsRequest)
		resp := NewBaseListResponse[[]*models.Token](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetSessions(ctx, req.UserId, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeDeleteSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteSessionRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteSession(ctx, req.Id)
		return &resp, nil
	}
}
