/*

 Copyright 2019 The KubeSphere Authors.

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
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/go-ldap/ldap"
	"github.com/gogo/protobuf/proto"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/utils/signals"
	"net"
	"time"
)

func NewLdapPool(ctx context.Context, options *LdapOptions) (pool Pool, err error) {
	logger := logs.GetContextLogger(ctx)
	if err = options.Valid(); err != nil {
		return nil, err
	}

	level.Debug(logger).Log("msg", "connect to ldap server", "host", options.Host, "manager_dn", options.ManagerDn)
	pool, err = NewChannelPool(ctx, 2, 64, "ldap", func(s string) (c ldap.Client, err error) {
		conn, err := (&net.Dialer{Timeout: ldap.DefaultTimeout}).DialContext(ctx, "tcp", options.Host)
		if err != nil {
			return nil, err
		}
		ldapConn := ldap.NewConn(conn, false)
		ldapConn.Start()
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
		stopCh.Done()
	}()

	level.Debug(logger).Log("msg", "connected to ldap server: "+options.Host)
	return pool, nil
}

type Client struct {
	pool    Pool
	options *LdapOptions
}

// Merge implement proto.Merger
func (c *Client) Merge(src proto.Message) {
	if s, ok := src.(*Client); ok {
		c.options = s.options
		c.pool = s.pool
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

type NopCloser struct {
	ldap.Client
}

func (NopCloser) Close() {}

func (l *Client) Session(ctx context.Context) ldap.Client {
	if conn := ctx.Value(global.LDAPConnName); conn != nil {
		switch db := conn.(type) {
		case ldap.Client:
			return &NopCloser{Client: db}
		default:
			logger := logs.GetContextLogger(ctx)
			level.Warn(logger).Log("msg", "未知的上下文属性(global.LDAPConnName)值", global.LDAPConnName, fmt.Sprintf("%#v", conn))
		}
	}
	s := &Session{ctx: ctx}
	if l.pool == nil {
		s.err = fmt.Errorf("LDAP connection pool not initialized")
		return s
	}
	s.c, s.err = l.pool.Get()
	// cannot connect to ldap server or pool is closed
	if s.err != nil {
		s.err = fmt.Errorf("failed to get ldap connection: %s. ", s.err)
		return s
	}
	s.err = s.c.Bind(l.options.ManagerDn, l.options.ManagerPassword)
	if s.err != nil {
		s.c.Close()
		s.err = fmt.Errorf("failed to connect to LDAP server: %s. ", s.err)
	}
	return s
}

func (l *Client) Options() *LdapOptions {
	return l.options
}

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
