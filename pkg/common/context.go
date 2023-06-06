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

package common

import (
	"context"
	"net/url"

	http2 "github.com/MicroOps-cn/fuck/http"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
)

type getWebURLOptions struct {
	subPages []string
	query    url.Values
}

type WithGetWebURLOptions func(*getWebURLOptions)

func WithSubPages(subPages ...string) WithGetWebURLOptions {
	return func(o *getWebURLOptions) {
		o.subPages = subPages
	}
}

func WithQuery(query url.Values) WithGetWebURLOptions {
	return func(o *getWebURLOptions) {
		o.query = query
	}
}

func GetWebURL(ctx context.Context, o ...WithGetWebURLOptions) (string, error) {
	var opts getWebURLOptions
	for _, options := range o {
		options(&opts)
	}

	webPrefix, ok := ctx.Value(global.HTTPWebPrefixKey).(string)
	if !ok {
		return "", errors.NewServerError(500, "webPrefix is null")
	}
	externalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string)
	if !ok {
		return "", errors.NewServerError(500, "externalURL is null")
	}
	extURL, err := url.Parse(externalURL)
	if err != nil {
		return "", errors.NewServerError(500, "")
	}
	p := []string{extURL.Path, webPrefix}
	if len(opts.subPages) > 0 {
		p = append(p, opts.subPages...)
	}
	extURL.Path = http2.JoinPath(p...)
	if len(opts.query) > 0 {
		q := extURL.Query()
		for name, vals := range opts.query {
			for _, val := range vals {
				if len(val) != 0 {
					q.Add(name, val)
				}
			}
		}
		extURL.RawQuery = q.Encode()
	}
	return extURL.String(), nil
}
