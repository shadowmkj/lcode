package main

import (
	"fmt"
	"os"
)

const (
	GraphQLEndpoint = "https://leetcode.com/graphql"
	UserAgent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: lcode <command> <args>")
		fmt.Println("Commands:")
		fmt.Println("  pick <slug|id>   Fetch problem and open editor")
		fmt.Println("  submit <file>    Submit solution to LeetCode")
		fmt.Println("  auth             Authenticate with LeetCode (Login with Chrome first)")
		return
	}
	if len(os.Args) < 3 && os.Args[1] != "auth" {
		fmt.Println("Usage: lcode <command> <args>")
		fmt.Println("Commands:")
		fmt.Println("  pick <slug|id>   Fetch problem and open editor")
		fmt.Println("  submit <file>    Submit solution to LeetCode")
		fmt.Println("  auth             Authenticate with LeetCode (Login with Chrome first)")
		return
	}

	command := os.Args[1]

	var httpClient HttpClient
	httpClient.initClient()
	switch command {
	case "pick":
		arg := os.Args[2]
		handlePick(httpClient, arg)
	case "submit":
		arg := os.Args[2]
		handleSubmit(httpClient, arg)
	case "auth":
		handleAuth(httpClient)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func handleAuth(client HttpClient) {
	initEnvFile()
}
