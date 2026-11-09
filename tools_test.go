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
	t.Logf("ouptut: %v", r["output"])
}
