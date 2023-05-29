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

package opts

import "github.com/MicroOps-cn/idas/pkg/errors"

type WithGetUserOptions func(o *GetUserOptions)

type GetUserOptions struct {
	Id          string
	Username    string
	Email       string
	PhoneNumber string
	Ext         bool
	Err         error
	Apps        bool
}

func WithPhoneNumber(no string) WithGetUserOptions {
	return func(o *GetUserOptions) {
		o.PhoneNumber = no
		if len(o.PhoneNumber) == 0 && o.Err != nil {
			o.Err = errors.LackParameterError("phoneNumber")
		}
	}
}

func WithUsername(username string) WithGetUserOptions {
	return func(o *GetUserOptions) {
		o.Username = username
		if len(o.Username) == 0 && o.Err != nil {
			o.Err = errors.LackParameterError("username")
		}
	}
}

func WithUserExt(o *GetUserOptions) {
	o.Ext = true
}

func WithoutUserExt(o *GetUserOptions) {
	o.Ext = false
}

func WithUserId(id string) WithGetUserOptions {
	return func(o *GetUserOptions) {
		o.Id = id
		if len(o.Id) == 0 && o.Err != nil {
			o.Err = errors.LackParameterError("id")
		}
	}
}

func WithApps(o *GetUserOptions) {
	o.Apps = true
}

func WithEmail(email string) WithGetUserOptions {
	return func(o *GetUserOptions) {
		o.Email = email
		if len(o.Email) == 0 && o.Err != nil {
			o.Err = errors.LackParameterError("email")
		}
	}
}

func NewGetUserOptions(opts ...WithGetUserOptions) *GetUserOptions {
	o := &GetUserOptions{}
	for _, option := range opts {
		option(o)
	}
	return o
}
