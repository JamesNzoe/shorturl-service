package controllers

import (
	"context"

	shorturl "github.com/fpay/lehuipay-shorturl-go"
	api "github.com/fpay/lehuipay-shorturl-go/api"
)

type GRPCOptions struct {
	Port int `yaml:"port" mapstructure:"port"`
}

type GrpcServer struct {
	shorturlService shorturl.ShortURLService
}

func NewGrpcServer(s shorturl.ShortURLService) *GrpcServer {
	return &GrpcServer{
		shorturlService: s,
	}
}

func (g *GrpcServer) CreateShortURL(ctx context.Context, req *api.ShortURLRequest) (*api.ShortURLResponse, error) {
	surl, err := g.shorturlService.CreateShortURL(ctx, req.Url, req.CustomName)
	if err != nil {
		return nil, err
	}
	return &api.ShortURLResponse{
		Shorturl: surl,
	}, nil
}
