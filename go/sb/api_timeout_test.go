package sb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutMiddlewareInjectsDeadlineAndReturnsTimeoutStatus(t *testing.T) {
	handler := TimeoutMiddleware(20 * time.Millisecond)(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Deadline(); !ok {
			t.Fatalf("expected request context deadline")
		}

		<-r.Context().Done()
		if err := r.Context().Err(); err != context.DeadlineExceeded {
			t.Fatalf("expected deadline exceeded, got %v", err)
		}

		status := normalizeHandlerStatus(r.Context(), RpcOk)
		if !checkStatus(w, status) {
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/timeout", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusRequestTimeout {
		t.Fatalf("expected status %d, got %d", http.StatusRequestTimeout, rec.Code)
	}
}

func TestTimeoutMiddlewareDisabledPassesThrough(t *testing.T) {
	handler := TimeoutMiddleware(0)(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Deadline(); ok {
			t.Fatalf("expected request context without deadline")
		}

		status := normalizeHandlerStatus(r.Context(), RpcOk)
		if !checkStatus(w, status) {
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/no-timeout", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
