package main

import "github.com/ollama/ollama/api"

var getMenuTool = api.Tool{
	Type:     "function",
	Function: api.ToolFunction{Name: "get_menu", Description: "returns the menu of our baristabot cafe"},
}
