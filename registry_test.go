package main

import (
	"context"
	"testing"

	"google.golang.org/genai"
)

type mockTool struct {
	name       string
	executions int
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name: m.name,
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: map[string]*genai.Schema{},
		},
	}
}

func (m *mockTool) Execute(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	m.executions++
	return map[string]any{"status": "ok"}, nil
}

func TestToolRegistry(t *testing.T) {
	registry := NewToolRegistry()
	tool := &mockTool{name: "test_tool"}

	// Register a tool.
	if err := registry.Register(tool); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Look up the registered tool.
	found, ok := registry.Lookup("test_tool")
	if !ok || found != tool {
		t.Fatal("Lookup failed to return the registered tool")
	}

	// Reject duplicate registration.
	if err := registry.Register(tool); err == nil {
		t.Fatal("Expected duplicate registration error")
	}

	// Reject an unregistered tool lookup
	if _, ok := registry.Lookup("unknown_tool"); ok {
		t.Fatal("Unexpected lookup success for unknown tool")
	}

	// Verify Gemini declarations.
	declarations := registry.Declarations()
	if len(declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(declarations))
	}

	if declarations[0].Name != "test_tool" {
		t.Fatalf("Unexpected declaration: %s", declarations[0].Name)
	}

}

func TestToolRegistryDeclarationOrder(t *testing.T) {
	registry := NewToolRegistry()

	// Intentionally register tools out of alphabetical order.
	for _, name := range []string{"zeta", "alpha", "middle"} {
		if err := registry.Register(&mockTool{name: name}); err != nil {
			t.Fatalf("Register(%q) failed: %v", name, err)
		}
	}

	expected := []string{"alpha", "middle", "zeta"}
	declarations := registry.Declarations()

	if len(declarations) != len(expected) {
		t.Fatalf("Expected %d declarations, got %d",
			len(expected), len(declarations))
	}

	for i, declaration := range declarations {
		if declaration.Name != expected[i] {
			t.Errorf("Declaration %d: expected %q, got %q",
				i, expected[i], declaration.Name)
		}
	}
}
