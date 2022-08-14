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
	"encoding/json"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-ldap/ldap"

	"github.com/MicroOps-cn/idas/pkg/logs"
	w "github.com/MicroOps-cn/idas/pkg/utils/wrapper"
)

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

func (s *Session) Bind(username, password string) error {
	if s.err != nil {
		return s.err
	}
	return s.c.Bind(username, password)
}

func (s *Session) SimpleBind(simpleBindRequest *ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.c.SimpleBind(simpleBindRequest)
}

func (s *Session) Add(addRequest *ldap.AddRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.WithCaller(4))
	level.Debug(logger).Log("msg", "create ldap object", "dn", addRequest.DN, "attributes", string(w.M[[]byte](json.Marshal(addRequest.Attributes))))
	return s.c.Add(addRequest)
}

func (s *Session) Del(delRequest *ldap.DelRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.WithCaller(4))
	level.Debug(logger).Log("msg", "delete ldap object", "dn", delRequest.DN)
	return s.c.Del(delRequest)
}

func (s *Session) Modify(modifyRequest *ldap.ModifyRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.WithCaller(4))
	level.Debug(logger).Log("msg", "modify ldap object", "dn", modifyRequest.DN, "attributes", string(w.M[[]byte](json.Marshal(modifyRequest.Changes))))
	return s.c.Modify(modifyRequest)
}

func (s *Session) ModifyDN(modifyDNRequest *ldap.ModifyDNRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.WithCaller(4))
	level.Debug(logger).Log("msg", "modify ldap object dn", "dn", modifyDNRequest.DN)
	return s.c.ModifyDN(modifyDNRequest)
}

func (s *Session) Compare(dn, attribute, value string) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.c.Compare(dn, attribute, value)
}

func (s *Session) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.WithCaller(4))
	level.Debug(logger).Log("msg", "modify ldap user password", "uid", passwordModifyRequest.UserIdentity)
	return s.c.PasswordModify(passwordModifyRequest)
}

func (s *Session) Search(searchRequest *ldap.SearchRequest) (result *ldap.SearchResult, err error) {
	defer func() {
		logger := logs.GetContextLogger(s.ctx, logs.WithCaller(5))
		logger = log.With(logger, "[baseDN]", searchRequest.BaseDN,
			"[filter]", searchRequest.Filter,
			"[scope]", searchRequest.Scope,
			"[attributes]", searchRequest.Attributes,
			"[derefAliases]", searchRequest.DerefAliases,
			"[sizeLimit]", searchRequest.SizeLimit,
			"[timeLimit]", searchRequest.TimeLimit,
			"[controls]", searchRequest.Controls,
			"[result]", result,
		)
		if err != nil {
			level.Error(logger).Log("msg", "failed to execute ldap search", "err", err)
		} else {
			level.Debug(logger).Log("msg", "ldap search")
		}
	}()
	if s.err != nil {
		return nil, s.err
	}
	return s.c.Search(searchRequest)
}

func (s *Session) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.c.SearchWithPaging(searchRequest, pagingSize)
}

var _ ldap.Client = &Session{}
