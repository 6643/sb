package sbv2

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	rt "sb/go/sb/v2"
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
	defaultMaxRespBytes int64 = 4 * 1024 * 1024
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Timeout time.Duration
	Retries int
	MaxRespBytes int64
	headers map[string]string
	headersMu sync.RWMutex
}

func NewClient(baseURL string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
		Timeout: 5 * time.Second,
		Retries: 3,
		MaxRespBytes: defaultMaxRespBytes,
		headers: make(map[string]string),
	}
}

func SetClientHeader(c *Client, key, value string) {
	if c == nil { return }
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
	if c.headers == nil { c.headers = make(map[string]string) }
	c.headers[key] = value
}
func GetClientHeader(c *Client, key string) string {
	if c == nil { return "" }
	c.headersMu.RLock()
	defer c.headersMu.RUnlock()
	return c.headers[key]
}
func RemoveClientHeader(c *Client, key string) {
	if c == nil { return }
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
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
	if ctx == nil { ctx = context.Background() }
	httpClient := c.HTTP
	if httpClient == nil { httpClient = &http.Client{} }
	maxRetries := c.Retries
	if maxRetries < 0 { maxRetries = 0 }
	maxRespBytes := c.MaxRespBytes
	if maxRespBytes <= 0 { maxRespBytes = defaultMaxRespBytes }
	var resp *http.Response
	var err error

	for i := 0; i <= maxRetries; i++ {
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
		c.headersMu.RLock()
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
		c.headersMu.RUnlock()
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/octet-stream")
		}
		resp, err = httpClient.Do(req)
		cancel()
		if err != nil {
			if isTimeout(err) && i < maxRetries { continue }
			if isTimeout(err) { return nil, RpcTimeout }
			return nil, RpcNoConn
		}
		if resp.StatusCode == http.StatusRequestTimeout && i < maxRetries {
			resp.Body.Close()
			continue
		}
		break
	}
	if resp == nil { return nil, RpcNoConn }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return nil, RpcErrCode(resp.StatusCode) }
	readLimit := maxRespBytes
	if readLimit < (1<<63 - 1) { readLimit++ }
	b, readErr := io.ReadAll(io.LimitReader(resp.Body, readLimit))
	if readErr != nil { return nil, RpcRespErr }
	if readLimit > maxRespBytes && int64(len(b)) > maxRespBytes { return nil, RpcRespErr }
	return b, RpcOk
}

// CallUserGetAbc 获取用户的id
func CallUserGetAbc(c *Client, ctx context.Context) (result OrderStatus, errCode RpcErrCode) {
	var buf bytes.Buffer

	body, status := doClient(c, ctx, "/user.get_abc", buf.Bytes())
	if status != RpcOk {
		return result, status
	}
	respBuf := bytes.NewBuffer(body)
	{
		value, err := rt.GetU8(respBuf)
		if err != nil { return result, RpcRespErr }
		item := OrderStatus(value)
		if !isOrderStatus(item) { return result, RpcRespErr }
		result = item
	}
	if respBuf.Len() != 0 { return result, RpcRespErr }
	return result, status
}

// CallUserGetAbcd 获取abcd
func CallUserGetAbcd(c *Client, ctx context.Context, page uint8, size uint8) (result OrderStatus, errCode RpcErrCode) {
	var buf bytes.Buffer
	{
		if err := rt.SetU8(&buf, page); err != nil { return result, RpcReqErr }
	}
	{
		if err := rt.SetU8(&buf, size); err != nil { return result, RpcReqErr }
	}

	body, status := doClient(c, ctx, "/user.get_abcd", buf.Bytes())
	if status != RpcOk {
		return result, status
	}
	respBuf := bytes.NewBuffer(body)
	{
		value, err := rt.GetU8(respBuf)
		if err != nil { return result, RpcRespErr }
		item := OrderStatus(value)
		if !isOrderStatus(item) { return result, RpcRespErr }
		result = item
	}
	if respBuf.Len() != 0 { return result, RpcRespErr }
	return result, status
}

// CallUserSetSimInfo 设置sim信息
// 无返回值
func CallUserSetSimInfo(c *Client, ctx context.Context, info *SimInfo) (errCode RpcErrCode) {
	var buf bytes.Buffer
	{
		if err := SetSimInfo(&buf, info); err != nil { return RpcReqErr }
	}

	body, status := doClient(c, ctx, "/user.set_sim_info", buf.Bytes())
	if status != RpcOk {
		return status
	}
	if len(body) != 0 { return RpcRespErr }
	return status
}

// CallGetCount 获取数量
func CallGetCount(c *Client, ctx context.Context, page uint8) (result uint8, errCode RpcErrCode) {
	var buf bytes.Buffer
	{
		if err := rt.SetU8(&buf, page); err != nil { return result, RpcReqErr }
	}

	body, status := doClient(c, ctx, "/get_count", buf.Bytes())
	if status != RpcOk {
		return result, status
	}
	respBuf := bytes.NewBuffer(body)
	{
		value, err := rt.GetU8(respBuf)
		if err != nil { return result, RpcRespErr }
		result = value
	}
	if respBuf.Len() != 0 { return result, RpcRespErr }
	return result, status
}

// CallGetBin 获取bin
func CallGetBin(c *Client, ctx context.Context, page uint8) (result []byte, errCode RpcErrCode) {
	var buf bytes.Buffer
	{
		if err := rt.SetU8(&buf, page); err != nil { return result, RpcReqErr }
	}

	body, status := doClient(c, ctx, "/get_bin", buf.Bytes())
	if status != RpcOk {
		return result, status
	}
	respBuf := bytes.NewBuffer(body)
	{
		state, err := rt.GetU8(respBuf)
		if err != nil { return result, RpcRespErr }
		value, err := rt.GetBinCompactInto(respBuf, state, nil)
		if err != nil { return result, RpcRespErr }
		result = value
	}
	if respBuf.Len() != 0 { return result, RpcRespErr }
	return result, status
}
