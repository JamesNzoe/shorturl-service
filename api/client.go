package api

import (
	"context"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultDialTimeout      = 10
	defaultKeepAliveTime    = 600
	defaultKeepAliveTimeout = 20
)

// Options 客户端连接参数。
type Options struct {
	// Address 服务的地址，IP和端口
	Address string `json:"address" yaml:"address" mapstructure:"address"`

	// DialTimeout 连接超时时间，单位秒
	DialTimeout int64 `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`

	// KeepAliveTime 连接保活周期，单位秒
	KeepAliveTime int64 `json:"keep_alive_time" yaml:"keep_alive_time" mapstructure:"keep_alive_time"`

	// KeepAliveTimeout 发送保活心跳包的超时时间，单位秒
	KeepAliveTimeout int64 `json:"keep_alive_timeout" yaml:"keep_alive_timeout" mapstructure:"keep_alive_timeout"`
}

// loadDefaults 加载参数默认值。
func (opts *Options) loadDefaults() {
	if opts.DialTimeout == 0 {
		opts.DialTimeout = defaultDialTimeout
	}
	if opts.KeepAliveTime == 0 {
		opts.KeepAliveTime = defaultKeepAliveTime
	}
	if opts.KeepAliveTimeout == 0 {
		opts.KeepAliveTimeout = defaultKeepAliveTimeout
	}
}

// buildDialOptions 根据参数构建 gRPC 连接参数。
func (opts *Options) buildDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(opts.KeepAliveTime) * time.Second,
			Timeout:             time.Duration(opts.KeepAliveTimeout) * time.Second,
			PermitWithoutStream: true,
		}),
	}
}

type ShortURLClient interface {
	CreateShortURL(ctx context.Context, url string) (string, error)
	Close() error
}

type Client struct {
	ShortURLServiceClient

	conn *grpc.ClientConn
}

func NewClient(opts Options) (*Client, error) {
	opts.loadDefaults()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.DialTimeout)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, opts.Address, opts.buildDialOptions()...)
	if err != nil {
		return nil, err
	}

	return &Client{
		ShortURLServiceClient: NewShortURLServiceClient(conn),
		conn:                  conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateShortURL(ctx context.Context, url string) (string, error) {
	in := &ShortURLRequest{
		Url: url,
	}
	resp, err := c.ShortURLServiceClient.CreateShortURL(ctx, in)
	if err != nil {
		return "", err
	}
	return resp.GetShorturl(), nil
}
