package {{.Package}}

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
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
		headers := make(map[string]string, len(c.headers))
		for k, v := range c.headers {
			headers[k] = v
		}
		c.headersMu.RUnlock()
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/octet-stream")
		}

		resp, err = httpClient.Do(req)
		cancel()
		if err != nil {
			if isTimeout(err) && i < maxRetries {
				continue
			}
			if isTimeout(err) {
				return nil, RpcTimeout
			}
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

	if resp.StatusCode != http.StatusOK {
		return nil, RpcErrCode(resp.StatusCode)
	}

	readLimit := maxRespBytes
	if readLimit < (1<<63 - 1) { readLimit++ }
	b, readErr := io.ReadAll(io.LimitReader(resp.Body, readLimit))
	if readErr != nil {
		return nil, RpcRespErr
	}
	if readLimit > maxRespBytes && int64(len(b)) > maxRespBytes { return nil, RpcRespErr }
	return b, RpcOk
}

{{range .Apis}}
{{- $resData := .Result -}}
// Call{{.Name | PascalCase}} {{.Note}}
func Call{{.Name | PascalCase}}(c *Client, ctx context.Context{{range .Args}}, {{.Name | CamelCase}} {{GoLogicType .Type}}{{end}}) ({{if eq $resData.Name "nil"}}errCode RpcErrCode{{else}}result {{GoLogicType .Result}}, errCode RpcErrCode{{end}}) {
	var buf bytes.Buffer
	{{- if .Args}}
	{{- range .Args}}
	{{- if and (IsStruct .Type) (not .Type.IsList)}}
	if {{.Name | CamelCase}} == nil {
		return {{if eq $resData.Name "nil"}}RpcReqErr{{else}}result, RpcReqErr{{end}}
	}
	{{- end}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	if !Is{{.Type.Name | PascalCase}}List(({{.Type.Name | PascalCase}}List)({{.Name | CamelCase}})) {
		return {{if eq $resData.Name "nil"}}RpcReqErr{{else}}result, RpcReqErr{{end}}
	}
	{{- else}}
	if !Is{{.Type.Name | PascalCase}}({{.Name | CamelCase}}) {
		return {{if eq $resData.Name "nil"}}RpcReqErr{{else}}result, RpcReqErr{{end}}
	}
	{{- end}}
	{{- end}}
	{{- end}}
	if err := SetAll(&buf{{range .Args}}, func(buf *bytes.Buffer) error {
		{{- if IsBaseType .Type}}
		return Set{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf, {{.Name | CamelCase}})
		{{- else if IsEnum .Type}}
		return Set{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf, {{if .Type.IsList}}({{.Type.Name | PascalCase}}List)({{.Name | CamelCase}}){{else}}{{.Name | CamelCase}}{{end}})
		{{- else if IsStruct .Type}}
		return Set{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf, {{if .Type.IsList}}({{.Type.Name | PascalCase}}List)({{.Name | CamelCase}}){{else}}{{.Name | CamelCase}}{{end}})
		{{- end}}
	}{{end}}); err != nil {
		return {{if eq $resData.Name "nil"}}RpcReqErr{{else}}result, RpcReqErr{{end}}
	}
	{{- end}}

	body, status := doClient(c, ctx, "/{{.Name}}", buf.Bytes())
	if status != RpcOk {
		return {{if eq $resData.Name "nil"}}status{{else}}result, status{{end}}
	}

	{{if ne $resData.Name "nil" -}}
	respBuf := bytes.NewBuffer(body)
	if err := GetAll(respBuf, func(buf *bytes.Buffer) error {
		val, err := {{if and (IsStruct .Result) (not .Result.IsList)}}Read{{.Result.Name | PascalCase}}{{else}}Get{{.Result.Name | PascalCase}}{{if .Result.IsList}}List{{end}}{{end}}(buf)
		if err != nil { return err }
		result = val
		return nil
	}); err != nil {
		return result, RpcRespErr
	}
	if respBuf.Len() != 0 { return result, RpcRespErr }
	{{- if and (IsEnum .Result) .Result.IsList}}
	if !Is{{.Result.Name | PascalCase}}List(({{.Result.Name | PascalCase}}List)(result)) { return result, RpcRespErr }
	{{- else if IsEnum .Result}}
	if !Is{{.Result.Name | PascalCase}}(result) { return result, RpcRespErr }
	{{- end}}
	return result, status
	{{- else -}}
	if len(body) != 0 { return RpcRespErr }
	return status
	{{- end}}
}
{{end}}
