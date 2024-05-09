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
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
)

type getURLOptions struct {
	subPages  []string
	isWebPage bool
	query     url.Values
	gv        *schema.GroupVersion
	isRoot    bool
}

type WithGetURLOptions func(*getURLOptions)

func WithSubPages(subPages ...string) WithGetURLOptions {
	return func(o *getURLOptions) {
		o.subPages = subPages
	}
}

func WithWebPage(o *getURLOptions) {
	o.isWebPage = true
}

func WithRoot(o *getURLOptions) {
	o.isRoot = true
}

func WithQuery(query url.Values) WithGetURLOptions {
	return func(o *getURLOptions) {
		for name, val := range query {
			o.query[name] = val
		}
	}
}

func WithParam(key, value string) WithGetURLOptions {
	return func(o *getURLOptions) {
		o.query.Set(key, value)
	}
}

func WithAPI(version, group string, subPages ...string) WithGetURLOptions {
	return func(o *getURLOptions) {
		o.gv = &schema.GroupVersion{
			Group:   group,
			Version: version,
		}
		o.subPages = append(o.subPages, subPages...)
	}
}

func GetWebURL(ctx context.Context, o ...WithGetURLOptions) (string, error) {
	return GetURL(ctx, append([]WithGetURLOptions{WithWebPage}, o...)...)
}

func GetURL(ctx context.Context, o ...WithGetURLOptions) (string, error) {
	opts := getURLOptions{query: make(url.Values)}
	for _, options := range o {
		options(&opts)
	}
	externalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string)
	if !ok {
		return "", errors.NewServerError(500, "externalURL is null")
	}
	extURL, err := url.Parse(externalURL)
	if err != nil {
		return "", errors.NewServerError(500, "")
	}
	if opts.isRoot {
		extURL.Path = ""
		extURL.RawQuery = ""
		return extURL.String(), nil
	}
	p := []string{extURL.Path}
	if opts.isWebPage {
		webPrefix, ok := ctx.Value(global.HTTPWebPrefixKey).(string)
		if !ok {
			return "", errors.NewServerError(500, "webPrefix is null")
		}
		p = append(p, webPrefix)
	} else if opts.gv != nil {
		p = append(p, "/api", opts.gv.Version, opts.gv.Group)
	}
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
