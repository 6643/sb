package sb

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func sampleEnvelope() *Envelope {
	return &Envelope{
		Item: &Item{
			Id:     7,
			Color:  ColorGreen,
			Tags:   []string{"hot", "new"},
			Active: true,
		},
		Items: []*Item{
			{
				Id:     8,
				Color:  ColorRed,
				Tags:   []string{"x"},
				Active: false,
			},
			{
				Id:     9,
				Color:  ColorBlue,
				Tags:   []string{"y", "z"},
				Active: true,
			},
		},
		Note: "core",
	}
}

func readAllBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body failed: %v", err)
	}
	return body
}

func encodeEnvelopeBytes(t *testing.T, env *Envelope) []byte {
	t.Helper()
	var out bytes.Buffer
	if err := SetAll(&out, func(buf *bytes.Buffer) error { return SetEnvelope(buf, env) }); err != nil {
		t.Fatalf("encode envelope failed: %v", err)
	}
	return out.Bytes()
}

func decodeEnvelopeBody(t *testing.T, body []byte) (*Envelope, bool) {
	t.Helper()
	buf := bytes.NewBuffer(body)
	env, err := ReadEnvelope(buf)
	if err != nil || buf.Len() != 0 || env == nil {
		return nil, false
	}
	return env, true
}

func writeOKBytes(t *testing.T, w http.ResponseWriter, payload []byte) {
	t.Helper()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		t.Fatalf("write response failed: %v", err)
	}
}

func newRPCServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/echo":
			env, ok := decodeEnvelopeBody(t, readAllBody(t, r))
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			writeOKBytes(t, w, encodeEnvelopeBytes(t, env))
		case "/echo_tail":
			env, ok := decodeEnvelopeBody(t, readAllBody(t, r))
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			encoded := encodeEnvelopeBytes(t, env)
			encoded = append(encoded, 0xEF)
			writeOKBytes(t, w, encoded)
		case "/get_color":
			var out bytes.Buffer
			if err := SetColor(&out, ColorBlue); err != nil {
				t.Fatalf("encode color failed: %v", err)
			}
			writeOKBytes(t, w, out.Bytes())
		case "/get_bad_color":
			writeOKBytes(t, w, []byte{0})
		case "/ping":
			if len(readAllBody(t, r)) != 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "/ping_junk":
			writeOKBytes(t, w, []byte{1})
		case "/pick":
			buf := bytes.NewBuffer(readAllBody(t, r))
			color, err := GetColor(buf)
			if err != nil || buf.Len() != 0 || !IsColor(color) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var out bytes.Buffer
			if err := SetColor(&out, color); err != nil {
				t.Fatalf("encode color failed: %v", err)
			}
			writeOKBytes(t, w, out.Bytes())
		case "/pick_bad":
			writeOKBytes(t, w, []byte{0})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestSerializationRoundTripAndGuards(t *testing.T) {
	env := sampleEnvelope()
	var buf bytes.Buffer
	if err := SetEnvelope(&buf, env); err != nil {
		t.Fatalf("set envelope failed: %v", err)
	}

	got, err := ReadEnvelope(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		t.Fatalf("read envelope failed: %v", err)
	}
	if !EqEnvelope(env, got) {
		t.Fatalf("round-trip mismatch")
	}

	if err := SetEnvelope(&bytes.Buffer{}, nil); err == nil {
		t.Fatalf("expected SetEnvelope(nil) to fail")
	}
	if err := ValidateEnvelopeList(EnvelopeList{nil}); err == nil {
		t.Fatalf("expected ValidateEnvelopeList with nil element to fail")
	}

	invalidItem := []byte{0x03, 0x01, 0x00, 0x00, 0x00, 0x00}
	if _, err := ReadItem(bytes.NewBuffer(invalidItem)); err == nil {
		t.Fatalf("expected invalid enum value in item decode to fail")
	}
}

func TestRPCHappyPath(t *testing.T) {
	srv := newRPCServer(t)
	defer srv.Close()

	client := NewClient(srv.URL + "/")
	env := sampleEnvelope()

	echo, code := CallEcho(client, context.Background(), env)
	if code != RpcOk {
		t.Fatalf("echo code = %v, want %v", code, RpcOk)
	}
	if !EqEnvelope(env, echo) {
		t.Fatalf("echo payload mismatch")
	}

	color, code := CallGetColor(client, nil)
	if code != RpcOk || color != ColorBlue {
		t.Fatalf("get_color unexpected: code=%v color=%v", code, color)
	}

	if code := CallPing(client, nil); code != RpcOk {
		t.Fatalf("ping code = %v, want %v", code, RpcOk)
	}

	picked, code := CallPick(client, context.Background(), ColorGreen)
	if code != RpcOk || picked != ColorGreen {
		t.Fatalf("pick unexpected: code=%v picked=%v", code, picked)
	}
}

func TestRPCRobustness(t *testing.T) {
	srv := newRPCServer(t)
	defer srv.Close()

	client := NewClient(srv.URL + "/")
	env := sampleEnvelope()

	if _, code := CallEcho(client, context.Background(), nil); code != RpcReqErr {
		t.Fatalf("nil struct arg should fail with RpcReqErr, got %v", code)
	}

	if _, code := CallPick(client, context.Background(), Color(0)); code != RpcReqErr {
		t.Fatalf("invalid enum arg should fail with RpcReqErr, got %v", code)
	}

	if _, code := CallGetBadColor(client, context.Background()); code != RpcRespErr {
		t.Fatalf("bad enum response should fail with RpcRespErr, got %v", code)
	}

	if _, code := CallPickBad(client, context.Background(), ColorRed); code != RpcRespErr {
		t.Fatalf("bad enum response should fail with RpcRespErr, got %v", code)
	}

	if _, code := CallEchoTail(client, context.Background(), env); code != RpcRespErr {
		t.Fatalf("tail-bytes response should fail with RpcRespErr, got %v", code)
	}

	if code := CallPingJunk(client, context.Background()); code != RpcRespErr {
		t.Fatalf("nil-return non-empty body should fail with RpcRespErr, got %v", code)
	}

	strictClient := NewClient(srv.URL + "/")
	strictClient.MaxRespBytes = 1
	if _, code := CallEcho(strictClient, context.Background(), env); code != RpcRespErr {
		t.Fatalf("response over max bytes should fail with RpcRespErr, got %v", code)
	}
}

func TestRPCTransportErrors(t *testing.T) {
	env := sampleEnvelope()

	if _, code := CallEcho(nil, context.Background(), env); code != RpcNoConn {
		t.Fatalf("nil client should fail with RpcNoConn, got %v", code)
	}

	timeoutSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer timeoutSrv.Close()

	timeoutClient := NewClient(timeoutSrv.URL + "/")
	timeoutClient.Timeout = 1 * time.Millisecond
	timeoutClient.Retries = 0
	if _, code := CallEcho(timeoutClient, context.Background(), env); code != RpcTimeout {
		t.Fatalf("timeout should fail with RpcTimeout, got %v", code)
	}

	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer authSrv.Close()

	authClient := NewClient(authSrv.URL + "/")
	if _, code := CallEcho(authClient, context.Background(), env); code != RpcNotAuth {
		t.Fatalf("401 should map to RpcNotAuth, got %v", code)
	}
}
