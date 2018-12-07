package itc

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseUrl = "https://api.appstoreconnect.apple.com/"

	defaultRetryCount    = 3
	defaultRetryInterval = time.Second
	defaultUserAgent     = "itc-go/v1"
)

type Options struct {
	RetryCount    int
	RetryInterval time.Duration

	UserAgent string
}

type Option func(*Options)

func RetryCount(count int) Option {
	return func(args *Options) {
		args.RetryCount = count
	}
}

func RetryInterval(interval time.Duration) Option {
	return func(args *Options) {
		args.RetryInterval = interval
	}
}

func UserAgent(ua string) Option {
	return func(args *Options) {
		args.UserAgent = ua
	}
}

type Client interface {
	GetSalesReport(ctx context.Context, data url.Values) (io.Reader, error)
	GetFinanceReport(ctx context.Context, data url.Values) (io.Reader, error)
}

type client struct {
	Options *Options

	client *http.Client
	jwt    *itcJWT
}

func NewClient(keyID, issuerID string, privateKey *ecdsa.PrivateKey, options ...Option) Client {
	args := &Options{
		RetryCount:    defaultRetryCount,
		RetryInterval: defaultRetryInterval,

		UserAgent: defaultUserAgent,
	}

	for _, option := range options {
		option(args)
	}

	return &client{
		jwt: &itcJWT{
			KeyID:      keyID,
			IssuerID:   issuerID,
			PrivateKey: privateKey,
		},
		Options: args,
		client:  http.DefaultClient,
	}
}

func (c *client) makeRequest(ctx context.Context, method string, headers map[string]string,
	pathPart string, data url.Values) (io.Reader, error) {
	sr := new(strings.Reader)

	if data != nil {
		switch method {
		case http.MethodPut:
			fallthrough
		case http.MethodPost:
			sr = strings.NewReader(data.Encode())
		case http.MethodGet:
			pathPart += "?" + data.Encode()
		}
	}

	url := baseUrl + pathPart

	var reader io.Reader

	err := try(c.Options.RetryCount, c.Options.RetryInterval, func() error {
		req, err := c.newRequest(method, url, sr)
		if err != nil {
			return err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := c.client.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			return errors.New(fmt.Sprintf("request failed: url=%s code=%d body=%s",
				req.URL.String(), resp.StatusCode, body))
		}

		reader = resp.Body

		return nil
	})

	return reader, err
}

func (c *client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	token, err := c.jwt.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode jwt: %s", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", c.Options.UserAgent)
	return req, nil
}

func try(count int, interval time.Duration, f func() error) error {
	var err error
	for i := 0; i < count; i++ {
		if err = f(); err == nil {
			return nil
		}
		time.Sleep(interval)
	}
	return err
}
