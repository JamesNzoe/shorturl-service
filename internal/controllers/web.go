package controllers

import (
	"net/http"

	"github.com/fpay/gopress"
	shorturl "github.com/fpay/lehuipay-shorturl-go"
)

type WEBOptions struct {
	Port int `yaml:"port" mapstructure:"port"`
}

type Web struct {
	shorturlService shorturl.ShortURLService
}

func NewWeb(s shorturl.ShortURLService) *Web {
	return &Web{
		shorturlService: s,
	}
}

func (w *Web) Index(ctx gopress.Context) error {
	hash := ctx.Param("hash")
	url, err := w.shorturlService.ExtractURL(ctx, hash)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Unknow short url.")
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}
