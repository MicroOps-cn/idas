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
	"fmt"
	"net"
	"regexp"
	"time"

	fuck_tls "github.com/MicroOps-cn/fuck/clients/tls"
	"github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/sets"
	"github.com/go-kit/log/level"
	"github.com/go-ldap/ldap/v3"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/MicroOps-cn/idas/api"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

func NewLdapPool(ctx context.Context, options *LdapOptions) (pool Pool, err error) {
	logger := log.GetContextLogger(ctx)
	if err = options.Valid(); err != nil {
		return nil, err
	}
	var tlsConfig *tls.Config
	if options.StartTLS || options.IsTLS {
		tlsConfig, err = fuck_tls.NewTLSConfig(options.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ldap tls config: %s", err)
		}
	}
	level.Debug(logger).Log("msg", "connect to ldap server", "host", options.Host, "manager_dn", options.ManagerDn, "isTLS", options.IsTLS, "startTLS", options.StartTLS)
	pool, err = NewChannelPool(ctx, 2, 64, "ldap", func(s string) (c ldap.Client, err error) {
		conn, err := (&net.Dialer{Timeout: ldap.DefaultTimeout}).DialContext(ctx, "tcp", options.Host)
		if err != nil {
			return nil, err
		}
		if options.IsTLS {
			conn = tls.Client(conn, tlsConfig)
		}
		ldapConn := ldap.NewConn(conn, options.IsTLS)
		ldapConn.Start()
		if options.StartTLS {
			if err = ldapConn.StartTLS(tlsConfig); err != nil {
				return nil, err
			}
		}
		return ldapConn, nil
	}, []uint16{ldap.LDAPResultAdminLimitExceeded, ldap.ErrorNetwork})

	if err != nil {
		return nil, err
	}
	client := &Client{
		pool:    pool,
		options: options,
	}
	if err = client.Session(ctx).(*Session).Error(); err != nil {
		return nil, err
	}
	stopCh := signals.SetupSignalHandler(logger)
	stopCh.Add(1)
	go func() {
		<-stopCh.Channel()
		stopCh.WaitRequest()
		if pool != nil {
			pool.Close()
		}
		level.Debug(logger).Log("msg", "LDAP connect closed")
		stopCh.Done()
	}()

	level.Debug(logger).Log("msg", "connected to ldap server: "+options.Host)
	return pool, nil
}

type Client struct {
	pool    Pool
	options *LdapOptions
}

func NewClient(ctx context.Context, o *LdapOptions) (*Client, error) {
	pool, err := NewLdapPool(ctx, o)
	if err != nil {
		return nil, err
	}
	return &Client{pool: pool, options: o}, nil
}

var _ api.CustomType = &Client{}

// Merge implement proto.Merger
func (l *Client) Merge(src proto.Message) {
	if s, ok := src.(*Client); ok {
		l.options = s.options
		l.pool = s.pool
	}
}

// String implement proto.Message
func (l Client) String() string {
	return l.options.String()
}

// ProtoMessage implement proto.Message
func (l *Client) ProtoMessage() {
	l.options.ProtoMessage()
}

// Reset *implement proto.Message*
func (l *Client) Reset() {
	l.options.Reset()
}

func (l Client) Marshal() ([]byte, error) {
	return proto.Marshal(l.options)
}

func (l *Client) Unmarshal(data []byte) (err error) {
	if l.options == nil {
		l.options = NewLdapOptions()
	}
	if l.options.AppObjectClass != "groupOfUniqueNames" && l.options.AppObjectClass != "groupOfNames" {
		return fmt.Errorf("the ldap.app_object_class config can only be groupOfUniqueNames or groupOfNames")
	}
	if err = proto.Unmarshal(data, l.options); err != nil {
		return err
	}
	if l.pool, err = NewLdapPool(context.Background(), l.options); err != nil {
		return err
	}
	return
}

func (l Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.options)
}

func (l *Client) UnmarshalJSON(data []byte) (err error) {
	if l.options == nil {
		l.options = NewLdapOptions()
	}
	if l.options.AppObjectClass != "groupOfUniqueNames" && l.options.AppObjectClass != "groupOfNames" {
		return fmt.Errorf("the ldap.app_object_class config can only be groupOfUniqueNames or groupOfNames")
	}
	if err = json.Unmarshal(data, l.options); err != nil {
		return err
	}
	if l.pool, err = NewLdapPool(context.Background(), l.options); err != nil {
		return err
	}
	return
}

func (l Client) Close() {
	l.pool.Close()
}

func (l *Client) Options() *LdapOptions {
	return l.options
}

func (l *Client) Session(ctx context.Context) ldap.Client {
	if conn := ctx.Value(ldapConnName{}); conn != nil {
		switch db := conn.(type) {
		case ldap.Client:
			return &NopCloser{Client: db}
		default:
			logger := log.GetContextLogger(ctx)
			level.Warn(logger).Log("msg", "Unknown context value type.", "name", fmt.Sprintf("%T", ldapConnName{}), "value", fmt.Sprintf("%T", conn))
		}
	}
	s := &Session{ctx: ctx}
	if l.pool == nil {
		s.err = errors.New("LDAP connection pool not initialized")
		return s
	}
	s.c, s.err = l.pool.Get()
	// cannot connect to ldap server or pool is closed
	if s.err != nil {
		s.err = errors.WithMessage(s.err, "failed to get ldap connection")
		return s
	}
	passwd, err := l.options.ManagerPassword.UnsafeString()
	if err != nil {
		s.err = err
		return s
	}
	s.err = s.c.Bind(l.options.ManagerDn, passwd)
	if s.err != nil {
		s.c.Close()
		s.err = errors.WithMessage(s.err, "failed to connect to LDAP server")
	}
	return s
}

func WithConnContext(ctx context.Context, client ldap.Client) context.Context {
	return context.WithValue(ctx, ldapConnName{}, client)
}

type ldapConnName struct{}

type NopCloser struct {
	ldap.Client
}

func (NopCloser) Close() {}

func init() {
	ldap.DefaultTimeout = 3 * time.Second
}

func IsLdapError(err error, errCode ...uint16) bool {
	if err == nil {
		return false
	}
	ldapErr, ok := err.(*ldap.Error)
	if !ok {
		return false
	} else if len(errCode) == 0 {
		return true
	}
	for _, code := range errCode {
		if ldapErr.ResultCode == code {
			return true
		}
	}
	return false
}

var classNameExp = regexp.MustCompile(`^\( [\d.]+ NAME '(\w+)'`)

func GetAvailableObjectClass(clt ldap.Client) (sets.Set[string], error) {
	subSchemaReq := ldap.NewSearchRequest("", ldap.ScopeBaseObject, ldap.DerefAlways,
		0, 0, false, `(objectClass=*)`,
		[]string{"subschemaSubentry"}, nil,
	)
	subSchemaResp, err := clt.Search(subSchemaReq)
	if err != nil {
		return nil, err
	}
	classNames := sets.New[string]()
	for _, entry := range subSchemaResp.Entries {
		for _, subSchemaSubEntry := range entry.GetAttributeValues("subschemaSubentry") {
			subEntryReq := ldap.NewSearchRequest(subSchemaSubEntry, ldap.ScopeBaseObject, ldap.DerefAlways,
				0, 0, false, `(objectClass=subschema)`,
				[]string{"objectClasses"}, nil,
			)
			subEntryResp, err := clt.Search(subEntryReq)
			if err != nil {
				return nil, err
			}
			for _, subEntry := range subEntryResp.Entries {
				for _, classes := range subEntry.GetAttributeValues("objectClasses") {
					classesMatch := classNameExp.FindStringSubmatch(classes)
					if len(classesMatch) == 2 {
						classNames.Insert(classesMatch[1])
					}
				}
			}
		}
	}
	return classNames, nil
}
