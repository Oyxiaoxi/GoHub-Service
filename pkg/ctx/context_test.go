package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "req-123456"
	
	ctx = WithRequestID(ctx, requestID)
	result := GetRequestID(ctx)
	
	assert.Equal(t, requestID, result)
}

func TestGetRequestID_NotSet(t *testing.T) {
	ctx := Background()
	result := GetRequestID(ctx)
	
	assert.Empty(t, result)
}

func TestWithUserID(t *testing.T) {
	ctx := Background()
	userID := "user-789"
	
	ctx = WithUserID(ctx, userID)
	result := GetUserID(ctx)
	
	assert.Equal(t, userID, result)
}

func TestGetUserID_NotSet(t *testing.T) {
	ctx := Background()
	result := GetUserID(ctx)
	
	assert.Empty(t, result)
}

func TestWithTraceID(t *testing.T) {
	ctx := Background()
	traceID := "trace-abc123"
	
	ctx = WithTraceID(ctx, traceID)
	result := GetTraceID(ctx)
	
	assert.Equal(t, traceID, result)
}

func TestGetTraceID_NotSet(t *testing.T) {
	ctx := Background()
	result := GetTraceID(ctx)
	
	assert.Empty(t, result)
}

func TestWithTimeout(t *testing.T) {
	parent := Background()
	timeout := 5 * time.Second
	
	ctx, cancel := WithTimeout(parent, timeout)
	defer cancel()
	
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(deadline) <= timeout)
}

func TestWithDefaultTimeout(t *testing.T) {
	parent := Background()
	
	ctx, cancel := WithDefaultTimeout(parent)
	defer cancel()
	
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(deadline) <= DefaultTimeout)
}

func TestContextChaining(t *testing.T) {
	ctx := Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithUserID(ctx, "user-456")
	ctx = WithTraceID(ctx, "trace-789")
	
	assert.Equal(t, "req-123", GetRequestID(ctx))
	assert.Equal(t, "user-456", GetUserID(ctx))
	assert.Equal(t, "trace-789", GetTraceID(ctx))
}

func TestContextWithTimeoutCancellation(t *testing.T) {
	parent := Background()
	ctx, cancel := WithTimeout(parent, 100*time.Millisecond)
	
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done immediately")
	case <-time.After(10 * time.Millisecond):
		// 正常情况
	}
	
	cancel()
	
	select {
	case <-ctx.Done():
		// 正常，context 应该被取消
	case <-time.After(100 * time.Millisecond):
		t.Fatal("context should be cancelled")
	}
}
