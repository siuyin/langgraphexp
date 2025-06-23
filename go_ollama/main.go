package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: systemInstruction,
		},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	fmt.Println(welcomeMessage)
	for {
		msg := getInput()
		if strings.ToLower(msg) == "q" {
			break
		}

		messages = append(messages, api.Message{Role: "user", Content: msg})
		req := &api.ChatRequest{
			Model:    dflt.EnvString("AI_MODEL", "gemma3:4b"),
			Messages: messages,
		}
		send(client, req, respFunc)
	}
}

func getInput() string {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	return sc.Text()
}

func send(client *api.Client, req *api.ChatRequest, respFunc func(r api.ChatResponse) error) {
	ctx := context.Background()
	err := client.Chat(ctx, req, respFunc)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
