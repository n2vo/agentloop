package main

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/genai"
)

var sparkHost = "spark-e2fe"

// gpuStatusDecl describes the gpu_status tool to the model.
var gpuStatusDecl = &genai.FunctionDeclaration{
	Name:			"gpu_status",
	Description:	"Returns the current status of each GPU on the DGX Spark as CSV lines: name, utilization percent, temperature in Celsius, power draw in watts.",
}


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
