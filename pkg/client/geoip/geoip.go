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

package geoip

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/oschwald/geoip2-golang"

	"github.com/MicroOps-cn/idas/api"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

type Client struct {
	c *geoip2.Reader
	o *GeoIPOptions
}

func (c Client) City(ip net.IP) (*geoip2.City, error) {
	return c.c.City(ip)
}

func NewClient(ctx context.Context, options *GeoIPOptions) (db *geoip2.Reader, err error) {
	logger := log.GetContextLogger(ctx)
	level.Debug(logger).Log("msg", "open geoip database", "file", options.Path)
	db, err = geoip2.Open(options.Path)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("failed to open GeoIP database: %s", options.Path), "err", err)
		return nil, errors.WithServerError(http.StatusInternalServerError, err, fmt.Sprintf("failed to open GeoIP database: %s", options.Path))
	}

	stopCh := signals.SetupSignalHandler(logger)
	stopCh.Add(1)
	go func() {
		<-stopCh.Channel()
		stopCh.WaitRequest()

		if err = db.Close(); err != nil {
			level.Warn(logger).Log("msg", "Failed to close GeoIP database", "err", err)
		} else {
			level.Debug(logger).Log("msg", "GeoIP database closed")
		}
		stopCh.Done()
	}()

	return db, nil
}

func NewGeoIPOptions() *GeoIPOptions {
	return &GeoIPOptions{}
}

var _ api.CustomType = &Client{}

// Merge implement proto.Merger
func (c *Client) Merge(src proto.Message) {
	if s, ok := src.(*Client); ok {
		c.o = s.o
		c.c = s.c
	}
}

// Reset *implement proto.Message*
func (c *Client) Reset() {
	c.o.Reset()
}

// String implement proto.Message
func (c Client) String() string {
	return c.o.String()
}

// ProtoMessage implement proto.Message
func (c *Client) ProtoMessage() {
	c.o.ProtoMessage()
}

func (c Client) Marshal() ([]byte, error) {
	return proto.Marshal(c.o)
}

func (c *Client) Unmarshal(data []byte) (err error) {
	if c.o == nil {
		c.o = NewGeoIPOptions()
	}
	if err = proto.Unmarshal(data, c.o); err != nil {
		return err
	}
	if c.c, err = NewClient(context.Background(), c.o); err != nil {
		return err
	}
	return
}

var _ proto.Unmarshaler = &Client{}

func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.o)
}

func (c *Client) UnmarshalJSON(data []byte) (err error) {
	if c.o == nil {
		c.o = NewGeoIPOptions()
	}
	if err = json.Unmarshal(data, c.o); err != nil {
		return err
	}
	if c.c, err = NewClient(context.Background(), c.o); err != nil {
		return err
	}
	return
}
