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

package redis

import (
	"context"
	"encoding/json"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"

	"github.com/MicroOps-cn/idas/api"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

type Client struct {
	client  *redis.Client
	options *RedisOptions
}

var _ api.CustomType = &Client{}

// Merge implement proto.Merger
func (r *Client) Merge(src proto.Message) {
	if s, ok := src.(*Client); ok {
		r.options = s.options
		r.client = s.client
	}
}

// String implement proto.Message
func (r Client) String() string {
	return r.options.String()
}

// ProtoMessage implement proto.Message
func (r *Client) ProtoMessage() {
	r.options.ProtoMessage()
}

// Reset *implement proto.Message*
func (r *Client) Reset() {
	r.options.Reset()
}

func (r Client) Marshal() ([]byte, error) {
	return proto.Marshal(r.options)
}

func (r *Client) Unmarshal(data []byte) (err error) {
	if r.options == nil {
		r.options = &RedisOptions{}
	}
	if err = proto.Unmarshal(data, r.options); err != nil {
		return err
	}
	if r.client, err = NewRedisClient(context.Background(), r.options); err != nil {
		return err
	}
	return
}

func (r Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.options)
}

func (r *Client) UnmarshalJSON(data []byte) (err error) {
	if r.options == nil {
		r.options = &RedisOptions{}
	}
	if err = json.Unmarshal(data, r.options); err != nil {
		return err
	}
	if r.client, err = NewRedisClient(context.Background(), r.options); err != nil {
		return err
	}
	return
}

func NewRedisClient(ctx context.Context, option *RedisOptions) (*redis.Client, error) {
	logger := logs.GetContextLogger(ctx)
	options, err := redis.ParseURL(option.Url)
	if err != nil {
		level.Error(logger).Log("msg", "解析Redis URL失败", "err", err)
		return nil, err
	}

	client := redis.NewClient(options)
	if err = client.Ping().Err(); err != nil {
		level.Error(logger).Log("msg", "Redis连接失败", "err", err)
		_ = client.Close()
		return nil, err
	}

	stopCh := signals.SetupSignalHandler(logger)
	if stopCh != nil {
		stopCh.Add(1)
		go func() {
			<-stopCh.Channel()
			stopCh.WaitRequest()
			if err = client.Close(); err != nil {
				level.Error(logger).Log("msg", "Redis客户端关闭出错", "err", err)
				time.Sleep(1 * time.Second)
			}
			level.Error(logger).Log("msg", "关闭Redis连接", "err", err)
			stopCh.Done()
		}()
	}
	return client, nil
}

func (r *Client) Redis(ctx context.Context) *redis.Client {
	return r.client.WithContext(ctx)
}
