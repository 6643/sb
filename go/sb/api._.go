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
	sendResponse(w, U8(result))
}
func UserGetAbcdHandler(w http.ResponseWriter, r *http.Request) {
	var page U8
	var size U8

	if !parseRequest(w, r, &page, &size) { return }

	result, status := user_get_abcd(r.Context(), uint8(page), uint8(size))
	if !checkStatus(w, status) { return }
	sendResponse(w, U8(result))
}
func UserSetSimInfoHandler(w http.ResponseWriter, r *http.Request) {
	var info SimInfo

	if !parseRequest(w, r, &info) { return }

	status := user_set_sim_info(r.Context(), &info)
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
}
func GetCountHandler(w http.ResponseWriter, r *http.Request) {
	var page U8

	if !parseRequest(w, r, &page) { return }

	result, status := get_count(r.Context(), uint8(page))
	if !checkStatus(w, status) { return }
	sendResponse(w, U8(result))
}
func GetBinHandler(w http.ResponseWriter, r *http.Request) {
	var page U8

	if !parseRequest(w, r, &page) { return }

	result, status := get_bin(r.Context(), uint8(page))
	if !checkStatus(w, status) { return }
	sendResponse(w, Bin(result))
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

func RegisterUser(mux *http.ServeMux, mws ...Middleware) {
	mw := composeMiddleware(mws...)
	mux.HandleFunc("POST /user.get_abc", mw(UserGetAbcHandler))
	mux.HandleFunc("POST /user.get_abcd", mw(UserGetAbcdHandler))
	mux.HandleFunc("POST /user.set_sim_info", mw(UserSetSimInfoHandler))
}


// --- 内部辅助函数 ---

func checkStatus(w http.ResponseWriter, status RpcErrCode) bool {
	if status == RpcOk { return true }
	w.WriteHeader(int(status)); return false
}

func parseRequest(w http.ResponseWriter, r *http.Request, args ...Deserializable) bool {
	if len(args) == 0 { return true }
	body, err := io.ReadAll(r.Body); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if err := GetAll(bytes.NewBuffer(body), args...); err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	return true
}

func sendResponse(w http.ResponseWriter, result Serializable) {
	var buf bytes.Buffer
	if err := SetAll(&buf, result); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	w.Write(buf.Bytes())
}
