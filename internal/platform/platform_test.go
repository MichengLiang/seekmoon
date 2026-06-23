package platform

import (
	"context"
	"testing"
	"time"
)

func TestHTTPClientDefaultTimeout(t *testing.T) {
	client := NewHTTPClient(0)
	if client.Timeout != 30*time.Second {
		t.Fatalf("timeout = %s", client.Timeout)
	}
}

func TestOSFSHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := (OSFS{}).MkdirAll(ctx, t.TempDir(), 0o755); err == nil {
		t.Fatal("MkdirAll should fail on canceled context")
	}
}
