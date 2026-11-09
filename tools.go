package main

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

const sparkHost = "spark-e2fe"

// gpuStatus runs nvidia-smi on the DGX Spark over SSH and returns its CSV output.
func gpuStatus(ctx context.Context) map[string]any {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh", "-o", "BatchMode=yes", sparkHost,
		"nvidia-smi",
		"--query-gpu=name,utilization.gpu,temperature.gpu,power.draw",
		"--format=csv,noheader")
	out, err := cmd.CombinedOutput()
	result := map[string]any{"output": strings.TrimSpace(string(out))}
	if err != nil {
		result["error"] = err.Error()
	}
	return result
}
