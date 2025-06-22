package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/siuyin/dflt"
	"google.golang.org/genai"
)

// const llm = "gemini-2.5-flash"
const llm = "gemini-2.5-flash-lite-preview-06-17"

var client *genai.Client

func init() {
	clientCfg := &genai.ClientConfig{
		APIKey:  dflt.EnvString("GOOGLE_API_KEY", "your-key-here"),
		Backend: genai.BackendGeminiAPI,
	}
	var err error

	client, err = genai.NewClient(context.Background(), clientCfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	thinkingBudget := int32(0) // 0 == Disable thinking
	genCfg := &genai.GenerateContentConfig{
		SystemInstruction: genai.Text("You are a philosopher, respond with one sentence.")[0],
		ThinkingConfig:    &genai.ThinkingConfig{ThinkingBudget: &thinkingBudget},
	}
	chatStream(genCfg)

}

func chatStream(genCfg *genai.GenerateContentConfig) {
	ctx := context.Background()

	chat, err := client.Chats.Create(ctx, llm, genCfg, nil) // ctx context.Context, model string, config *GenerateContentConfig, history []*Content
	if err != nil {
		log.Fatal(err)
	}

	for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: "Hi, my name is Siu Yin. I live in Bukit Timah, Singapore"}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(result.Candidates[0].Content.Parts[0].Text)
	}

	fmt.Println("\n---1")
	for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: "What is my name and address?"}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(result.Candidates[0].Content.Parts[0].Text)
	}

	fmt.Println("\n---2")
	for _, h := range chat.History(false) {
		fmt.Println(h.Role + ": " + h.Parts[0].Text)
	}

}

func debugPrint[T any](r *T) {

	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response))
}

func genContentStream(p string, genCfg *genai.GenerateContentConfig) {
	ctx := context.Background()
	prompt := genai.Text(p)
	for r, err := range client.Models.GenerateContentStream(ctx, llm, prompt, genCfg) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(r.Candidates[0].Content.Parts[0].Text)
	}
}
