package sb

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func sampleSizeBlobWrap() *BlobWrap {
	nums := make([]uint8, 255)
	for i := range nums {
		nums[i] = uint8(i)
	}
	return &BlobWrap{
		TextValue: strings.Repeat("x", 1024),
		BinValue:  bytes.Repeat([]byte{0xAB}, 1024),
		Nums:      nums,
		Level:     LevelHigh,
	}
}

func sizeReadAllBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body failed: %v", err)
	}
	return body
}

func sizeEncodeBlobWrap(t *testing.T, payload *BlobWrap) []byte {
	t.Helper()
	var out bytes.Buffer
	if err := SetBlobWrap(&out, payload); err != nil {
		t.Fatalf("encode blob_wrap failed: %v", err)
	}
	return out.Bytes()
}

func sizeDecodeBlobWrap(body []byte) (*BlobWrap, bool) {
	buf := bytes.NewBuffer(body)
	payload, err := ReadBlobWrap(buf)
	if err != nil || payload == nil || buf.Len() != 0 {
		return nil, false
	}
	return payload, true
}

func TestSizeSerializationBoundaryAndEnumGuard(t *testing.T) {
	payload := sampleSizeBlobWrap()
	var buf bytes.Buffer
	if err := SetBlobWrap(&buf, payload); err != nil {
		t.Fatalf("set blob_wrap failed: %v", err)
	}

	got, err := ReadBlobWrap(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		t.Fatalf("read blob_wrap failed: %v", err)
	}
	if !EqBlobWrap(payload, got) {
		t.Fatalf("round-trip mismatch")
	}

	bad := bytes.NewBuffer(nil)
	if err := SetAll(
		bad,
		func(buf *bytes.Buffer) error { return SetText(buf, "bad") },
		func(buf *bytes.Buffer) error { return SetBin(buf, []byte{1}) },
		func(buf *bytes.Buffer) error { return SetU8List(buf, []uint8{1, 2}) },
		func(buf *bytes.Buffer) error { return SetU8(buf, 0) },
	); err != nil {
		t.Fatalf("build invalid blob_wrap bytes failed: %v", err)
	}
	if _, err := ReadBlobWrap(bytes.NewBuffer(bad.Bytes())); err == nil {
		t.Fatalf("expected invalid enum value in blob_wrap decode to fail")
	}
}

func TestSizeRPCRoundTripAndRequestGuard(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/round" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		payload, ok := sizeDecodeBlobWrap(sizeReadAllBody(t, r))
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(sizeEncodeBlobWrap(t, payload)); err != nil {
			t.Fatalf("write response failed: %v", err)
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL + "/")
	payload := sampleSizeBlobWrap()

	got, code := CallRound(client, context.Background(), payload)
	if code != RpcOk {
		t.Fatalf("round code = %v, want %v", code, RpcOk)
	}
	if !EqBlobWrap(payload, got) {
		t.Fatalf("rpc round payload mismatch")
	}

	if _, code := CallRound(client, context.Background(), nil); code != RpcReqErr {
		t.Fatalf("nil struct arg should fail with RpcReqErr, got %v", code)
	}

	strict := NewClient(srv.URL + "/")
	strict.MaxRespBytes = 16
	if _, code := CallRound(strict, context.Background(), payload); code != RpcRespErr {
		t.Fatalf("response over max bytes should fail with RpcRespErr, got %v", code)
	}
}

func TestSizeRPCMalformedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/round" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte{1}); err != nil {
			t.Fatalf("write response failed: %v", err)
		}
	}))
	defer srv.Close()

	client := NewClient(srv.URL + "/")
	if _, code := CallRound(client, context.Background(), sampleSizeBlobWrap()); code != RpcRespErr {
		t.Fatalf("malformed response should fail with RpcRespErr, got %v", code)
	}
}
