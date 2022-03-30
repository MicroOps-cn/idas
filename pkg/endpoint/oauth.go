package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"idas/pkg/service"
)

type OAuthGrantType string

const (
	OAuthTypeRefreshToken      OAuthGrantType = "refresh_token"
	OAuthTypeAuthorizationCode OAuthGrantType = "authorization_code"
	OAuthTypePassword          OAuthGrantType = "password"
	OAuthTypeClientCredentials OAuthGrantType = "client_credentials"
)

type OAuthTokenRequest struct {
	BaseRequest
	Code         string         `json:"code"`
	GrantType    OAuthGrantType `json:"grant_type"`
	RedirectURI  string         `json:"redirect_uri"`
	ClientId     string         `json:"client_id"`
	ClientSecret string         `json:"client_secret"`
	Password     string         `json:"password"`
	Username     string         `json:"username"`
	RefreshToken string         `json:"refresh_token"`
}

type OAuthTokenResponse struct {
	Error        string `json:"error"`
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: "Bearer"}
		if restfulReq := req.GetRestfulRequest(); restfulReq != nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			switch req.GrantType {
			case OAuthTypeAuthorizationCode:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByAuthorizationCode(ctx, req.Code, req.ClientId, req.RedirectURI)
			case OAuthTypePassword:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, req.Username, req.Password)
			case OAuthTypeClientCredentials:
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, username, password)
				} else {
					err = fmt.Errorf("invalid_request")
				}
			case OAuthTypeRefreshToken:
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
			if restfulResp := req.GetRestfulResponse(); restfulResp != nil {
				restfulResp.WriteHeader(400)
			}
		}
		return &resp, nil
	}
}

type OAuthAuthorizeRequest struct {
	BaseRequest
	ResponseType string `json:"response_type"`
	ClientId     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
}

type OAuthAuthorizeResponse struct {
	BaseResponse `json:",inline"`
}

func MakeOAuthAuthorizeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*OAuthAuthorizeRequest)
		resp := OAuthAuthorizeResponse{}
		var redirect string
		if redirect, err = s.OAuthAuthorize(ctx, req.ResponseType, req.ClientId, req.RedirectURI); err == nil {
			if req.GetRestfulResponse() != nil && len(redirect) > 0 {
				req.GetRestfulResponse().AddHeader("Location", redirect)
				req.GetRestfulResponse().ResponseWriter.WriteHeader(302)
			}
		} else {
			resp.Error = err
		}
		return &resp, nil
	}
}
