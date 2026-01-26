package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/charmbracelet/glamour"
)

const (
	GraphQLEndpoint = "https://leetcode.com/graphql"
	UserAgent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

func main() {
	// var LEETCODE_SESSION, LEETCODE_CSRF_TOKEN = loadEnv()
	if len(os.Args) < 3 || os.Args[1] != "pick" {
		fmt.Println("Usage: lcode pick <problem-slug>")
		return
	}
	input := os.Args[2]
	slug := input

	var HttpClient HttpClient
	HttpClient.initClient()
	if _, err := strconv.Atoi(input); err == nil {
		fmt.Printf("Searching for Problem ID: %s...\n", input)
		foundSlug, err := HttpClient.getSlugFromID(input)
		if err != nil {
			log.Fatalf("Could not find problem with ID %s: %v", input, err)
		}
		slug = foundSlug
		fmt.Printf("Found: %s\n", slug)
	}

	query := `
    query getQuestionDetail($titleSlug: String!) {
      question(titleSlug: $titleSlug) {
        questionId
        title
        content
        codeSnippets {
          lang
          langSlug
          code
        }
      }
    }
    `
	payload := GQLRequest{
		OperationName: "getQuestionDetail",
		Variables:     map[string]any{"titleSlug": slug},
		Query:         query,
	}
	var result GQLResponse
	resp, err := HttpClient.post(payload, &result)

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	if resp.IsError() {
		log.Fatalf("API Error: %s", resp.Status())
	}

	q := result.Data.Question
	if q.Title == "" {
		log.Fatal("Problem not found or auth failed. Check your slug and .env")
	}
	markdown, err := htmltomarkdown.ConvertString(q.Content)
	markdown = fmt.Sprintf("# %s\n\n%s", q.Title, markdown)
	if err != nil {
		log.Fatal(err)
	}
	out, err := glamour.Render(markdown, "dark")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)

	f, err := os.Create(fmt.Sprintf("%s.py", slug))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, snippet := range q.CodeSnippets {
		if snippet.LangSlug == "python3" {
			var types []string
			for _, t := range []string{"List", "Optional"} {
				if strings.Contains(snippet.Code, t+"[") {
					types = append(types, t)

				}
			}
			if len(types) > 0 {
				snippet.Code = "from typing import " + strings.Join(types, ", ") + "\n\n" + snippet.Code
			}
			_, err := f.WriteString(snippet.Code)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}

	openInEditor(f.Name())
}
