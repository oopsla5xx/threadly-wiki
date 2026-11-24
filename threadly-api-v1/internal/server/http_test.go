package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestRunShutsDownGracefullyOnContextCancel(t *testing.T) {
	// Reserve a free port so the test can dial the server it starts.
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := probe.Addr().String()
	_ = probe.Close()

	httpServer := &http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() { done <- Run(ctx, httpServer, zap.NewNop(), time.Second) }()

	var response *http.Response
	for range 50 {
		if response, err = http.Get("http://" + address); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("server never became reachable: %v", err)
	}
	_ = response.Body.Close()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after context cancel")
	}

	if _, err := net.DialTimeout("tcp", address, 100*time.Millisecond); err == nil {
		t.Fatal("server still accepting connections after shutdown")
	}
}

func TestRunReturnsErrorWhenAddressInUse(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()

	httpServer := &http.Server{Addr: occupied.Addr().String()}

	if err := Run(context.Background(), httpServer, zap.NewNop(), time.Second); err == nil {
		t.Fatal("Run returned nil, want listen error")
	}
}
