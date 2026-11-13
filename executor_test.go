package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

type slowTool struct {
	mockTool
}

func (s *slowTool) Execute(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestToolExecutor(t *testing.T) {
	registry := NewToolRegistry()
	tool := &mockTool{name: "test_tool"}

	if err := registry.Register(tool); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	executor := NewToolExecutor(registry)
	ctx := context.Background()

	t.Run("valid execution", func(t *testing.T) {
		result, err := executor.Execute(
			ctx,
			"test_tool",
			nil,
		)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result["status"] != "ok" {
			t.Fatalf("Unexpected result: %v", result)
		}
	})

	t.Run("unknown tool", func(t *testing.T) {
		_, err := executor.Execute(ctx, "unknown_tool", nil)
		if err == nil || !strings.Contains(err.Error(), "unknown_tool") {
			t.Fatalf("Expected unknown tool error, got %v", err)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		emptyExecutor := NewToolExecutor(nil)
		_, err := emptyExecutor.Execute(ctx, "test_tool", nil)
		if err == nil {
			t.Fatalf("Expected nil registry error")
		}
	})

	t.Run("invalid arguments", func(t *testing.T) {
		_, err := executor.Execute(
			ctx,
			"test_tool",
			map[string]any{"host": "unauthorized-host"},
		)

		if err == nil {
			t.Fatal("Expected argument validation error")
		}

		if !strings.Contains(err.Error(), "invalid arguments") {
			t.Fatalf("Expected validation error, got %v", err)
		}
	})

	if tool.executions != 1 {
		t.Fatalf(
			"Expected only the earlier valid execution, got %d executions",
			tool.executions,
		)
	}
}

func TestToolExecutorTimeout(t *testing.T) {
	registry := NewToolRegistry()

	if err := registry.Register(&slowTool{
		mockTool: mockTool{name: "slow_tool"},
	}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	executor := NewToolExecutor(registry)
	executor.timeout = 50 * time.Millisecond

	start := time.Now()

	_, err := executor.Execute(
		context.Background(),
		"slow_tool",
		nil,
	)

	if err != context.DeadlineExceeded {
		t.Fatalf("Expected deadline exceeded, got %v", err)
	}

	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("Timeout took too long: %v", elapsed)
	}
}

func TestToolExecutorCancellation(t *testing.T) {
	registry := NewToolRegistry()

	if err := registry.Register(&slowTool{
		mockTool: mockTool{name: "slow_tool"},
	}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	executor := NewToolExecutor(registry)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := executor.Execute(ctx, "slow_tool", nil)

	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, got %v", err)
	}
}
