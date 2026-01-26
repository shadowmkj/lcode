package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/charmbracelet/glamour"
)

func handlePick(httpClient HttpClient, input string) {

	slug := input
	if _, err := strconv.Atoi(input); err == nil {
		fmt.Printf("Searching for Problem ID: %s...\n", input)
		foundSlug, err := httpClient.getSlugFromID(input)
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
	resp, err := httpClient.post(payload, &result)

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
	fmt.Println(out)

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
			snippet.Code = fmt.Sprintf("# @lc app=leetcode id=%s lang=python3\n\n%s\n# @lc code=end\n", q.QuestionID, snippet.Code)
			_, err := f.WriteString(snippet.Code)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}

	cmd := exec.Command("bat")

	reader, writer := io.Pipe()
	cmd.Stdin = reader

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start less: %v\n", err)
		return
	}

	go func() {
		defer writer.Close()
		fmt.Fprint(writer, out)
	}()

	openInEditor(f.Name())
	cmd.Wait()
}

func handleSubmit(httpClient HttpClient, filename string) {
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	content := string(contentBytes)

	questionID := parseQuestionID(content)
	if questionID == "" {
		log.Fatal("Could not find Question ID in file. Did you remove the header?")
	}

	slug := strings.TrimSuffix(filename, ".py")

	fmt.Printf("Submitting %s (ID: %s)...\n", slug, questionID)

	httpClient.setHeader("Referer", fmt.Sprintf("https://leetcode.com/problems/%s/", slug))
	subPayload := SubmissionPayload{
		Lang:       "python3",
		QuestionID: questionID,
		TypedCode:  content,
	}

	var subResp SubmissionResponse
	resp, err := httpClient.client.R().
		SetBody(subPayload).
		SetResult(&subResp).
		Post(fmt.Sprintf("https://leetcode.com/problems/%s/submit/", slug))

	if err != nil {
		log.Fatalf("Submission failed: %s", resp.String())
	}

	if resp.IsError() {
		if resp.StatusCode() == 403 {
			log.Fatal("Submission failed: Forbidden (403). Check your authentication.")
		}
		fmt.Println("Unauthenticated! Please run 'lcode auth' to set up your environment.")
		return
	}

	fmt.Printf("Submission sent! ID: %d. Waiting for results...\n", subResp.SubmissionID)
	for {
		time.Sleep(1 * time.Second)
		var check SubmissionCheckResult
		_, err := httpClient.client.R().
			SetResult(&check).
			Get(fmt.Sprintf("https://leetcode.com/submissions/detail/%d/check/", subResp.SubmissionID))

		if err != nil {
			continue
		}

		if check.State == "SUCCESS" {
			printResult(check)
			break
		}
		fmt.Print(".")
	}
}

func printResult(res SubmissionCheckResult) {
	fmt.Println("\n-----------------------------")
	if res.StatusMsg == "Accepted" {
		fmt.Printf("✅ ACCEPTED\n")
		fmt.Printf("Runtime: %s (Better than %.2f%%)\n", res.StatusRuntime, res.RuntimePercent)
		fmt.Printf("Memory:  %s (Better than %.2f%%)\n", res.StatusMemory, res.MemoryPercent)
	} else {
		fmt.Printf("❌ %s\n", res.StatusMsg)
		fmt.Printf("Testcase: %s\n", res.LastTestcase)
		fmt.Printf("Output: %s\n", res.CodeOutput)
		fmt.Printf("Expected: %s\n", res.ExpectedOutput)
		fmt.Printf("Passed: %d / %d testcases\n", res.TotalCorrect, res.TotalTestcases)
		if res.FullCompileError != "" {
			fmt.Printf("Error: %s\n", res.FullCompileError)
		}
	}
	fmt.Println("-----------------------------")
}


func parseQuestionID(content string) string {
	// Looks for: # @lc app=leetcode id=123 lang=python3
	re := regexp.MustCompile(`id=(\d+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
