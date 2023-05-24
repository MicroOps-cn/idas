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
	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/log/level"
	"net/http"
	"strings"
)

type Client struct {
	o *Option
}

type UserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

type TokenRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func (c *Client) GetToken(ctx context.Context, code string) (*TokenResponse, error) {
	reqData, err := json.Marshal(&TokenRequest{Code: code, GrantType: "authorization_code", ClientId: c.o.ClientId, ClientSecret: c.o.ClientSecret})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.o.TokenUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to create token request")
	}
	req.SetBasicAuth(c.o.ClientId, c.o.ClientSecret)
	req.Header.Set("content-type", restful.MIME_JSON)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to get token")
	}
	if resp.StatusCode != 200 {
		return nil, errors.WithServerError(500, fmt.Errorf("response code: %d", resp.StatusCode), "failed to get token: statusCode:"+resp.Status)
	}
	defer resp.Body.Close()
	var tokenResp TokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		level.Error(logs.GetContextLogger(ctx)).Log("err", "err", "msg", "failed to decode oauth token")
		return nil, errors.WithServerError(500, err, "failed to decode oauth token")
	}
	return &tokenResp, nil
}
func (c *Client) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	tokenData, err := c.GetToken(ctx, code)
	if err != nil {
		return nil, err
	}
	var userInfo UserInfo
	if len(tokenData.IdToken) > 0 {
		fmt.Println(tokenData.IdToken)
		chunks := strings.Split(tokenData.IdToken, ".")
		if len(chunks) == 3 {
			idInfo, err := base64.StdEncoding.DecodeString(chunks[1])
			if err == nil {
				fmt.Println(idInfo)
				err = json.Unmarshal(idInfo, &userInfo)
				if err == nil {
					if userInfo.Username != "" {
						return &userInfo, nil
					}
				}
			}
		}
	}
	userInfoReq, err := http.NewRequest("GET", c.o.ApiUrl, nil)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to get user info: failed to create request")
	}
	fmt.Println(tokenData.AccessToken)
	userInfoReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userInfoResp, err := http.DefaultClient.Do(userInfoReq)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to get user info")
	}

	if userInfoResp.StatusCode != 200 {
		return nil, errors.WithServerError(500, fmt.Errorf("response code: %d", userInfoResp.StatusCode), "failed to get userinfo")
	}
	defer userInfoResp.Body.Close()
	err = json.NewDecoder(userInfoResp.Body).Decode(&userInfo)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to decode user info")
	}
	return &userInfo, nil
}

func NewClient(o *Option) *Client {
	if o == nil {
		return nil
	}
	return &Client{o: o}
}
