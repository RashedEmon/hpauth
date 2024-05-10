package hpauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetHotelPlannerAuthToken(apiKey string, secretKey string, accountId string) string {
	unixEpoch := strconv.FormatInt(time.Now().Unix(), 10)
	encodedAPIKey := base64.URLEncoding.EncodeToString([]byte(apiKey))
	signatureKey := fmt.Sprintf("%s|%s|%s", encodedAPIKey, accountId, unixEpoch)
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(signatureKey))
	hashValue := hash.Sum(nil)
	encodedHashValue := base64.URLEncoding.EncodeToString(hashValue)
	authToken := fmt.Sprintf("%s.%s", encodedAPIKey, encodedHashValue)
	return authToken
}

// Config the plugin configuration.
type Config struct {
	Headers  map[string]string
	HpConfig map[string]string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers:  make(map[string]string),
		HpConfig: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	hpconfig map[string]string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &Demo{
		headers:  config.Headers,
		next:     next,
		hpconfig: config.HpConfig,
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for key, value := range a.headers {
		req.Header.Set(key, value)
	}
	apiKey := a.hpconfig["apiKey"]
	secretKey := a.hpconfig["secretKey"]
	accountId := a.hpconfig["accountId"]

	req.Header.Set("Authorization", GetHotelPlannerAuthToken(apiKey, secretKey, accountId))

	a.next.ServeHTTP(rw, req)
}
