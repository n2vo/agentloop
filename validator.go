package main

import (
	"fmt"

	"google.golang.org/genai"
)

// ValidateArgs checks tool arguments against the declared input schema.
// Initially supports parameterless tools only.
func ValidateArgs(
	decl *genai.FunctionDeclaration,
	args map[string]any,
) error {
	if decl == nil {
		return fmt.Errorf("nil tool declaration")
	}

	if decl.Parameters == nil {
		return fmt.Errorf("tool %q has no input schema", decl.Name)
	}

	if decl.Parameters.Type != genai.TypeObject {
		return fmt.Errorf("tool %q requires an object schema", decl.Name)
	}

	if len(decl.Parameters.Properties) == 0 && len(args) != 0 {
		return fmt.Errorf(
			"tool %q accepts no arguments; received %d",
			decl.Name, len(args),
		)
	}

	return nil
}
