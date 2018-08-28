package search

type ValidationError struct {
	Err string   `json:"reason"`
	Code int64   `json:"code"`
	Field string `json:"field"`
}



type ApiError struct {
	Code int64               `json:"code"`
	Reason string            `json:"reason"`
	Errors []ValidationError `json:validationErrors`
}
