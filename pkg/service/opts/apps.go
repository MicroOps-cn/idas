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

package opts

type GetAppOptions struct {
	DisableGetUsers,
	DisableGetAccessController,
	DisableGetProxy bool
	Id, Name string
	UserId   []string
	GetTags  bool
}

func NewAppOptions(opts ...WithGetAppOptions) *GetAppOptions {
	o := &GetAppOptions{}
	for _, option := range opts {
		option(o)
	}
	return o
}

type WithGetAppOptions func(o *GetAppOptions)

func WithoutUsers(o *GetAppOptions) {
	o.DisableGetUsers = true
}

func WithGetTags(o *GetAppOptions) {
	o.GetTags = true
}

func WithoutACL(o *GetAppOptions) {
	o.DisableGetAccessController = true
}

func WithUsers(id ...string) func(o *GetAppOptions) {
	return func(o *GetAppOptions) {
		o.DisableGetUsers = false
		o.UserId = id
	}
}

func WithBasic(o *GetAppOptions) {
	WithoutACL(o)
	WithoutProxy(o)
	WithoutUsers(o)
	o.DisableGetAccessController = true
}

func WithoutProxy(o *GetAppOptions) {
	o.DisableGetProxy = true
}

func WithAppId(id string) func(o *GetAppOptions) {
	return func(o *GetAppOptions) {
		o.Id = id
	}
}

func WithAppName(name string) func(o *GetAppOptions) {
	return func(o *GetAppOptions) {
		o.Name = name
	}
}
