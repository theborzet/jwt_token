package pagination

type Paginator struct {
	CurrentPage  int          `json:"current_page"`
	TotalPages   int          `json:"total_pages"`
	PageSize     int          `json:"page_size"`
	TotalItems   int          `json:"total_items"`
	HasPrevious  bool         `json:"has_previous"`
	HasNext      bool         `json:"has_next"`
	PreviousPage int          `json:"previous_page"`
	NextPage     int          `json:"next_page"`
	PageNumbers  []PageNumber `json:"page_numbers"`
}

type PageNumber struct {
	Number   int    `json:"number"`
	URL      string `json:"url"`
	IsActive bool   `json:"is_active"`
}
