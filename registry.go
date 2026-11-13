package main

import (
	"fmt"
	"sort"

	"google.golang.org/genai"
)

// ToolRegistry manages the tools available to the agent.
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates an empty registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool, rejecting duplicate names.
func (r *ToolRegistry) Register(t Tool) error {
	if t == nil {
		return fmt.Errorf("cannot register nil tool")
	}

	name := t.Name()
	if name == "" {
		return fmt.Errorf("cannot register tool with empty name")
	}

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %q already registered", name)
	}

	r.tools[name] = t
	return nil
}

// Lookup retrieves a registered tool by name.
func (r *ToolRegistry) Lookup(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// Declarations returns tool specifications for Gemini.
func (r *ToolRegistry) Declarations() []*genai.FunctionDeclaration {
	names := make([]string, 0, len(r.tools))

	for name := range r.tools {
		names = append(names, name)
	}

	sort.Strings(names)

	declarations := make([]*genai.FunctionDeclaration, 0, len(names))

	for _, name := range names {
		declarations = append(declarations, r.tools[name].Declaration())
	}

	return declarations
}
