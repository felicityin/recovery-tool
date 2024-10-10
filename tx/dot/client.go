package dot

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"recovery-tool/common"
	"time"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

func (c *Client) Get(ctx context.Context, path string, queryValues url.Values, response interface{}) error {
	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	if queryValues != nil {
		values := uri.Query()
		for k, v := range values {
			queryValues[k] = v
		}
		uri.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)

	if err != nil {
		return err
	}

	// http client and send request
	httpclient := &http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := httpclient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}

	// return result
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Request(ctx context.Context, method string, params interface{}, response interface{}) error {
	// post data
	j, err := json.Marshal(params)
	if err != nil {
		common.Logger.Errorf("marshal params err: %s", err.Error())
		return err
	}

	// post request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%v/%v", c.endpoint, method), bytes.NewBuffer(j))
	if err != nil {
		common.Logger.Errorf("post request err: %s", err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	httpclient := NoKeepAliveTransport(120 * time.Second)
	res, err := httpclient.Do(req)
	if err != nil {
		common.Logger.Errorf("request err: %s", err.Error())
		return err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.Logger.Errorf("parse body err: %s", err.Error())
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			common.Logger.Errorf("unmarshal err: %s", err.Error())
			return err
		}
	}

	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		common.Logger.Errorf("get status code: %d", res.StatusCode)
		return fmt.Errorf("get status code: %d", res.StatusCode)
	}
	return nil
}

func CheckRpcResult(res GeneralResponse, err error) error {
	if res.Error != "" {
		common.Logger.Errorf("res.err: %s", res.Error)
		common.Logger.Errorf("cause: %s", res.Cause)
		return fmt.Errorf("%s", res.Cause)
	}

	if res.Message != "" {
		common.Logger.Errorf("msg: %s", res.Message)
		return fmt.Errorf(res.Message)
	}

	if err != nil {
		common.Logger.Errorf("err: %s %s", err.Error(), res.Cause)
		return err
	}

	return nil
}

func NoKeepAliveTransport(timeout time.Duration) *http.Client {
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
