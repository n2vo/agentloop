package main

import (
	"context"
	"testing"
)

func TestGPUStatus(t *testing.T) {
	r := gpuStatus(context.Background())
	if e, ok := r["error"]; ok {
		t.Fatalf("gpuStatus error: %v (output: %v)", e, r["output"])
	}
	t.Logf("output: %v", r["output"])
}

func TestGPUStatusToolExecutor(t *testing.T) {
	registry := NewToolRegistry()

	if err := registry.Register(&GPUStatusTool{}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	executor := NewToolExecutor(registry)

	result, err := executor.Execute(
		context.Background(),
		"gpu_status",
		nil,
	)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if toolErr, ok := result["error"]; ok {
		t.Fatalf("gpu_status failed: %v", toolErr)
	}

	if result["output"] == nil || result["output"] == "" {
		t.Fatalf("Expected non-empty GPU telemetry")
	}

	t.Logf("GPU telemetry: %v", result["output"])
}
