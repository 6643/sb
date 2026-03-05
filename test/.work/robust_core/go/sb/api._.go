package sb

import (
	"bytes"
	"io"
	"net/http"
)

// --- API Handlers ---

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var env *Envelope

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := ReadEnvelope(buf)
		if err != nil { return err }
		env = val
		return nil
	}) { return }
	if err := ValidateEnvelope(env); err != nil { w.WriteHeader(http.StatusBadRequest); return }

	result, status := echo(r.Context(), env)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetEnvelope(buf, result) })
}
func EchoTailHandler(w http.ResponseWriter, r *http.Request) {
	var env *Envelope

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := ReadEnvelope(buf)
		if err != nil { return err }
		env = val
		return nil
	}) { return }
	if err := ValidateEnvelope(env); err != nil { w.WriteHeader(http.StatusBadRequest); return }

	result, status := echo_tail(r.Context(), env)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetEnvelope(buf, result) })
}
func GetColorHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	result, status := get_color(r.Context())
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetColor(buf, result) })
}
func GetBadColorHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	result, status := get_bad_color(r.Context())
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetColor(buf, result) })
}
func PingHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	status := ping(r.Context())
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
}
func PingJunkHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	status := ping_junk(r.Context())
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
}
func PickHandler(w http.ResponseWriter, r *http.Request) {
	var color Color

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := GetColor(buf)
		if err != nil { return err }
		color = val
		return nil
	}) { return }
	if !IsColor(color) { w.WriteHeader(http.StatusBadRequest); return }

	result, status := pick(r.Context(), color)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetColor(buf, result) })
}
func PickBadHandler(w http.ResponseWriter, r *http.Request) {
	var color Color

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := GetColor(buf)
		if err != nil { return err }
		color = val
		return nil
	}) { return }
	if !IsColor(color) { w.WriteHeader(http.StatusBadRequest); return }

	result, status := pick_bad(r.Context(), color)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetColor(buf, result) })
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
	mux.HandleFunc("POST /echo", mw(EchoHandler))
	mux.HandleFunc("POST /echo_tail", mw(EchoTailHandler))
	mux.HandleFunc("POST /get_color", mw(GetColorHandler))
	mux.HandleFunc("POST /get_bad_color", mw(GetBadColorHandler))
	mux.HandleFunc("POST /ping", mw(PingHandler))
	mux.HandleFunc("POST /ping_junk", mw(PingJunkHandler))
	mux.HandleFunc("POST /pick", mw(PickHandler))
	mux.HandleFunc("POST /pick_bad", mw(PickBadHandler))
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
