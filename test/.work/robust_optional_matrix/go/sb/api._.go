package sb

import (
	"bytes"
	"io"
	"net/http"
)

// --- API Handlers ---

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var m *Matrix

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := ReadMatrix(buf)
		if err != nil { return err }
		m = val
		return nil
	}) { return }
	if err := ValidateMatrix(m); err != nil { w.WriteHeader(http.StatusBadRequest); return }

	status := submit(r.Context(), m)
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
}
func FetchHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	result, status := fetch(r.Context())
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetMatrix(buf, result) })
}


// --- 路由注册 ---

type Middleware func(http.HandlerFunc) http.HandlerFunc

func composeMiddleware(mws ...Middleware) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }
		return h
	}
}


func RegisterApi(mux *http.ServeMux, mws ...Middleware) {
	mw := composeMiddleware(mws...)
	mux.HandleFunc("POST /submit", mw(SubmitHandler))
	mux.HandleFunc("POST /fetch", mw(FetchHandler))
}


// --- 内部辅助函数 ---

const defaultMaxReqBytes int64 = 4 * 1024 * 1024

func checkStatus(w http.ResponseWriter, status RpcErrCode) bool {
	if status == RpcOk { return true }
	w.WriteHeader(int(status)); return false
}

func parseRequest(w http.ResponseWriter, r *http.Request, args ...Getter) bool {
	if len(args) == 0 {
		body, err := io.ReadAll(io.LimitReader(r.Body, 1)); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
		if len(body) != 0 { w.WriteHeader(http.StatusBadRequest); return false }
		return true
	}
	readLimit := defaultMaxReqBytes + 1
	body, err := io.ReadAll(io.LimitReader(r.Body, readLimit)); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if int64(len(body)) > defaultMaxReqBytes { w.WriteHeader(http.StatusRequestEntityTooLarge); return false }
	buf := bytes.NewBuffer(body)
	if err := GetAll(buf, args...); err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if buf.Len() != 0 { w.WriteHeader(http.StatusBadRequest); return false }
	return true
}

func sendResponse(w http.ResponseWriter, setter Setter) {
	var buf bytes.Buffer
	if err := SetAll(&buf, setter); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	if _, err := w.Write(buf.Bytes()); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
}
