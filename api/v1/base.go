package v1

type apiError struct {
	Errors []string `json:"errors"`
}

type PaginationInfo struct {
	Limit          int    `json:"limit"`
	PreviousOffset int    `json:"previous_offset"`
	PreviousUrl    int    `json:"previous_url"`
	CurrentOffset  int    `json:"current_offset"`
	NextOffset     int    `json:"next_offset"`
	NextUrl        string `json:"next_url"`
}
