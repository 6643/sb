package sb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerStatusMapsRpcNoConnTo500(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	if checkStatus(rec, RpcNoConn) {
		t.Fatalf("checkStatus should reject RpcNoConn")
	}
	if got, want := rec.Code, http.StatusInternalServerError; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
}

func TestHandlerStatusPreservesKnownHTTPStatuses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		status RpcErrCode
		want   int
	}{
		{name: "timeout", status: RpcTimeout, want: http.StatusRequestTimeout},
		{name: "not auth", status: RpcNotAuth, want: http.StatusUnauthorized},
		{name: "not exist", status: RpcNotExist, want: http.StatusNotFound},
		{name: "bad request", status: RpcReqErr, want: http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			if checkStatus(rec, tc.status) {
				t.Fatalf("checkStatus should reject %v", tc.status)
			}
			if got := rec.Code; got != tc.want {
				t.Fatalf("status = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestHandlerStatusMapsUnknownInternalStatusTo500(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	if checkStatus(rec, RpcErrCode(499)) {
		t.Fatalf("checkStatus should reject unknown internal status")
	}
	if got, want := rec.Code, http.StatusInternalServerError; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
}
