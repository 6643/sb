package sb

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type RpcErrCode int

const (
	RpcOk        RpcErrCode = 200
	RpcNoConn    RpcErrCode = 0
	RpcTimeout   RpcErrCode = 408
	RpcReqErr    RpcErrCode = 400
	RpcRespErr   RpcErrCode = 500
	RpcNotAuth   RpcErrCode = 401
	RpcNotExist  RpcErrCode = 404
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Timeout time.Duration
	Retries int
	headers map[string]string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 5 * time.Second},
		Timeout: 5 * time.Second,
		Retries: 3,
		headers: make(map[string]string),
	}
}

func (c *Client) SetHeader(key, value string) { c.headers[key] = value }
func (c *Client) GetHeader(key string) string  { return c.headers[key] }
func (c *Client) RemoveHeader(key string)      { delete(c.headers, key) }

func (c *Client) SetAuthorization(token string) { c.SetHeader("Authorization", "Bearer "+token) }
func (c *Client) GetAuthorization() string      { return c.GetHeader("Authorization") }
func (c *Client) RemoveAuthorization()          { c.RemoveHeader("Authorization") }
func (c *Client) IsAuthorized() bool            { return c.GetAuthorization() != "" }

func isTimeout(err error) bool {
	if err == nil { return false }
	if err == context.DeadlineExceeded { return true }
	if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() { return true }
	return false
}

func (c *Client) do(ctx context.Context, path string, body []byte) ([]byte, RpcErrCode) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.Retries; i++ {
		if i > 0 {
			timer := time.NewTimer(time.Duration(i) * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, RpcTimeout
			case <-timer.C:
			}
		}

		req, reqErr := http.NewRequestWithContext(ctx, "POST", c.BaseURL+path, bytes.NewReader(body))
		if reqErr != nil {
			return nil, RpcNoConn
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		resp, err = c.HTTP.Do(req)
		if err != nil {
			if isTimeout(err) && i < c.Retries {
				continue
			}
			if isTimeout(err) {
				return nil, RpcTimeout
			}
			return nil, RpcNoConn
		}

		if resp.StatusCode == http.StatusRequestTimeout && i < c.Retries {
			resp.Body.Close()
			continue
		}
		break
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, RpcErrCode(resp.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)
	return b, RpcOk
}

// UserGetAbc 获取用户的id
func (c *Client) UserGetAbc(ctx context.Context) (result OrderStatus, errCode RpcErrCode) {
	var res U8
	var buf bytes.Buffer

	body, status := c.do(ctx, "/user.get_abc", buf.Bytes())
	if status != RpcOk {
		return OrderStatus(res), status
	}

	if err := GetAll(bytes.NewBuffer(body), (*U8)(&res)); err != nil {
		return OrderStatus(res), RpcRespErr
	}
	return OrderStatus(res), status
}
// UserGetAbcd 获取abcd
func (c *Client) UserGetAbcd(ctx context.Context, page uint8, size uint8) (result OrderStatus, errCode RpcErrCode) {
	var res U8
	var buf bytes.Buffer
	if err := SetAll(&buf, U8(page), U8(size)); err != nil {
		return OrderStatus(res), RpcReqErr
	}

	body, status := c.do(ctx, "/user.get_abcd", buf.Bytes())
	if status != RpcOk {
		return OrderStatus(res), status
	}

	if err := GetAll(bytes.NewBuffer(body), (*U8)(&res)); err != nil {
		return OrderStatus(res), RpcRespErr
	}
	return OrderStatus(res), status
}
// UserSetSimInfo 设置sim信息
func (c *Client) UserSetSimInfo(ctx context.Context, info *SimInfo) (errCode RpcErrCode) {
	
	var buf bytes.Buffer
	if err := SetAll(&buf, info); err != nil {
		return RpcReqErr
	}

	_, status := c.do(ctx, "/user.set_sim_info", buf.Bytes())
	if status != RpcOk {
		return status
	}

	return status
}
// GetCount 获取数量
func (c *Client) GetCount(ctx context.Context, page uint8) (result uint8, errCode RpcErrCode) {
	var res U8
	var buf bytes.Buffer
	if err := SetAll(&buf, U8(page)); err != nil {
		return uint8(res), RpcReqErr
	}

	body, status := c.do(ctx, "/get_count", buf.Bytes())
	if status != RpcOk {
		return uint8(res), status
	}

	if err := GetAll(bytes.NewBuffer(body), &res); err != nil {
		return uint8(res), RpcRespErr
	}
	return uint8(res), status
}
// GetBin 获取bin
func (c *Client) GetBin(ctx context.Context, page uint8) (result []byte, errCode RpcErrCode) {
	var res Bin
	var buf bytes.Buffer
	if err := SetAll(&buf, U8(page)); err != nil {
		return []byte(res), RpcReqErr
	}

	body, status := c.do(ctx, "/get_bin", buf.Bytes())
	if status != RpcOk {
		return []byte(res), status
	}

	if err := GetAll(bytes.NewBuffer(body), &res); err != nil {
		return []byte(res), RpcRespErr
	}
	return []byte(res), status
}
