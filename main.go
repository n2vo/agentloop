package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

const (
	model = "gemini-3.5-flash"
	maxSteps = 5
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendGeminiAPI})
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{gpuStatusDecl}}},
	}
	question := "Which GPU is installed, and what are its utilization and temperature right now?"
	contents := []*genai.Content{genai.NewContentFromText(question, genai.RoleUser)}
	
	for step := 1; step <= maxSteps; step++ {
		resp, err := client.Models.GenerateContent(ctx, model, contents, config)
		if err != nil {
			log.Fatal(err)
		}
		calls := resp.FunctionCalls()
		if len(calls) == 0 {
			fmt.Println(resp.Text())
			return
		}
		contents = append(contents, resp.Candidates[0].Content)
		var parts []*genai.Part
		for _, call := range calls {
			var result map[string]any
			switch call.Name {
			case gpuStatusDecl.Name:
				result = gpuStatus(ctx)
			default:
				result = map[string]any{"error": "unknown tool: " + call.Name}
			}
			part := genai.NewPartFromFunctionResponse(call.Name, result)
			part.FunctionResponse.ID = call.ID
			parts = append(parts, part)
		}
		contents = append(contents, genai.NewContentFromParts(parts, genai.RoleUser))
	}
	log.Fatalf("stopped: reached step cap of %d without a final answer", maxSteps)
}
