package ton

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	endpoint string
	auth     string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

func (c *Client) SetAuth(auth string) {
	c.auth = auth
}

func (s *Client) GetRequest(ctx context.Context, url string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if s.auth != "" {
		req.Header.Add("X-API-Key", s.auth)
	}
	// http client and send request
	// 跳过证书验证
	httpclient := NoKeepAliveTransport(120 * time.Second)
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}
	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", res.StatusCode)
	}
	return nil
}

func (s *Client) PostRequest(ctx context.Context, url string, params interface{}, response interface{}) error {
	// post request
	//fmt.Printf("params: %+v\n", params)
	j, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if s.auth != "" {
		req.Header.Add("X-API-Key", s.auth)
	}
	// http client and send request
	//跳过证书验证
	httpclient := NoKeepAliveTransport(120 * time.Second)
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}
	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", res.StatusCode)
	}
	return nil
}

// http短连接设置
func NoKeepAliveTransport(timeout time.Duration) *http.Client {
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return &client
}
