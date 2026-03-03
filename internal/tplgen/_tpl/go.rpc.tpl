package {{.Package}}

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

{{range .Apis}}
{{- $resData := .Result -}}
// {{.Name | PascalCase}} {{.Note}}
func (c *Client) {{.Name | PascalCase}}(ctx context.Context{{range .Args}}, {{.Name | CamelCase}} {{GoLogicType .Type}}{{end}}) ({{if eq $resData.Name "nil"}}errCode RpcErrCode{{else}}result {{GoLogicType .Result}}, errCode RpcErrCode{{end}}) {
	{{if ne $resData.Name "nil"}}var res {{GoRpcType $resData}}{{end}}
	var buf bytes.Buffer
	{{- if .Args}}
	if err := SetAll(&buf{{range .Args}}, {{if or (IsBaseType .Type) (IsEnum .Type)}}{{GoRpcType .Type}}({{.Name | CamelCase}}){{else if IsStruct .Type}}{{.Name | CamelCase}}{{else}}{{.Name | CamelCase}}{{end}}{{end}}); err != nil {
		return {{if eq $resData.Name "nil"}}RpcReqErr{{else}}{{if or (IsBaseType $resData) (IsEnum $resData)}}{{GoLogicType $resData}}(res){{else}}res{{end}}, RpcReqErr{{end}}
	}
	{{- end}}

	{{if eq $resData.Name "nil"}}_{{else}}body{{end}}, status := c.do(ctx, "/{{.Name}}", buf.Bytes())
	if status != RpcOk {
		return {{if eq $resData.Name "nil"}}status{{else}}{{if or (IsBaseType $resData) (IsEnum $resData)}}{{GoLogicType $resData}}(res){{else}}res{{end}}, status{{end}}
	}

	{{if ne $resData.Name "nil" -}}
	if err := GetAll(bytes.NewBuffer(body), {{if or (IsStruct $resData) (IsList $resData)}}&res{{else if IsEnum $resData}}(*U8)(&res){{else}}&res{{end}}); err != nil {
		return {{if or (IsBaseType $resData) (IsEnum $resData)}}{{GoLogicType $resData}}(res){{else}}res{{end}}, RpcRespErr
	}
	return {{if or (IsBaseType $resData) (IsEnum $resData)}}{{GoLogicType $resData}}(res){{else}}res{{end}}, status
	{{- else -}}
	return status
	{{- end}}
}
{{end}}