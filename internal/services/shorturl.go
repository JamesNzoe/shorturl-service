package services

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/fpay/foundation-go/cache"
	"github.com/fpay/gopress"
	shorturl "github.com/fpay/lehuipay-shorturl-go"
)

var _ shorturl.ShortURLService = (*ShortURLService)(nil)

const (
	AutoIncrementKey = "short:url:increment"
	SourceURLPreKey  = "source:url"
	ShortURLPreKey   = "short:url"
)

type ShortURLService struct {
	cache  *cache.RedisCache
	domain string
}

func NewShortURLService(ca *cache.RedisCache, domain string) *ShortURLService {
	return &ShortURLService{
		cache:  ca,
		domain: domain,
	}
}

// CreateShortURL 创建短URL
func (s *ShortURLService) CreateShortURL(ctx context.Context, url string, customName string) (shorturl string, err error) {
	urlKey := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	hashid, err := s.cache.Redis().HGet(SourceURLPreKey, urlKey).Result()
	if err == nil && hashid != "" {
		return s.domain + hashid, nil
	}

	id, err := s.cache.IncrBy(AutoIncrementKey, 1)
	if err != nil {
		return "", err
	}
	hasher, err := NewHasher()
	if err != nil {
		return "", err
	}
	if customName != "" {
		hashid = customName
	} else {
		hashid, err = hasher.EncodeInt64([]int64{id})
	}
	if err != nil {
		return "", err
	}
	_, err = s.cache.Redis().HSet(ShortURLPreKey, hashid, url).Result()
	_, err = s.cache.Redis().HSet(SourceURLPreKey, urlKey, hashid).Result()
	if err != nil {
		return "", err
	}
	return s.domain + hashid, nil
}

// ExtractURL 析出原始的URL
func (s *ShortURLService) ExtractURL(ctx gopress.Context, hashid string) (url string, err error) {
	srcUrl, err := s.cache.Redis().HGet(ShortURLPreKey, hashid).Result()
	if err != nil {
		return "", err
	}

	return srcUrl, nil
}
