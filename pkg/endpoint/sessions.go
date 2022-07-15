package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service"
	"idas/pkg/service/models"
	"net/http"
	"strings"
	"time"
)

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UserLoginRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		if loginCookie, err := s.CreateLoginSession(ctx, req.Username, req.Password, req.AutoLogin); err == nil {
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", strings.Join(loginCookie, ","))
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Wrong user name or password")
		}
		return &resp, nil
	}
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(Requester).GetRequestData().(*UserLogoutRequest)
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
		return s.Authentication(ctx, req.AuthMethod, req.AuthAlgorithm, req.AuthKey, req.AuthSecret)
	}
}

func MakeGetLoginSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sessionId := request.([]string)
		var resp []*models.User
		if len(sessionId) > 0 {
			if resp, err = s.GetLoginSession(ctx, sessionId); err != nil {
				err = errors.NotLoginError
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
		resp.Data, resp.Total, resp.BaseResponse.Error = s.GetSessions(ctx, req.UserId, req.Current, req.PageSize)
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
