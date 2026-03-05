package sb

import (
	"bytes"
	"io"
	"net/http"
)

// --- API Handlers ---

func UserGetAbcHandler(w http.ResponseWriter, r *http.Request) {

	if !parseRequest(w, r) { return }

	result, status := user_get_abc(r.Context())
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetOrderStatus(buf, result) })
}
func UserGetAbcdHandler(w http.ResponseWriter, r *http.Request) {
	var page uint8
	var size uint8

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := GetU8(buf)
		if err != nil { return err }
		page = val
		return nil
	}, func(buf *bytes.Buffer) error {
		val, err := GetU8(buf)
		if err != nil { return err }
		size = val
		return nil
	}) { return }

	result, status := user_get_abcd(r.Context(), page, size)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetOrderStatus(buf, result) })
}
func UserSetSimInfoHandler(w http.ResponseWriter, r *http.Request) {
	var info *SimInfo

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := ReadSimInfo(buf)
		if err != nil { return err }
		info = val
		return nil
	}) { return }
	if err := ValidateSimInfo(info); err != nil { w.WriteHeader(http.StatusBadRequest); return }

	status := user_set_sim_info(r.Context(), info)
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
}
func GetCountHandler(w http.ResponseWriter, r *http.Request) {
	var page uint8

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := GetU8(buf)
		if err != nil { return err }
		page = val
		return nil
	}) { return }

	result, status := get_count(r.Context(), page)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetU8(buf, result) })
}
func GetBinHandler(w http.ResponseWriter, r *http.Request) {
	var page uint8

	if !parseRequest(w, r, func(buf *bytes.Buffer) error {
		val, err := GetU8(buf)
		if err != nil { return err }
		page = val
		return nil
	}) { return }

	result, status := get_bin(r.Context(), page)
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return SetBin(buf, result) })
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
	mux.HandleFunc("POST /get_count", mw(GetCountHandler))
	mux.HandleFunc("POST /get_bin", mw(GetBinHandler))
}

func RegisterUserApi(mux *http.ServeMux, mws ...Middleware) {
	mw := composeMiddleware(mws...)
	mux.HandleFunc("POST /user.get_abc", mw(UserGetAbcHandler))
	mux.HandleFunc("POST /user.get_abcd", mw(UserGetAbcdHandler))
	mux.HandleFunc("POST /user.set_sim_info", mw(UserSetSimInfoHandler))
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
