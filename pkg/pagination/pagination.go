package pagination

type Meta struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

type Params struct {
	Page  int
	Limit int
}

func (p Params) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func NewMeta(total, perPage, currentPage int) Meta {
	count := perPage
	if total-((currentPage-1)*perPage) < perPage {
		count = total - ((currentPage - 1) * perPage)
	}
	if count < 0 {
		count = 0
	}
	totalPages := (total + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}
	return Meta{
		Total:       total,
		Count:       count,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}
}
