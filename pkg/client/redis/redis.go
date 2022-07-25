package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"

	"idas/pkg/logs"
	"idas/pkg/utils/signals"
)

type Client struct {
	client  *redis.Client
	options *RedisOptions
}

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
