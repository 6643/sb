package sb

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(baseURL string, timeout time.Duration, retries int) *Client {
	client := NewClient(baseURL)
	client.Timeout = timeout
	client.Retries = retries
	return client
}

func TestDoClientSlowBodyWithinTimeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x11})
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write([]byte{0x22})
	}))
	defer server.Close()

	body, status := doClient(newTestClient(server.URL, 200*time.Millisecond, 0), context.Background(), "/", nil)
	if status != RpcOk {
		t.Fatalf("status = %v, want %v", status, RpcOk)
	}
	if got, want := fmt.Sprintf("%x", body), "1122"; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func TestDoClientBodyReadTimeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x11})
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(80 * time.Millisecond)
		_, _ = w.Write([]byte{0x22})
	}))
	defer server.Close()

	_, status := doClient(newTestClient(server.URL, 30*time.Millisecond, 0), context.Background(), "/", nil)
	if status != RpcTimeout {
		t.Fatalf("status = %v, want %v", status, RpcTimeout)
	}
}

func TestDoClientTruncatedBodyReturnsRespErr(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatalf("response writer does not support hijacking")
		}
		conn, rw, err := hj.Hijack()
		if err != nil {
			t.Fatalf("Hijack failed: %v", err)
		}
		defer conn.Close()
		_, _ = rw.WriteString("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: 2\r\n\r\n")
		_, _ = rw.Write([]byte{0x11})
		_ = rw.Flush()
	}))
	defer server.Close()

	_, status := doClient(newTestClient(server.URL, 200*time.Millisecond, 0), context.Background(), "/", nil)
	if status != RpcRespErr {
		t.Fatalf("status = %v, want %v", status, RpcRespErr)
	}
}

func TestDoClientBodyReadTimeoutThenRetrySuccess(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x11})
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		if attempts == 1 {
			time.Sleep(80 * time.Millisecond)
		}
		_, _ = w.Write([]byte{0x22})
	}))
	defer server.Close()

	body, status := doClient(newTestClient(server.URL, 30*time.Millisecond, 1), context.Background(), "/", nil)
	if status != RpcOk {
		t.Fatalf("status = %v, want %v", status, RpcOk)
	}
	if got, want := fmt.Sprintf("%x", body), "1122"; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestDoClientRetriesOnHTTP408(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
		_, _ = w.Write([]byte{0x22})
	}))
	defer server.Close()

	body, status := doClient(newTestClient(server.URL, 200*time.Millisecond, 1), context.Background(), "/", nil)
	if status != RpcOk {
		t.Fatalf("status = %v, want %v", status, RpcOk)
	}
	if got, want := fmt.Sprintf("%x", body), "22"; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestDoClientCancelDuringBackoff(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestTimeout)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, status := doClient(newTestClient(server.URL, 200*time.Millisecond, 2), ctx, "/", nil)
	if status != RpcTimeout {
		t.Fatalf("status = %v, want %v", status, RpcTimeout)
	}
}

var _ net.Conn
var _ *bufio.ReadWriter
