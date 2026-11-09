package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

const model = "gemini-3.5-flash"

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendGeminiAPI})
	if err != nil{
		log.Fatal(err)
	}
	resp, err := client.Models.GenerateContent(ctx, model, genai.Text("Reply with the single word: ready"), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Text())
}
