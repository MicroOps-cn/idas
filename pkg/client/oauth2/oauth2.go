/*
 Copyright © 2023 MicroOps-cn.

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

package oauth2

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	http2 "github.com/MicroOps-cn/fuck/http"
	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/tidwall/gjson"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
)

type Client struct {
	o *Options
}

type UserInfo struct {
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Avatar      string `json:"avatar"`
}

func (i UserInfo) String() string {
	return fmt.Sprintf("%#v", i)
}

type TokenRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RedirectURI  string `json:"redirect_uri"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func (c *Client) GetToken(ctx context.Context, code string, redirectURI string) (*TokenResponse, error) {
	secret, err := c.o.ClientSecret.UnsafeString()
	if err != nil {
		return nil, err
	}
	redirectURL, err := c.o.GetRedirectURL(ctx, redirectURI)
	if err != nil {
		return nil, err
	}
	reqData, err := json.Marshal(&TokenRequest{Code: code, GrantType: "authorization_code", ClientId: c.o.ClientId, ClientSecret: secret, RedirectURI: redirectURL.String()})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.o.TokenUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to create token request")
	}
	req.SetBasicAuth(c.o.ClientId, secret)
	req.Header.Set("content-type", restful.MIME_JSON)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		level.Warn(logs.GetContextLogger(ctx)).Log("msg", "failed to get token, switch content-type to application/x-www-form-urlencoded")
		req, err = http.NewRequest("POST", c.o.TokenUrl, bytes.NewBuffer([]byte(url.Values{
			"grant_type":    []string{"authorization_code"},
			"client_id":     []string{c.o.ClientId},
			"client_secret": []string{secret},
			"code":          []string{code},
			"redirect_uri":  []string{redirectURL.String()},
		}.Encode())))
		if err != nil {
			return nil, errors.WithServerError(500, err, "failed to create token request")
		}
		req.SetBasicAuth(c.o.ClientId, secret)
		req.Header.Set("content-type", restful.MIME_JSON)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		} else if resp.StatusCode != 200 {
			return nil, errors.WithServerError(500, fmt.Errorf("response code: %d", resp.StatusCode), "failed to get token: statusCode:"+resp.Status)
		}
	}
	defer resp.Body.Close()
	var tokenResp TokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		level.Error(logs.GetContextLogger(ctx)).Log("err", "err", "msg", "failed to decode oauth token")
		return nil, errors.WithServerError(500, err, "failed to decode oauth token")
	}
	return &tokenResp, nil
}

func (c *Client) decodeUserInfo(_ context.Context, raw string) *UserInfo {
	var userInfo UserInfo
	j := gjson.Parse(raw)
	userInfo.Role = j.Get(c.o.RoleAttributePath).String()
	userInfo.PhoneNumber = j.Get(c.o.PhoneNumberAttributePath).String()
	userInfo.Email = j.Get(c.o.EmailAttributePath).String()
	userInfo.FullName = j.Get(c.o.FullNameAttributePath).String()
	userInfo.Username = j.Get(c.o.UsernameAttributePath).String()
	userInfo.Avatar = j.Get(c.o.AvatarAttributePath).String()
	return &userInfo
}

func (c *Client) GetUserInfo(ctx context.Context, code string, redirectURI string) (*UserInfo, error) {
	tokenData, err := c.GetToken(ctx, code, redirectURI)
	if err != nil {
		return nil, err
	}
	var userInfo *UserInfo
	if len(tokenData.IdToken) > 0 {
		chunks := strings.Split(tokenData.IdToken, ".")
		if len(chunks) == 3 {
			idInfo, err := base64.StdEncoding.DecodeString(chunks[1])
			if err == nil {
				userInfo = c.decodeUserInfo(ctx, string(idInfo))
				if len(userInfo.Username) > 0 || len(userInfo.Email) > 0 || len(userInfo.PhoneNumber) > 0 {
					return userInfo, nil
				}
			}
		}
	}
	userInfoReq, err := http.NewRequest("GET", c.o.ApiUrl, nil)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to get user info: failed to create request")
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userInfoResp, err := http.DefaultClient.Do(userInfoReq)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to get user info")
	}

	if userInfoResp.StatusCode != 200 {
		return nil, errors.WithServerError(500, fmt.Errorf("response code: %d", userInfoResp.StatusCode), "failed to get userinfo")
	}
	defer userInfoResp.Body.Close()
	raw, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return nil, err
	}
	userInfo = c.decodeUserInfo(ctx, string(raw))
	return userInfo, nil
}

func NewClient(o *Options) *Client {
	if o == nil {
		return nil
	}
	return &Client{o: o}
}

type pbOptions Options

func (p *pbOptions) Reset() {
	(*Options)(p).Reset()
}

func (p *pbOptions) String() string {
	return (*Options)(p).String()
}

func (p *pbOptions) ProtoMessage() {
	(*Options)(p).Reset()
}

func (x *Options) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewOptions()
	x.LoginId = options.LoginId
	x.EmailAttributePath = options.EmailAttributePath
	x.UsernameAttributePath = options.UsernameAttributePath
	x.PhoneNumberAttributePath = options.PhoneNumberAttributePath
	x.FullNameAttributePath = options.FullNameAttributePath
	err := unmarshaller.Unmarshal(bytes.NewReader(b), (*pbOptions)(x))
	if err != nil {
		return err
	}
	return nil
}

// NewOptions return a default option
// which host field point to nowhere.
func NewOptions() *Options {
	return &Options{
		LoginId:                  "username",
		EmailAttributePath:       "email",
		UsernameAttributePath:    "username",
		PhoneNumberAttributePath: "phoneNumber",
		FullNameAttributePath:    "fullName",
		RoleAttributePath:        "role",
	}
}

func (x *Options) GetRedirectURL(ctx context.Context, oriRedirectURI string) (*url.URL, error) {
	httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string)
	if !ok {
		return nil, errors.StatusNotFound("httpExternalURL")
	}
	redirectURI, err := url.Parse(httpExternalURL)
	if !ok {
		return nil, err
	}
	redirectURI.Path = http2.JoinPath(redirectURI.Path, "api/v1/user/oauth/"+x.Id)
	if len(oriRedirectURI) > 0 {
		rq := redirectURI.Query()
		rq.Set("redirect_uri", oriRedirectURI)
		redirectURI.RawQuery = rq.Encode()
	}
	return redirectURI, nil
}
