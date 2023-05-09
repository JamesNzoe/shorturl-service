package shorturl

import (
	"context"

	"github.com/fpay/gopress"
)

type ShortURLService interface {
	CreateShortURL(ctx context.Context, url string, customName string) (shorturl string, err error)
	ExtractURL(ctx gopress.Context, shorturl string) (url string, err error)
}
