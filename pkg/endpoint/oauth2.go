package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service"
	"idas/pkg/service/models"
	"net/url"
	"strconv"
)

type OAuthTokenResponse struct {
	Error        string `json:"error"`
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: "Bearer"}
		if restfulReq := request.(RestfulRequester).GetRestfulRequest(); restfulReq != nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			switch req.GrantType {
			case OAuthTokenRequest_AuthorizationCode:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByAuthorizationCode(ctx, req.Code, req.ClientId)
			case OAuthTokenRequest_Password:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, req.Username, req.Password)
			case OAuthTokenRequest_ClientCredentials:
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, username, password)
				} else {
					err = fmt.Errorf("invalid_request")
				}
			case OAuthTokenRequest_RefreshToken:
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, username, password)
				} else if len(req.Username) != 0 && len(req.Password) != 0 {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, req.Username, req.Password)
				} else {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByAuthorizationCode(ctx, req.RefreshToken, req.ClientId, req.ClientSecret)
				}
			default:
				err = fmt.Errorf("unsupported_grant_type")
			}
		}

		if err != nil {
			resp.Error = err.Error()
			if restfulResp := request.(RestfulRequester).GetRestfulResponse(); restfulResp != nil {
				restfulResp.WriteHeader(400)
			}
		}
		return &resp, nil
	}
}

func MakeOAuthAuthorizeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*OAuthAuthorizeRequest)
		resp := BaseResponse[interface{}]{}
		var code string

		stdResp := request.(RestfulRequester).GetRestfulResponse()

		if len(req.ClientId) == 0 {
			return nil, errors.ParameterError("client_id")
		}
		users, ok := request.(RestfulRequester).GetRestfulRequest().Attribute(global.AttrUser).([]*models.User)
		if !ok || len(users) == 0 {
			return nil, errors.NotLoginError
		}

		uri, err := url.Parse(req.RedirectUri)
		if err != nil {
			return nil, errors.ParameterError("redirect_uri")
		}

		for _, user := range users {
			if code, err = s.GetAuthCodeByClientId(ctx, req.ClientId, user.Id, user.Storage); err != nil {
				return nil, err
			}
		}

		query := uri.Query()
		switch req.ResponseType {
		case OAuthAuthorizeRequest_code, OAuthAuthorizeRequest_default:
			query.Add("code", code)
			uri.RawQuery = query.Encode()
			stdResp.AddHeader("Location", uri.String())
			stdResp.ResponseWriter.WriteHeader(302)
		case OAuthAuthorizeRequest_token:
			accessToken, refreshToken, expiresIn, err := s.GetOAuthTokenByAuthorizationCode(ctx, code, req.ClientId)
			if err != nil {
				return nil, err
			}
			query.Add("access_token", accessToken)
			query.Add("refresh_token", refreshToken)
			query.Add("expires_in", strconv.Itoa(expiresIn))
			uri.RawQuery = query.Encode()
			stdResp.AddHeader("Location", uri.String())
			stdResp.ResponseWriter.WriteHeader(302)
		}
		return &resp, nil
	}
}
