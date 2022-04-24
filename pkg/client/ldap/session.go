package ldap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/go-kit/log/level"
	"github.com/go-ldap/ldap"
	"idas/pkg/logs"
	"idas/pkg/utils/wrapper"
	"time"
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
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
	level.Debug(logger).Log("msg", "create ldap object", "dn", addRequest.DN, "attributes", string(wrapper.Must[[]byte](json.Marshal(addRequest.Attributes))))
	return s.c.Add(addRequest)
}

func (s *Session) Del(delRequest *ldap.DelRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
	level.Debug(logger).Log("msg", "delete ldap object", "dn", delRequest.DN)
	return s.c.Del(delRequest)
}

func (s *Session) Modify(modifyRequest *ldap.ModifyRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
	level.Debug(logger).Log("msg", "modify ldap object", "dn", modifyRequest.DN, "attributes", string(wrapper.Must[[]byte](json.Marshal(modifyRequest.Changes))))
	return s.c.Modify(modifyRequest)
}

func (s *Session) ModifyDN(modifyDNRequest *ldap.ModifyDNRequest) error {
	if s.err != nil {
		return s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
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
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
	level.Debug(logger).Log("msg", "modify ldap user password", "uid", passwordModifyRequest.UserIdentity)
	return s.c.PasswordModify(passwordModifyRequest)
}

func (s *Session) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	logger := logs.GetContextLogger(s.ctx, logs.Caller(4))
	level.Debug(logger).Log("msg", "ldap search", "baseDN", searchRequest.BaseDN, "filter", searchRequest.Filter)
	return s.c.Search(searchRequest)
}

func (s *Session) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.c.SearchWithPaging(searchRequest, pagingSize)
}

var _ ldap.Client = &Session{}
