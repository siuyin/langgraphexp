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
		if len(resp.Message.ToolCalls) > 0 {
			log.Printf("Role: %#v, %#v\n", resp.Message.Role, resp.Message.ToolCalls)
			for _, toolCall := range resp.Message.ToolCalls {
				if toolCall.Function.Name == "get_menu" {
					msg := api.Message{Role: "tool", Content: get_menu()}
					log.Printf("%#v\n", msg)
					messages = append(messages, msg)
				}
			}
			return nil
		}
		fmt.Print(resp.Message.Content)
		return nil
	}

	fmt.Println(welcomeMessage)
	for {
		fmt.Printf("len: %d\n", len(messages))
		for i, m := range messages {
			fmt.Printf("%d: %s ", i, m.Role)
		}
		fmt.Println()
		lastMsg := messages[len(messages)-1]
		if lastMsg.Role == "tool" {
			req := &api.ChatRequest{
				Model:    dflt.EnvString("AI_MODEL", "gemma3:4b"),
				Messages: messages,
				Tools:    []api.Tool{getMenuTool},
			}
			send(client, req, respFunc)
			//fmt.Printf("LastMsg %#v,\n ************ sent\n", lastMsg)

		}

		msg := getInput()
		if strings.ToLower(msg) == "q" {
			break
		}

		messages = append(messages, api.Message{Role: "user", Content: msg})
		req := &api.ChatRequest{
			Model:    dflt.EnvString("AI_MODEL", "gemma3:4b"),
			Messages: messages,
			Tools:    []api.Tool{getMenuTool},
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
