package main

import (
	"testing"

	"google.golang.org/genai"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		decl    *genai.FunctionDeclaration
		args    map[string]any
		wantErr bool
	}{
		{
			name: "valid empty arguments",
			decl: gpuStatusDecl,
			args: map[string]any{},
		},
		{
			name: "nil arguments",
			decl: gpuStatusDecl,
			args: nil,
		},
		{
			name:    "unexpected arguments",
			decl:    gpuStatusDecl,
			args:    map[string]any{"host": "other-host"},
			wantErr: true,
		},
		{
			name:    "nil declaration",
			decl:    nil,
			args:    nil,
			wantErr: true,
		},
		{
			name: "missing schema",
			decl: &genai.FunctionDeclaration{
				Name: "no_schema",
			},
			wantErr: true,
		},
		{
			name: "non-object schema",
			decl: &genai.FunctionDeclaration{
				Name: "invalid_schema",
				Parameters: &genai.Schema{
					Type: genai.TypeString,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateArgs(tt.decl, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ValidateArgs() error = %v, wantErr = %v",
					err, tt.wantErr,
				)
			}
		})
	}
}
