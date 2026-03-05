package sb

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sampleOptionalMatrix() *Matrix {
	return &Matrix{
		Single: &Child{Id: 7},
		Many: []*Child{
			{Id: 8},
			{Id: 9},
		},
		Kind:  KindB,
		Kinds: []Kind{KindA, KindB},
	}
}

func optionalReadAllBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body failed: %v", err)
	}
	return body
}

func optionalEncodeMatrix(t *testing.T, m *Matrix) []byte {
	t.Helper()
	var out bytes.Buffer
	if err := SetMatrix(&out, m); err != nil {
		t.Fatalf("encode matrix failed: %v", err)
	}
	return out.Bytes()
}

func optionalDecodeMatrix(body []byte) (*Matrix, bool) {
	buf := bytes.NewBuffer(body)
	m, err := ReadMatrix(buf)
	if err != nil || m == nil || buf.Len() != 0 {
		return nil, false
	}
	return m, true
}

func TestOptionalMatrixSerializationAndGuards(t *testing.T) {
	full := sampleOptionalMatrix()
	var buf bytes.Buffer
	if err := SetMatrix(&buf, full); err != nil {
		t.Fatalf("set matrix failed: %v", err)
	}
	got, err := ReadMatrix(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		t.Fatalf("read matrix failed: %v", err)
	}
	if !EqMatrix(full, got) {
		t.Fatalf("full matrix round-trip mismatch")
	}

	empty := &Matrix{}
	buf.Reset()
	if err := SetMatrix(&buf, empty); err != nil {
		t.Fatalf("set empty matrix failed: %v", err)
	}
	emptyGot, err := ReadMatrix(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		t.Fatalf("read empty matrix failed: %v", err)
	}
	if !EqMatrix(empty, emptyGot) {
		t.Fatalf("empty matrix round-trip mismatch")
	}

	if err := SetMatrix(&bytes.Buffer{}, nil); err == nil {
		t.Fatalf("expected SetMatrix(nil) to fail")
	}
	if err := ValidateMatrixList(MatrixList{nil}); err == nil {
		t.Fatalf("expected ValidateMatrixList with nil element to fail")
	}

	bad := bytes.NewBuffer(nil)
	if err := SetAll(
		bad,
		func(buf *bytes.Buffer) error { return SetU8(buf, 1<<2) },
		func(buf *bytes.Buffer) error { return SetU8(buf, 0) },
	); err != nil {
		t.Fatalf("build invalid matrix bytes failed: %v", err)
	}
	if _, err := ReadMatrix(bytes.NewBuffer(bad.Bytes())); err == nil {
		t.Fatalf("expected invalid enum value in matrix decode to fail")
	}
}

func TestOptionalMatrixRPCHappyAndRequestGuard(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/submit":
			m, ok := optionalDecodeMatrix(optionalReadAllBody(t, r))
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err := ValidateMatrix(m); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "/fetch":
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(optionalEncodeMatrix(t, sampleOptionalMatrix())); err != nil {
				t.Fatalf("write response failed: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL + "/")
	want := sampleOptionalMatrix()

	if code := CallSubmit(client, context.Background(), want); code != RpcOk {
		t.Fatalf("submit code = %v, want %v", code, RpcOk)
	}

	got, code := CallFetch(client, context.Background())
	if code != RpcOk {
		t.Fatalf("fetch code = %v, want %v", code, RpcOk)
	}
	if !EqMatrix(want, got) {
		t.Fatalf("fetch payload mismatch")
	}

	if code := CallSubmit(client, context.Background(), nil); code != RpcReqErr {
		t.Fatalf("nil struct arg should fail with RpcReqErr, got %v", code)
	}
	badReq := sampleOptionalMatrix()
	badReq.Kind = Kind(9)
	if code := CallSubmit(client, context.Background(), badReq); code != RpcReqErr {
		t.Fatalf("invalid enum arg should fail with RpcReqErr, got %v", code)
	}

	strict := NewClient(srv.URL + "/")
	strict.MaxRespBytes = 16
	if _, code := CallFetch(strict, context.Background()); code != RpcRespErr {
		t.Fatalf("response over max bytes should fail with RpcRespErr, got %v", code)
	}
}

func TestOptionalMatrixRPCResponseGuards(t *testing.T) {
	submitJunkSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/submit" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte{1}); err != nil {
			t.Fatalf("write response failed: %v", err)
		}
	}))
	defer submitJunkSrv.Close()

	submitJunkClient := NewClient(submitJunkSrv.URL + "/")
	if code := CallSubmit(submitJunkClient, context.Background(), sampleOptionalMatrix()); code != RpcRespErr {
		t.Fatalf("nil-return non-empty body should fail with RpcRespErr, got %v", code)
	}

	fetchTailSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/fetch" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		payload := optionalEncodeMatrix(t, sampleOptionalMatrix())
		payload = append(payload, 0xEE)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(payload); err != nil {
			t.Fatalf("write response failed: %v", err)
		}
	}))
	defer fetchTailSrv.Close()

	fetchTailClient := NewClient(fetchTailSrv.URL + "/")
	if _, code := CallFetch(fetchTailClient, context.Background()); code != RpcRespErr {
		t.Fatalf("tail-bytes response should fail with RpcRespErr, got %v", code)
	}

	fetchBadSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/fetch" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		bad := bytes.NewBuffer(nil)
		if err := SetAll(
			bad,
			func(buf *bytes.Buffer) error { return SetU8(buf, 1<<2) },
			func(buf *bytes.Buffer) error { return SetU8(buf, 0) },
		); err != nil {
			t.Fatalf("build invalid matrix bytes failed: %v", err)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bad.Bytes()); err != nil {
			t.Fatalf("write response failed: %v", err)
		}
	}))
	defer fetchBadSrv.Close()

	fetchBadClient := NewClient(fetchBadSrv.URL + "/")
	if _, code := CallFetch(fetchBadClient, context.Background()); code != RpcRespErr {
		t.Fatalf("invalid enum response should fail with RpcRespErr, got %v", code)
	}
}
