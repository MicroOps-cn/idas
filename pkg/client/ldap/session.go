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

package ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-ldap/ldap"

	"github.com/MicroOps-cn/idas/pkg/logs"
)

func checkStatus(ctx context.Context, name string, req interface{}, err error) {
	logger := log.GetContextLogger(ctx, log.WithCaller(6))
	typeOf := getTypeElem(reflect.TypeOf(req))
	valOf := getValElem(reflect.ValueOf(req))
	for i := 0; i < valOf.NumField(); i++ {
		field := valOf.Field(i)
		elem := getValElem(field)
		switch elem.Kind() {
		case reflect.Slice, reflect.Array:
			logger = kitlog.With(logger, logs.WrapKeyName(typeOf.Field(i).Name), w.JSONStringer(elem.Interface()))
		default:
			fieldName := typeOf.Field(i).Name
			if fieldName == "Controls" && elem.Interface() == nil {
				continue
			}
			logger = kitlog.With(logger, logs.WrapKeyName(fieldName), elem.Interface())
		}
	}

	if err != nil {
		level.Error(logger).Log("msg", "failed to execute: "+name, "err", err)
	} else {
		level.Debug(logger).Log("msg", "execute:  "+name)
	}
}

func getValElem(valOf reflect.Value) reflect.Value {
	if valOf.Kind() == reflect.Ptr || valOf.Kind() == reflect.Pointer {
		return getValElem(valOf.Elem())
	}
	return valOf
}

func getTypeElem(typeOf reflect.Type) reflect.Type {
	if typeOf.Kind() == reflect.Ptr || typeOf.Kind() == reflect.Pointer {
		return getTypeElem(typeOf.Elem())
	}
	return typeOf
}

type Session struct {
	c   ldap.Client
	err error
	ctx context.Context
}

func (s *Session) Error() error {
	return s.err
}

func (s *Session) Start() {
	if s.c != nil {
		s.c.Start()
	}
}

func (s *Session) StartTLS(config *tls.Config) error {
	if s.err != nil {
		return s.err
	}
	return s.c.StartTLS(config)
}

func (s *Session) Close() {
	if s.c != nil {
		s.c.Close()
	}
}

func (s *Session) SetTimeout(duration time.Duration) {
	if s.c != nil {
		s.c.SetTimeout(duration)
	}
}

type bindRequest struct {
	Username string
	Password string
}

func (s *Session) Bind(username, password string) (err error) {
	defer func() {
		checkStatus(s.ctx, "Bind", bindRequest{Username: username, Password: "<secret>"}, err)
	}()
	if s.err != nil {
		return s.err
	}
	return s.c.Bind(username, password)
}

func (s *Session) SimpleBind(simpleBindRequest *ldap.SimpleBindRequest) (ret *ldap.SimpleBindResult, err error) {
	defer func() {
		var req ldap.SimpleBindRequest
		if simpleBindRequest != nil {
			req = *simpleBindRequest
		}
		req.Password = "<secret>"
		checkStatus(s.ctx, "Simple Bind", req, err)
	}()
	if s.err != nil {
		return nil, s.err
	}
	return s.c.SimpleBind(simpleBindRequest)
}

func (s *Session) Add(addRequest *ldap.AddRequest) (err error) {
	defer func() {
		checkStatus(s.ctx, "add LDAP entry", addRequest, err)
	}()
	if s.err != nil {
		return s.err
	}
	return s.c.Add(addRequest)
}

func (s *Session) Del(delRequest *ldap.DelRequest) (err error) {
	defer func() {
		checkStatus(s.ctx, "del LDAP entry", delRequest, err)
	}()
	if s.err != nil {
		return s.err
	}
	return s.c.Del(delRequest)
}

func (s *Session) Modify(modifyRequest *ldap.ModifyRequest) (err error) {
	defer func() {
		checkStatus(s.ctx, "modify LDAP entry", modifyRequest, err)
	}()
	if s.err != nil {
		return s.err
	}
	return s.c.Modify(modifyRequest)
}

func (s *Session) ModifyDN(modifyDNRequest *ldap.ModifyDNRequest) (err error) {
	defer func() {
		checkStatus(s.ctx, "modify LDAP entry DN", modifyDNRequest, err)
	}()
	if s.err != nil {
		return s.err
	}
	return s.c.ModifyDN(modifyDNRequest)
}

func (s *Session) Compare(dn, attribute, value string) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.c.Compare(dn, attribute, value)
}

func (s *Session) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (ret *ldap.PasswordModifyResult, err error) {
	defer func() {
		req := ldap.PasswordModifyRequest{UserIdentity: passwordModifyRequest.UserIdentity, OldPassword: "<secret>", NewPassword: "<secret>"}
		checkStatus(s.ctx, "Bind", req, err)
	}()
	if s.err != nil {
		return nil, s.err
	}
	logger := log.GetContextLogger(s.ctx, log.WithCaller(4))
	level.Debug(logger).Log("msg", "modify ldap user password", "uid", passwordModifyRequest.UserIdentity)
	return s.c.PasswordModify(passwordModifyRequest)
}

func (s *Session) Search(searchRequest *ldap.SearchRequest) (result *ldap.SearchResult, err error) {
	defer func() {
		logger := log.GetContextLogger(s.ctx, log.WithCaller(5))
		logger = kitlog.With(logger, logs.WrapKeyName("baseDN"), searchRequest.BaseDN,
			logs.WrapKeyName("filter"), searchRequest.Filter,
			logs.WrapKeyName("attributes"), w.JSONStringer(searchRequest.Attributes),
			logs.WrapKeyName("scope"), ldap.ScopeMap[searchRequest.Scope],
			logs.WrapKeyName("derefAliases"), ldap.DerefMap[searchRequest.DerefAliases],
			logs.WrapKeyName("limits"), fmt.Sprintf("sizeLimit=%d&timeLimit=%d", searchRequest.SizeLimit, searchRequest.TimeLimit),
		)

		if len(searchRequest.Controls) != 0 {
			logger = kitlog.With(logger, logs.WrapKeyName("controls"), searchRequest.Controls)
		}
		if result != nil {
			logger = kitlog.With(logger, logs.WrapKeyName("resultCount"), len(result.Entries))
		}

		if err == nil || IsLdapError(err, 32) {
			if err != nil {
				logger = kitlog.With(logger, "err", err)
			}
			level.Debug(logger).Log("msg", "execute: ldap search")
		} else {
			level.Error(logger).Log("msg", "failed to execute ldap search", "err", err)
		}
	}()
	if s.err != nil {
		return nil, s.err
	}
	return s.c.Search(searchRequest)
}

func (s *Session) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (ret *ldap.SearchResult, err error) {
	defer func() {
		checkStatus(s.ctx, "modify LDAP entry DN", searchRequest, err)
	}()
	if s.err != nil {
		return nil, s.err
	}
	return s.c.SearchWithPaging(searchRequest, pagingSize)
}

var _ ldap.Client = &Session{}
