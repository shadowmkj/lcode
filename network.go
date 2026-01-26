package main

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	client *resty.Client
}

func (c *HttpClient) initClient() {
	LEETCODE_SESSION, LEETCODE_CSRF_TOKEN := loadEnv()
	c.client = resty.New()
	c.client.SetHeader("User-Agent", UserAgent)
	c.client.SetHeader("Content-Type", "application/json")
	c.client.SetHeader("x-csrftoken", LEETCODE_CSRF_TOKEN)
	c.client.SetCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: LEETCODE_SESSION})
	c.client.SetCookie(&http.Cookie{Name: "csrftoken", Value: LEETCODE_CSRF_TOKEN})
}

func (c *HttpClient) post(payload GQLRequest, result *GQLResponse) (*resty.Response, error) {
	return c.client.R().
		SetBody(payload).
		SetResult(&result).
		Post(GraphQLEndpoint)
}

func (c *HttpClient) get(url string) (*resty.Response, error) {
	return c.client.R().Get(url)
}

func (c *HttpClient) getClient() *resty.Client {
	return c.client
}

func (c *HttpClient) setHeader(key, value string) {
	c.client.SetHeader(key, value)
}

func (c *HttpClient) setCookie(name, value string) {
	c.client.SetCookie(&http.Cookie{Name: name, Value: value})
}

func (c *HttpClient) setCookies(cookies map[string]string) {
	for name, value := range cookies {
		c.client.SetCookie(&http.Cookie{Name: name, Value: value})
	}
}

func (c *HttpClient) getSlugFromID(id string) (string, error) {
	query := `
    query problemsetQuestionList($categorySlug: String, $limit: Int, $filters: QuestionListFilterInput) {
      problemsetQuestionList: questionList(
        categorySlug: $categorySlug
        limit: $limit
        filters: $filters
      ) {
        data {
          questionFrontendId
          titleSlug
        }
      }
    }
    `
	payload := GQLRequest{
		OperationName: "problemsetQuestionList",
		Variables: map[string]any{
			"categorySlug": "",
			"limit":        10, // Fetch a few to ensure exact match
			"filters":      map[string]string{"searchKeywords": id},
		},
		Query: query,
	}

	var result ProblemListResponse
	_, err := c.client.R().SetBody(payload).SetResult(&result).Post(GraphQLEndpoint)
	if err != nil {
		return "", err
	}

	// Iterate to find exact ID match (search is fuzzy)
	for _, p := range result.Data.ProblemsetQuestionList.Data {
		if p.QuestionFrontendId == id {
			return p.TitleSlug, nil
		}
	}
	return "", fmt.Errorf("id not found")
}
