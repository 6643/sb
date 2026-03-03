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
		HTTP:    &http.Client{},
		Timeout: 5 * time.Second,
		Retries: 3,
		headers: make(map[string]string),
	}
}

func SetClientHeader(c *Client, key, value string) {
	if c == nil { return }
	if c.headers == nil { c.headers = make(map[string]string) }
	c.headers[key] = value
}
func GetClientHeader(c *Client, key string) string {
	if c == nil { return "" }
	return c.headers[key]
}
func RemoveClientHeader(c *Client, key string) {
	if c == nil { return }
	delete(c.headers, key)
}

func SetClientAuthorization(c *Client, token string) { SetClientHeader(c, "Authorization", "Bearer "+token) }
func GetClientAuthorization(c *Client) string        { return GetClientHeader(c, "Authorization") }
func RemoveClientAuthorization(c *Client)            { RemoveClientHeader(c, "Authorization") }
func IsClientAuthorized(c *Client) bool              { return GetClientAuthorization(c) != "" }

func isTimeout(err error) bool {
	if err == nil { return false }
	if err == context.DeadlineExceeded { return true }
	if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() { return true }
	return false
}

func doClient(c *Client, ctx context.Context, path string, body []byte) ([]byte, RpcErrCode) {
	if c == nil { return nil, RpcNoConn }
	if c.HTTP == nil { c.HTTP = &http.Client{} }
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

		reqCtx := ctx
		cancel := func() {}
		if c.Timeout > 0 {
			reqCtx, cancel = context.WithTimeout(ctx, c.Timeout)
		}
		req, reqErr := http.NewRequestWithContext(reqCtx, "POST", c.BaseURL+path, bytes.NewReader(body))
		if reqErr != nil {
			cancel()
			return nil, RpcNoConn
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		resp, err = c.HTTP.Do(req)
		cancel()
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

	b, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, RpcRespErr
	}
	return b, RpcOk
}

// CallUserGetAbc 获取用户的id
func CallUserGetAbc(c *Client, ctx context.Context) (result OrderStatus, errCode RpcErrCode) {
	var buf bytes.Buffer

	body, status := doClient(c, ctx, "/user.get_abc", buf.Bytes())
	if status != RpcOk {
		return result, status
	}

	if err := GetAll(bytes.NewBuffer(body), func(buf *bytes.Buffer) error {
		val, err := GetOrderStatus(buf)
		if err != nil { return err }
		result = val
		return nil
	}); err != nil {
		return result, RpcRespErr
	}
	return result, status
}
// CallUserGetAbcd 获取abcd
func CallUserGetAbcd(c *Client, ctx context.Context, page uint8, size uint8) (result OrderStatus, errCode RpcErrCode) {
	var buf bytes.Buffer
	if err := SetAll(&buf, func(buf *bytes.Buffer) error {
		return SetU8(buf, page)
	}, func(buf *bytes.Buffer) error {
		return SetU8(buf, size)
	}); err != nil {
		return result, RpcReqErr
	}

	body, status := doClient(c, ctx, "/user.get_abcd", buf.Bytes())
	if status != RpcOk {
		return result, status
	}

	if err := GetAll(bytes.NewBuffer(body), func(buf *bytes.Buffer) error {
		val, err := GetOrderStatus(buf)
		if err != nil { return err }
		result = val
		return nil
	}); err != nil {
		return result, RpcRespErr
	}
	return result, status
}
// CallUserSetSimInfo 设置sim信息
func CallUserSetSimInfo(c *Client, ctx context.Context, info *SimInfo) (errCode RpcErrCode) {
	var buf bytes.Buffer
	if err := SetAll(&buf, func(buf *bytes.Buffer) error {
		return SetSimInfo(buf, info)
	}); err != nil {
		return RpcReqErr
	}

	_, status := doClient(c, ctx, "/user.set_sim_info", buf.Bytes())
	if status != RpcOk {
		return status
	}

	return status
}
// CallGetCount 获取数量
func CallGetCount(c *Client, ctx context.Context, page uint8) (result uint8, errCode RpcErrCode) {
	var buf bytes.Buffer
	if err := SetAll(&buf, func(buf *bytes.Buffer) error {
		return SetU8(buf, page)
	}); err != nil {
		return result, RpcReqErr
	}

	body, status := doClient(c, ctx, "/get_count", buf.Bytes())
	if status != RpcOk {
		return result, status
	}

	if err := GetAll(bytes.NewBuffer(body), func(buf *bytes.Buffer) error {
		val, err := GetU8(buf)
		if err != nil { return err }
		result = val
		return nil
	}); err != nil {
		return result, RpcRespErr
	}
	return result, status
}
// CallGetBin 获取bin
func CallGetBin(c *Client, ctx context.Context, page uint8) (result []byte, errCode RpcErrCode) {
	var buf bytes.Buffer
	if err := SetAll(&buf, func(buf *bytes.Buffer) error {
		return SetU8(buf, page)
	}); err != nil {
		return result, RpcReqErr
	}

	body, status := doClient(c, ctx, "/get_bin", buf.Bytes())
	if status != RpcOk {
		return result, status
	}

	if err := GetAll(bytes.NewBuffer(body), func(buf *bytes.Buffer) error {
		val, err := GetBin(buf)
		if err != nil { return err }
		result = val
		return nil
	}); err != nil {
		return result, RpcRespErr
	}
	return result, status
}

