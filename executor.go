package main

import (
	"context"
	"fmt"
	"time"
)

// ToolExecutor validates and executes registered tools.
type ToolExecutor struct {
	registry *ToolRegistry
	timeout  time.Duration
}

// NewToolExecutor creates an executor for the given registry.
func NewToolExecutor(registry *ToolRegistry) *ToolExecutor {
	return &ToolExecutor{
		registry: registry,
		timeout:  10 * time.Second,
	}
}

// Execute resolves a tool, validates its arguments, and invokes it.
func (e *ToolExecutor) Execute(
	ctx context.Context,
	name string,
	args map[string]any,
) (map[string]any, error) {
	if e.registry == nil {
		return nil, fmt.Errorf("tool registry is nil")
	}

	tool, ok := e.registry.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}

	if err := ValidateArgs(tool.Declaration(), args); err != nil {
		return nil, fmt.Errorf("invalid arguments for %q: %w", name, err)
	}

	//return tool.Execute(ctx, args)
	toolCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	return tool.Execute(toolCtx, args)
}
