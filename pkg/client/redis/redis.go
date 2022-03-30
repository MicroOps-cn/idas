package redis

import (
	"context"
	"time"

	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"

	"idas/config"
	"idas/pkg/logs"
	"idas/pkg/utils/signals"
)

type Client struct {
	client *redis.Client
}

func NewRedisClientOrDie(ctx context.Context, options *config.RedisOptions) *Client {
	client, err := NewRedisClient(ctx, options)
	if err != nil {
		panic(any(err))
	}

	return client
}

func NewRedisClient(ctx context.Context, option *config.RedisOptions) (*Client, error) {
	var r Client
	logger := logs.GetContextLogger(ctx)
	options, err := redis.ParseURL(option.Url)
	if err != nil {
		level.Error(logger).Log("msg", "解析Redis URL失败", "err", err)
		return nil, err
	}

	r.client = redis.NewClient(options)
	if err = r.client.Ping().Err(); err != nil {
		level.Error(logger).Log("msg", "Redis连接失败", "err", err)
		_ = r.client.Close()
		return nil, err
	}

	stopCh := signals.SetupSignalHandler(logger)
	if stopCh != nil {
		stopCh.Add(1)
		go func() {
			<-stopCh.Channel()
			stopCh.WaitRequest()
			if err = r.client.Close(); err != nil {
				level.Error(logger).Log("msg", "Redis客户端关闭出错", "err", err)
				time.Sleep(1 * time.Second)
			}
			level.Error(logger).Log("msg", "关闭Redis连接", "err", err)
			stopCh.Done()
		}()
	}
	return &r, nil
}

func (r *Client) Redis(ctx context.Context) *redis.Client {
	return r.client.WithContext(ctx)
}
