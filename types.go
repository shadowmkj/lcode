package main

type GQLRequest struct {
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
	Query         string         `json:"query"`
}

type GQLResponse struct {
	Data struct {
		Question struct {
			QuestionID   string `json:"questionId"`
			Title        string `json:"title"`
			Content      string `json:"content"` // This is HTML
			CodeSnippets []struct {
				Lang     string `json:"lang"`
				LangSlug string `json:"langSlug"`
				Code     string `json:"code"`
			} `json:"codeSnippets"`
		} `json:"question"`
	} `json:"data"`
}

type QuestionDetailResponse struct {
	Data struct {
		Question struct {
			QuestionID   string `json:"questionId"`
			Title        string `json:"title"`
			Content      string `json:"content"`
			CodeSnippets []struct {
				Lang     string `json:"lang"`
				LangSlug string `json:"langSlug"`
				Code     string `json:"code"`
			} `json:"codeSnippets"`
		} `json:"question"`
	} `json:"data"`
}

type ProblemListResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Data []struct {
				QuestionFrontendId string `json:"questionFrontendId"`
				TitleSlug          string `json:"titleSlug"`
			} `json:"data"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

type SubmissionPayload struct {
	Lang       string `json:"lang"`
	QuestionID string `json:"question_id"`
	TypedCode  string `json:"typed_code"`
}

type SubmissionResponse struct {
	SubmissionID int64 `json:"submission_id"`
}

type SubmissionCheckResult struct {
	State            string  `json:"state"` // PENDING, STARTED, SUCCESS
	StatusMsg        string  `json:"status_msg"`
	TotalCorrect     int     `json:"total_correct"`
	TotalTestcases   int     `json:"total_testcases"`
	StatusRuntime    string  `json:"status_runtime"`
	StatusMemory     string  `json:"status_memory"`
	RuntimePercent   float64 `json:"runtime_percentile"`
	MemoryPercent    float64 `json:"memory_percentile"`
	FullCompileError string  `json:"full_compile_error"`
	LastTestcase     string  `json:"last_testcase"`
	CodeOutput       string  `json:"code_output"`
	ExpectedOutput   string  `json:"expected_output"`
}
